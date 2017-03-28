# smartconf
smartconf是提供配置服务的SDK, 提供本地或者remote etcd的配置信息装载以及配置更改事件，方便使用者实现配置文件热加载。

##本地配置文件  
创建SmartConf对象,传入本地path:

```
client.NewLocalSmartConf("./");
```
将本地的文件装载成内存对象，对象由调用者自己定义,第一个参数是path路径下的文件名,必须以.toml和.json结尾:

```
err := mc.LoadObject("local_conf_sample.json", &jsObj);
```
可以通过notify接口获取更新对象更新事件:

```
if ev, err := mc.Notify(); err != nil {
         log.Fatal(err)
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
```
一个完整的本地配置文件例子:

```
package main

import (
     "github.com/Meiqia/SmartConf/client"
     "log"
     "time"
)

type sample struct {
    Name string
    Age  int
}

func main() {
    mc := client.NewLocalSmartConf("./")
    var jsObj sample
    if err := mc.LoadObject("local_conf_sample.json", &jsObj); err != nil {
        log.Fatal(err)
    } else {
        log.Printf("local_conf_sample.json object content: %#v", jsObj)
    }
    var tomlObj sample
    if err := mc.LoadObject("local_conf_sample.toml", &tomlObj); err != nil {
        log.Fatal(err)
    } else {
        log.Printf("local_conf_sample.toml object content: %#v", tomlObj)
    }

    if ev, err := mc.Notify(); err != nil {
        log.Fatal(err)
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
    log.Print("please edit file then you can see the notify in 1 minutes")
    time.Sleep(1 * time.Minute)
}
```

##远程配置
创建SmartConf对象,传入etcd endpoints例如http://127.0.0.1:2379,第二个参数是etcd中的path:

```
client.NewRemoteSmartConf([]string{"http://127.0.0.1:2379"}, "/config/proj1");
```
将path下的对象装载成内存对象，对象由调用者自己定义，第一个参数是etcd中在该path下的key,必须以.toml和.json结尾:

```
err := mc.LoadObject("local_conf_sample.json", &jsObj);
```

可以通过notify接口获取更新对象更新事件:

```
if ev, err := mc.Notify(); err != nil {
         log.Fatal(err)
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
```
一个完整远程配置文件的例子:

```
package main

import (
    "github.com/Meiqia/SmartConf/client"
    "log"
    "time"
)

type sample struct {
    Name string
    Age  int
}

func main() {
    mc := client.NewRemoteSmartConf([]string{"http://127.0.0.1:2379"}, "/config/proj1")
    var jsonObj sample
    if err := mc.LoadObject("sample.json", &jsonObj); err != nil {
        log.Fatal(err)
    }
    log.Printf("sample.json %#v", jsonObj)
    notify, err := mc.Notify()
    if err != nil {
        log.Fatal(err)
    }

    go func() {
        for {
            ev := <-notify
            var obj sample
            if err := mc.LoadObject(ev.Name, &obj); err != nil {
                log.Fatal(err)
            }
            log.Printf("%s has been changed to %#v", ev.Name, obj)
        }
    }()

    log.Print("please edit file then you can see the notify in 1 minutes")
    time.Sleep(1 * time.Minute)
}
```


