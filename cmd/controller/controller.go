package main

import (
	"context"
	"gitcom.com/valentin-carl/f4s/pkg/controller"
	"github.com/fatih/color"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

const ETCD_URL = "localhost:2379"

func main() {
	log.Print(color.GreenString("starting controller"))

	//controller.StartController()

	// etcd
	etcdCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		log.Print("cancelling etcd context")
		cancel()
	}()

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{ETCD_URL},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		log.Fatal(color.RedString("could not create connection to etcd node: %s", err.Error()))
	}
	defer etcdClient.Close()

	log.Print(color.GreenString("connected to etcd"))

	err = controller.InitializeFunctionInstanceRecord(etcdCtx, etcdClient, "echo", "test1", "localhost:8080", 10)
	if err != nil {
		log.Print(err.Error())
	} else {
		log.Print(color.GreenString("works"))
	}

	err = controller.InitializeFunctionInstanceRecord(etcdCtx, etcdClient, "echo", "test2", "localhost:8081", 11)
	if err != nil {
		log.Print(err.Error())
	} else {
		log.Print(color.GreenString("works"))
	}

	err = controller.InitializeFunctionInstanceRecord(etcdCtx, etcdClient, "echo", "test3", "localhost:8082", 12)
	if err != nil {
		log.Print(err.Error())
	} else {
		log.Print(color.GreenString("works"))
	}

	cwc, err := controller.DecrementCurrentWorkerCount(etcdCtx, etcdClient, "echo", "test2")
	if err != nil {
		log.Print(err.Error())
	} else {
		log.Printf("updated current worker count %d", cwc)
	}

	workerCounts, err := controller.GetCurrentWorkersForFunction(etcdCtx, etcdClient, "echo")
	if err != nil {
		log.Print("error while trying to get current worker amounts")
	} else {
		for key, value := range workerCounts {
			log.Print(key, ": ", value)
		}
	}

	// TODO move most of this code to test file

}
