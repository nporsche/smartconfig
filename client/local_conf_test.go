package client

import (
	"log"
	"testing"
)

type sample struct {
	Name string
	Age  int
}

func TestLocalConf(t *testing.T) {
	mc := NewLocalMeConf("./")
	var jsObj sample
	if err := mc.LoadObject("local_conf_sample.json", &jsObj); err != nil {
		t.Fatal(err)
	} else {
		log.Printf("local_conf_sample.json object content: %#v", jsObj)
	}

	var tomlObj sample
	if err := mc.LoadObject("local_conf_sample.toml", &tomlObj); err != nil {
		t.Fatal(err)
	} else {
		log.Printf("local_conf_sample.toml object content: %#v", tomlObj)
	}

	if ev, err := mc.Notify(); err != nil {
		t.Fatal(err)
	} else {
		go func() {
			for {
				x := <-ev
				var obj sample
				mc.LoadObject(x.Name, &obj)
				log.Printf("%s has been changed to %#v", x.Name, obj)
			}
		}()
	}

	/*
		uncomment below to check update logic

		log.Print("please edit file then you can see the notify in 1 minutes")
		time.Sleep(1 * time.Minute)
	*/
}
