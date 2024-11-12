package metadata

import (
	pb_collector "github.com/eviltomorrow/king/lib/grpc/pb/king-collector"
	pb_storage "github.com/eviltomorrow/king/lib/grpc/pb/king-storage"
)

var (
	IP string

	begin string
	end   string

	mode string
)

var (
	ClientCollector pb_collector.CollectorClient
	ClientStorage   pb_storage.StorageClient
)
