package metadata

import (
	pb_collector "github.com/eviltomorrow/king/lib/grpc/pb/king-collector"
	pb_storage "github.com/eviltomorrow/king/lib/grpc/pb/king-storage"
)

var (
	IPVar string

	beginVar string
	endVar   string

	modeVar string
)

var (
	ClientCollector pb_collector.CollectorClient
	ClientStorage   pb_storage.StorageClient
)
