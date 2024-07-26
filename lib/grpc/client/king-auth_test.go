package client

import (
	"context"
	"testing"
	"time"

	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-auth"
)

var (
	account  = "shepard"
	password = "123456"
)

func TestAuthAndVerify(t *testing.T) {
	client, closeFunc, err := NewAuthWithTarget("127.0.0.1:5277")
	if err != nil {
		t.Fatal(err)
	}
	defer closeFunc()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	token, err := client.Auth(ctx, &pb.AuthReq{Method: pb.AuthReq_PASSWORD, Credential: &pb.AuthReq_Credential{Account: account, Key: password}})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("token: %v\r\n", token.String())

	if _, err = client.VerifyToken(ctx, token); err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)
	token, err = client.RenewToken(ctx, token)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("token: %v\r\n", token.String())
}

func TestRegister(t *testing.T) {
	client, closeFunc, err := NewAuthWithTarget("127.0.0.1:5277")
	if err != nil {
		t.Fatal(err)
	}
	defer closeFunc()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	value, err := client.Register(ctx, &pb.RegisterReq{Account: account, Password: password})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("id: %s\r\n", value.Value)
}
