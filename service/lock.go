package service

import (
	"time"

	"errors"

	"github.com/coreos/etcd/client"
	"github.com/revinate/go-utils/helper"
	"github.com/satori/go.uuid"
	"context"
	"fmt"
)

const (
	LockDir = "service/lock"
)

type LockService struct {
	ctx context.Context
	client client.Client
}

func NewLockService(connUrl string, ctx context.Context) (*LockService, error) {
	client, err := getEtcdClient(connUrl)
	return &LockService{client: client, ctx: ctx}, err
}

func getEtcdClient(connUrl string) (client.Client, error) {
	config := client.Config{
		Endpoints:               []string{"http://" + connUrl},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	etcdClient, err := client.New(config)
	return etcdClient, err
}

func (l LockService) AcquireLock(name string, ttl time.Duration, isBlocking bool) (string, error) {
	if ttl <= 0 {
		return "", errors.New("TTL has to be above or equal to 1 second")
	}
	api := client.NewKeysAPI(l.client)
	value := uuid.NewV4().String()
	name = getLockName(name)
	ticker := time.NewTicker(ttl / 2)
	defer ticker.Stop()
	var err error
loop:
	for {
		select {
		case <-ticker.C:
			options := client.SetOptions{PrevExist: client.PrevNoExist, TTL: ttl}
			_, err = api.Set(l.ctx, name, value, &options)
			if err == nil { // Got the lock
				return value, nil
			}
			if !isBlocking {
				break loop
			}
			helper.Debug("Waiting to acquire lock: "+name, nil)
		case <-l.ctx.Done():
			break loop
		}
	}
	return "", err
}

func (l LockService) RenewLock(name string, value string, ttl time.Duration) error {
	if ttl <= 0 {
		return errors.New("TTL has to be above or equal to 1 second")
	}
	api := client.NewKeysAPI(l.client)
	name = getLockName(name)
	options := client.SetOptions{PrevValue: value, TTL: ttl}
	_, err := api.Set(context.Background(), name, value, &options)
	return err
}

func (l LockService) AcquireAndKeepLock(name string, ttl time.Duration, isBlocking bool) (string, error) {
	value, err := l.AcquireLock(name, ttl, isBlocking)
	if err != nil {
		return value, err
	}
	go func() {
		ticker := time.NewTicker(ttl / 2)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := l.RenewLock(name, value, ttl)
				if err != nil {
					// TODO: Unsure what to do if RenewLock fails. It should not, but what if?
					value, err = l.AcquireLock(name, ttl, isBlocking)
					helper.Error(err, "Renewing lock failed")
				}
			case <-l.ctx.Done():
				return
			}
		}
	}()
	return value, nil
}

func (l LockService)ReleaseLock(name string, value string) error {
	name = getLockName(name)
	api := client.NewKeysAPI(l.client)
	_, err := api.Delete(context.Background(), name, &client.DeleteOptions{PrevValue: value, Recursive: true})
	return err
}

func getLockName(name string) string {
	return fmt.Sprintf("%s/%s",LockDir,name)
}