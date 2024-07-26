package controller

import (
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-account"
	"google.golang.org/grpc"
)

type TransactionRecord struct {
	pb.UnimplementedTransactionRecordServer
}

func NewTransactionRecord() *TransactionRecord {
	return &TransactionRecord{}
}

func (g *TransactionRecord) Service() func(*grpc.Server) {
	return func(server *grpc.Server) {
		pb.RegisterTransactionRecordServer(server, g)
	}
}
