package setting

import "time"

const (
	DEFUALT_HANDLE_10_SECOND     = 10 * time.Second
	DEFUALT_HANDLE_30_SECOND     = 30 * time.Second
	BATCH_HANDLE_LIMIT           = 50
	DB_QUERY_LIMIT               = 50
	GRPC_UNARY_TIMEOUT_10_SECOND = DEFUALT_HANDLE_10_SECOND
	GRPC_UNARY_TIMEOUT_30_SECOND = DEFUALT_HANDLE_30_SECOND
	GRPC_UNARY_TIMEOUT_60_SECOND = 60 * time.Second
)
