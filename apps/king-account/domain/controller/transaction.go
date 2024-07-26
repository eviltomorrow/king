package controller

import (
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-account"
	"google.golang.org/grpc"
)

type Transaction struct {
	pb.UnimplementedTransactionServer
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

func (g *Transaction) Service() func(*grpc.Server) {
	return func(server *grpc.Server) {
		pb.RegisterTransactionServer(server, g)
	}
}
