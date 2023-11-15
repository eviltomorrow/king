package server

import clientv3 "go.etcd.io/etcd/client/v3"

type HTTP struct {
	EtcdClient *clientv3.Client
	Host       string
	Port       int
	AppName    string
}

func (h *HTTP) Startup() error {
	return nil
}

func (h *HTTP) Stop() error {
	return nil
}
