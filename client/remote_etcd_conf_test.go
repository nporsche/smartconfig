package client

import (
	"log"
	"testing"
	"time"
)

func TestRemoteEtcdConf(t *testing.T) {
	//etcdctl set /config/proj1/sample.json "{\"Name\":\"name_sample_json\",\"Age\":21}"
	mc := NewRemoteSmartConf([]string{"http://127.0.0.1:2379"}, "/config/proj1")
	var jsonObj sample
	if err := mc.LoadObject("sample.json", &jsonObj); err != nil {
		t.Fatal(err)
	}
	log.Printf("sample.json %#v", jsonObj)

	notify, err := mc.Notify()
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for {
			ev := <-notify
			var obj sample
			if err := mc.LoadObject(ev.Name, &obj); err != nil {
				t.Fatal(err)
			}
			log.Printf("%s has been changed to %#v", ev.Name, obj)
		}
	}()

	log.Print("please edit file then you can see the notify in 1 minutes")
	time.Sleep(1 * time.Minute)
}
