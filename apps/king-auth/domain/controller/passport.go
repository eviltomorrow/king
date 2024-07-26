package controller

import (
	"context"
	"fmt"

	"github.com/eviltomorrow/king/apps/king-auth/domain/service"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Passport struct {
	pb.UnimplementedPassportServer
}

func NewPassport() *Passport {
	return &Passport{}
}

func (p *Passport) Service() func(*grpc.Server) {
	return func(server *grpc.Server) {
		pb.RegisterPassportServer(server, p)
	}
}

func (p *Passport) Register(ctx context.Context, req *pb.RegisterReq) (*wrapperspb.StringValue, error) {
	if req == nil {
		return nil, fmt.Errorf("req is nil")
	}

	accountId, err := service.PassportWithRegister(ctx, req.Account, req.Password)
	if err != nil {
		return nil, err
	}
	return &wrapperspb.StringValue{Value: accountId}, nil
}

func (p *Passport) Auth(ctx context.Context, req *pb.AuthReq) (*pb.Token, error) {
	if req == nil || req.Credential == nil {
		return nil, fmt.Errorf("req/credential is nil")
	}

	var authMethod service.PassportAuthMethod

	switch req.Method {
	case pb.AuthReq_PASSWORD:
		authMethod = service.PASSWORD
	case pb.AuthReq_SMS:
		authMethod = service.SMS
	default:
		return nil, fmt.Errorf("invalid auth method")
	}

	passport, err := service.PassportWithAuth(ctx, authMethod, req.Credential.Account, req.Credential.Key)
	if err != nil {
		return nil, err
	}

	token, err := service.TokenWithApply(ctx, passport.Id, "admin")
	if err != nil {
		return nil, err
	}

	return &pb.Token{AccessToken: token.AccessToken, TokenType: token.TokenType, RefreshToken: token.RefreshToken, ExpiresIn: token.ExpiresIn}, nil
}

func (p *Passport) RenewToken(ctx context.Context, req *pb.Token) (*pb.Token, error) {
	if req == nil {
		return nil, fmt.Errorf("req is nil")
	}

	newToken, err := service.TokenWithRenew(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	resp := &pb.Token{
		AccessToken:  newToken.AccessToken,
		TokenType:    newToken.TokenType,
		RefreshToken: newToken.RefreshToken,
		ExpiresIn:    newToken.ExpiresIn,
	}

	return resp, nil
}

func (p *Passport) VerifyToken(ctx context.Context, req *pb.Token) (*emptypb.Empty, error) {
	if req == nil {
		return nil, fmt.Errorf("req is nil")
	}

	token := service.Token{
		AccessToken:  req.AccessToken,
		TokenType:    req.TokenType,
		RefreshToken: req.RefreshToken,
		ExpiresIn:    req.ExpiresIn,
	}

	return &emptypb.Empty{}, service.TokenWithVerify(ctx, token)
}

func (p *Passport) RevokeToken(ctx context.Context, req *pb.Token) (*emptypb.Empty, error) {
	if req == nil {
		return nil, fmt.Errorf("req is nil")
	}
	token := service.Token{
		AccessToken:  req.AccessToken,
		TokenType:    req.TokenType,
		RefreshToken: req.RefreshToken,
		ExpiresIn:    req.ExpiresIn,
	}
	return &emptypb.Empty{}, service.TokenWithRevokeByToken(ctx, token)
}

func (p *Passport) Exist(ctx context.Context, req *wrapperspb.StringValue) (*wrapperspb.BoolValue, error) {
	if req == nil {
		return nil, fmt.Errorf("req is nil")
	}
	ok, err := service.PassportWithExist(ctx, req.Value)
	return &wrapperspb.BoolValue{Value: ok}, err
}

func (p *Passport) Lock(ctx context.Context, req *wrapperspb.StringValue) (*emptypb.Empty, error) {
	if req == nil {
		return nil, fmt.Errorf("req is nil")
	}

	if err := service.TokenWithRevokeByAccountId(ctx, req.Value); err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, service.PassportWithChangeStatus(ctx, service.LOCK, req.Value)
}

func (p *Passport) Unlock(ctx context.Context, req *wrapperspb.StringValue) (*emptypb.Empty, error) {
	if req == nil {
		return nil, fmt.Errorf("req is nil")
	}
	return &emptypb.Empty{}, service.PassportWithChangeStatus(ctx, service.NORMAL, req.Value)
}

func (p *Passport) Get(ctx context.Context, req *wrapperspb.StringValue) (*pb.User, error) {
	if req == nil {
		return nil, fmt.Errorf("req is nil")
	}
	passport, err := service.PassportWithGet(ctx, req.Value)
	if err != nil {
		return nil, err
	}
	return &pb.User{Id: passport.Id, Account: passport.Account, Code: passport.Code, Email: passport.Email, Phone: passport.Phone, Status: passport.Status, RegisterDatetime: passport.CreateTimestamp.Unix()}, nil
}

func (p *Passport) Remove(ctx context.Context, req *wrapperspb.StringValue) (*emptypb.Empty, error) {
	if req == nil {
		return nil, fmt.Errorf("req is nil")
	}
	return &emptypb.Empty{}, status.Error(codes.Unimplemented, "unimplemented")
}

func (p *Passport) ModifyPassword(ctx context.Context, req *pb.ModifyPasswordReq) (*emptypb.Empty, error) {
	if req == nil {
		return nil, fmt.Errorf("req is nil")
	}
	if err := service.TokenWithRevokeByAccountId(ctx, req.Id); err != nil {
		return &emptypb.Empty{}, err
	}

	if err := service.PassportWithChangePassword(ctx, req.Password, req.Id); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
