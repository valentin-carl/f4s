package controller

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"log"
	"net"
	"regexp"
)

const (
	GRPC_PORT = 54321
)

var (
	// TODO make sure these are initialized before use
	EtcdClient *clientv3.Client
	EtcdCtx    context.Context
)

type FunctionStateServer struct {
	UnimplementedFunctionStateServiceServer
}

func (f *FunctionStateServer) GetInstance(ctx context.Context, r *FunctionInstanceRequest) (*FunctionInstanceResponse, error) {

	// TODO
	//  - get current worker information from etcd
	//  - if instances are available: pick the one with the least amount of workers currently, increment current worker counter
	//  - else: create new container, store the info in etcd, increment current worker counter
	//  - return function url of instance that was chosen

	fn := r.GetFunctionName()

	cw, err := GetCurrentWorkersForFunction(EtcdCtx, EtcdClient, fn)
	if err != nil {
		log.Print(color.YellowString("grpc: error while trying to get currentWorkerCounts for %s: %s", fn, err.Error()))
		return nil, err
	}

	if len(cw) > 0 {
		const MAX_INT = 2147483647
		minAvailableInstance, currMinCount := "", MAX_INT
		// key structure: '/functions/<functionName>/containers/<containerName>/currentWorkerCount'
		for key, value := range cw {
			if value < currMinCount {
				minAvailableInstance, currMinCount = key, value
			}
		}

		if minAvailableInstance == "" || currMinCount == MAX_INT {
			// some instances already exist but no more workers left to add
			log.Printf("no available instance found for %s, creating new function instances", fn)
			goto NewInstance

		} else {
			// found available instance
			containerId := regexp.MustCompile(`.*/containers/(.*)/currentWorkerCount`).FindStringSubmatch(minAvailableInstance)[1] // TODO test + add error handling (array index, ...)
			url := GetUrlForContainer(ctx, EtcdClient, fn, containerId)
			log.Printf("found url for %s:%s: %s", fn, containerId, url)
			cwc, err := IncrementCurrentWorkerCount(EtcdCtx, EtcdClient, fn, containerId)
			if err != nil {
				log.Printf("error while incrementing currentWorkerCounter for %s:%s: %s", fn, containerId, err.Error())
				return nil, err
			}
			log.Printf("incremented currentWorkerCounter for %s:%s to %d", fn, containerId, cwc)
			return &FunctionInstanceResponse{
				Url: url,
			}, nil
		}
	}

	// create new instance instead
NewInstance:
	log.Printf("creating new function instance for %s", fn)
	// TODO
	//  no available instance found
	//  => add when doing the Docker stuff

}

func (f *FunctionStateServer) NotifyStreamClosed(ctx context.Context, r *StreamClosedNotification) (*StreamClosedNotificationResponse, error) {

	// TODO
	//  - find correct function instance
	//  - decrement counter of function instance
	//  - if currentWorkerCounter is now zero: destroy the container

}

func StartController() error {

	// TODO create etcd context & client

	// create & start gRPC server for updating platform state
	// TODO make sure the networking works once the proxy and controller are in separate containers
	grpcAddr := fmt.Sprintf(":%d", GRPC_PORT)
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Print(color.RedString("controller could not create TCP listener to wait for gRPC connections: %s", err.Error()))
	}

	grpcServer := grpc.NewServer()
	RegisterFunctionStateServiceServer(grpcServer, &FunctionStateServer{})

	go func() {
		log.Printf("gRPC server listening at %s", grpcAddr)
		err := grpcServer.Serve(grpcListener)
		if err != nil {
			log.Print(color.RedString("gRPC server finished with error: %s", err.Error()))
		}
		log.Print(color.GreenString("gRPC server done"))
	}()

	// TODO create & start HTTP server for deployments

	// TODO if (very) bored: change this to wait for interrupt (and other components too) and write code to stop everything with one command
	<-make(chan any)
}
