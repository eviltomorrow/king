package impl

import (
	"github.com/eviltomorrow/king/apps/king-notification/conf"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
)

type NotificationServer struct {
	NFTY *conf.NFTY
	pb.UnimplementedNotificationServer
}

// func (s *NotificationServer)
