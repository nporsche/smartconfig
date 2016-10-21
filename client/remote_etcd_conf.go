package client

import (
	//	"context"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/BurntSushi/toml"
	etcd "github.com/coreos/etcd/client"
)

type remoteEtcdMeConf struct {
	etcdClient  etcd.Client
	path        string
	subscribers []chan *Event
}

func newRemoteEtcdMeConf(endPoints []string, path string) *remoteEtcdMeConf {
	rc := &remoteEtcdMeConf{
		path:        path,
		subscribers: make([]chan *Event, 0),
	}
	cfg := etcd.Config{
		Endpoints: endPoints,
		Transport: etcd.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	var err error
	rc.etcdClient, err = etcd.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	rc.monitor()

	return rc
}

func (rc *remoteEtcdMeConf) LoadObject(objectName string, v interface{}) error {
	api := etcd.NewKeysAPI(rc.etcdClient)
	resp, err := api.Get(context.Background(), path.Join(rc.path, objectName), nil)
	if err != nil {
		return err
	}

	if strings.HasSuffix(objectName, ".toml") {
		_, err = toml.Decode(resp.Node.Value, v)
		return err
	} else if strings.HasSuffix(objectName, ".json") {
		return json.Unmarshal([]byte(resp.Node.Value), v)
	} else {
		return fmt.Errorf("cannot deduce the format type from %s", objectName)
	}
}

func (rc *remoteEtcdMeConf) Notify() (event <-chan *Event, err error) {
	ev := make(chan *Event, 1)
	rc.subscribers = append(rc.subscribers, ev)

	return ev, nil
}

func (rc *remoteEtcdMeConf) monitor() error {
	api := etcd.NewKeysAPI(rc.etcdClient)
	wo := &etcd.WatcherOptions{Recursive: true}
	watcher := api.Watcher(rc.path, wo)
	go func() {
		for {
			resp, err := watcher.Next(context.Background())
			if err == nil {
				if resp.Action == "set" {
					ev := &Event{
						Name: path.Base(resp.Node.Key),
					}
					for _, sub := range rc.subscribers {
						select {
						case sub <- ev:
						default:
						}
					}
				}
			} else {
				//log.Println(err)
			}
		}
	}()

	return nil
}
