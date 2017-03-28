package client

type Event struct {
	Name string
}

type SmartConf interface {
	LoadObject(objectName string, v interface{}) (err error)
	Notify() (event <-chan *Event, err error)
}

func NewLocalSmartConf(path string) SmartConf {
	return newLocalConf(path)
}

func NewRemoteSmartConf(addresses []string, path string) SmartConf {
	return newRemoteEtcdSmartConf(addresses, path)
}
