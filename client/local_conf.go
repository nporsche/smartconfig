package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/howeyc/fsnotify"
)

type localConf struct {
	path        string
	watcher     *fsnotify.Watcher
	watchOK     bool
	subscribers []chan *Event
}

func newLocalConf(path string) (lf *localConf) {
	lc := &localConf{
		path:        path,
		subscribers: make([]chan *Event, 0),
	}
	if err := lc.monitor(); err == nil {
		//log
		lc.watchOK = true
	} else {
		lc.watchOK = false
	}

	return lc
}

func (lc *localConf) LoadObject(objectName string, v interface{}) error {
	if strings.HasSuffix(objectName, ".toml") {
		_, err := toml.DecodeFile(path.Join(lc.path, objectName), v)
		return err
	} else if strings.HasSuffix(objectName, ".json") {
		bs, err := ioutil.ReadFile(path.Join(lc.path, objectName))
		if err != nil {
			return err
		}

		return json.Unmarshal(bs, v)
	} else {
		return fmt.Errorf("cannot distinguish file type from %s", objectName)
	}

	return nil
}

func (lc *localConf) Notify() (event <-chan *Event, err error) {
	if !lc.watchOK {
		return nil, errors.New("file watch failed")
	}
	ev := make(chan *Event, 1)
	lc.subscribers = append(lc.subscribers, ev)
	return ev, nil
}

func (lc *localConf) Close() {
	lc.watcher.Close()
}

func (lc *localConf) monitor() error {
	var err error
	lc.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// Process events
	go func() {
		for {
			select {
			case ev := <-lc.watcher.Event:
				{
					if ev.IsModify() && !ev.IsAttrib() {
						event := &Event{
							Name: ev.Name,
						}
						for _, sub := range lc.subscribers {
							select {
							case sub <- event:
							default:
							}
						}
					}
				}

				//case err := <-lc.watcher.Error:
				//ignore
			}
		}
	}()

	err = lc.watcher.Watch(lc.path)
	if err != nil {
		return err
	}

	return nil
}
