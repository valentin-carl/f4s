package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"strconv"
	"strings"
	"sync"
)

const (
	curr = "currentWorkerCount"
	maxw = "maxWorkerCount"
)

var InstanceLocks = make(map[string]sync.Mutex) // TODO might not be necessary anymore if I get the transactions right

func InitializeFunctionInstanceRecord(ctx context.Context, c *clientv3.Client, functionName string, containerId string, functionUrl string, maxWorkers int) error {

	// 1 function instance == 1 docker container
	// 1 function can have multiple function instances
	// 1 function instance can have multiple workers

	base := fmt.Sprintf("functions/%s/containers/%s/", functionName, containerId)

	txResponse, err := c.Txn(ctx).If(
		clientv3.Compare(clientv3.Version(base+"url"), "=", 0),
		clientv3.Compare(clientv3.Version(base+curr), "=", 0),
		clientv3.Compare(clientv3.Version(base+maxw), "=", 0),
	).Then(
		clientv3.OpPut(base+"url", functionUrl),
		clientv3.OpPut(base+curr, "0"),
		clientv3.OpPut(base+maxw, strconv.Itoa(maxWorkers)),
	).Commit()

	if err != nil {
		log.Print(color.YellowString("etcd error while trying to create function record: %s", err.Error()))
		return err
	}

	log.Printf("transaction to create instance record %s:%s returned %v", functionName, containerId, txResponse.Succeeded)

	for _, r := range txResponse.Responses {
		log.Print(r.String())
	}

	if !txResponse.Succeeded {
		log.Print("transaction unsuccessful")
		return errors.New("did not create instance record because one of its fields already exists")
	}

	return nil
}

func GetCurrentWorkersForFunction(ctx context.Context, client *clientv3.Client, functionName string) (map[string]int, error) {

	res := make(map[string]int)
	_ = res

	// TODO make sure this works and GoLand is just giving me a weird warning about nothing
	l := InstanceLocks[functionName]
	l.Lock()
	defer l.Unlock()

	r, err := client.Get(ctx, fmt.Sprintf("functions/%s/containers", functionName), clientv3.WithPrefix())
	if err != nil {
		log.Print(color.YellowString("etcd error while trying to get containers for %s: %s", functionName, err.Error()))
		return nil, err
	}

	for _, c := range r.Kvs {
		key := string(c.Key)
		value := string(c.Value)
		if strings.Contains(key, "/"+curr) {
			cwc, err := strconv.Atoi(value)
			if err != nil {
				log.Print(color.YellowString("etcd error could not convert value %s of key % to int", value, key))
				continue
			}
			res[key] = cwc
		}
	}

	return res, nil
}

func IncrementCurrentWorkerCount(ctx context.Context, client *clientv3.Client, functionName string, containerId string) (int, error) {

	// returns (updated) currentWorkerCount, error

	basePath := fmt.Sprintf("functions/%s/containers/%s/", functionName, containerId)

	if _, ok := InstanceLocks[functionName]; !ok {
		InstanceLocks[functionName] = sync.Mutex{}
		log.Printf("lock for function %s created", functionName)
	}

	// TODO make sure this works and GoLand is just giving me a weird warning about nothing
	l := InstanceLocks[functionName]
	l.Lock()
	defer l.Unlock()

	currentWorkerCountResult, err := client.Get(ctx, basePath+curr)
	if err != nil {
		log.Print(color.YellowString("etcd error could not get currentWorkerCount for %s:%s: %s", functionName, containerId, err.Error()))
		return -1, err
	}
	if len(currentWorkerCountResult.Kvs) != 1 {
		errorMsg := fmt.Sprintf("etcd error unexpected amount of results for currentWorkerCount of %s:%s", functionName, containerId)
		log.Print(color.YellowString(errorMsg))
		return -1, errors.New(errorMsg)
	}
	currentWorkerCount, err := strconv.Atoi(string(currentWorkerCountResult.Kvs[0].Value))
	if err != nil {
		log.Print(color.YellowString("etcd error could not convert currentWorkerCount result to int: %s", err.Error()))
		return -1, err
	}

	newCurrentWorkerCount := currentWorkerCount + 1
	txResult, err := client.Txn(ctx).If(
		// TODO figure out why "<" and ">" only work up to currentWorkerCount 2
		clientv3.Compare(clientv3.Value(basePath+maxw), "!=", strconv.Itoa(currentWorkerCount)),
	).Then(
		clientv3.OpPut(basePath+curr, strconv.Itoa(newCurrentWorkerCount)),
	).Commit()
	if err != nil {
		log.Print(color.YellowString("etcd error while trying to increment currentWorkerCount for %s:%s: %s", functionName, containerId, err.Error()))
		return -1, err
	}

	log.Print(txResult)

	if txResult.Succeeded {
		log.Printf("etcd currently incremented currentWorkerCounter for %s:%s", functionName, containerId)
		return newCurrentWorkerCount, nil
	} else {
		errorMsg := fmt.Sprintf("etcd transaction to increment currentWorkerCount for %s:%s did not succeed", functionName, containerId)
		log.Print(color.YellowString(errorMsg))
		return -1, errors.New(errorMsg)
	}
}

func DecrementCurrentWorkerCount(ctx context.Context, client *clientv3.Client, functionName string, containerId string) (int, error) {

	// returns (updated) currentWorkerCount, error

	basePath := fmt.Sprintf("functions/%s/containers/%s/", functionName, containerId)

	if _, ok := InstanceLocks[functionName]; !ok {
		InstanceLocks[functionName] = sync.Mutex{}
		log.Printf("lock for function %s created", functionName)
	}

	// TODO make sure this works and GoLand is just giving me a weird warning about nothing
	l := InstanceLocks[functionName]
	l.Lock()
	defer l.Unlock()

	currentWorkerCountResult, err := client.Get(ctx, basePath+curr)
	if err != nil {
		log.Print(color.YellowString("etcd error could not get currentWorkerCount for %s:%s: %s", functionName, containerId, err.Error()))
		return -1, err
	}
	if len(currentWorkerCountResult.Kvs) != 1 {
		errorMsg := fmt.Sprintf("etcd error unexpected amount of results for currentWorkerCount of %s:%s", functionName, containerId)
		log.Print(color.YellowString(errorMsg))
		return -1, errors.New(errorMsg)
	}
	currentWorkerCount, err := strconv.Atoi(string(currentWorkerCountResult.Kvs[0].Value))
	if err != nil {
		log.Print(color.YellowString("etcd error could not convert currentWorkerCount result to int: %s", err.Error()))
		return -1, err
	}

	newCurrentWorkerCount := currentWorkerCount - 1
	txResult, err := client.Txn(ctx).If(
		// TODO figure out why "<" and ">" only work up to currentWorkerCount 2
		clientv3.Compare(clientv3.Value(basePath+curr), "!=", strconv.Itoa(0)),
	).Then(
		clientv3.OpPut(basePath+curr, strconv.Itoa(newCurrentWorkerCount)),
	).Commit()
	if err != nil {
		log.Print(color.YellowString("etcd error while trying to decrement currentWorkerCount for %s:%s: %s", functionName, containerId, err.Error()))
		return -1, err
	}

	log.Print(txResult)

	if txResult.Succeeded {
		log.Printf("etcd currently decremented currentWorkerCounter for %s:%s", functionName, containerId)
		return newCurrentWorkerCount, nil
	} else {
		errorMsg := fmt.Sprintf("etcd transaction to decrement currentWorkerCount for %s:%s did not succeed", functionName, containerId)
		log.Print(color.YellowString(errorMsg))
		return -1, errors.New(errorMsg)
	}
}

// TODO add code to delete key/value pairs if function instances are destroyed
