package client

type Event struct {
	Name string
}

type MeConf interface {
	LoadObject(objectName string, v interface{}) (err error)
	Notify() (event <-chan *Event, err error)
	Close()
}

func NewLocalMeConf(path string) MeConf {
	return newLocalConf(path)
}

func NewRemoteMeConf(addresses []string, path string) MeConf {
	return newRemoteEtcdMeConf(addresses, path)
}
