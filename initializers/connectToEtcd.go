package initializers

import (
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var Etcd_cli *clientv3.Client

func ConnectToEtcd() {

	var err error

	Etcd_cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.144.129:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
}
