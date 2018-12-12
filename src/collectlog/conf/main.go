package conf

import (
	. "collectlog/logconfig"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego/config"
	"go.etcd.io/etcd/clientv3"
)

var AppConf AppConfig

func InitConf() {
	AppConf.GenAppConf()
	// AppConf.InsertConf()
	AppConf.GetCollect()
	go AppConf.WatchEtcd()
}

type Collects struct {
	Logpath string
	Topic   string
}

type AppConfig struct {
	Logpath        string
	Loglevel       string
	ChannelSize    int
	KafkaConn      string
	KafkaParallel  int
	EtcdConns      []string
	EtcdPrefix     string
	Collects       []Collects
	CollectChannel chan []Collects
	ips            []string
}

func (acf *AppConfig) GenAppConf() {
	conf, err := config.NewConfig("ini", "config.ini")
	if err != nil {
		fmt.Println("new config failed, err:", err)
		return
	}

	logPath := conf.String("logs::logpath")
	logLevel := conf.String("logs::loglevel")
	acf.Logpath = logPath
	acf.Loglevel = logLevel
	InitLog(logPath, logLevel)
	// time.Sleep(time.Second * 5)
	channelSize := conf.DefaultInt("parallel::channelsize", 100)
	kafkaConn := conf.String("kafka::conn")
	kafkaParallel := conf.DefaultInt("kafka::parallel", 1)

	etcdConns := conf.Strings("etcd::conn")
	etcdPrefix := conf.String("etcd::prefix")
	if !strings.HasSuffix(etcdPrefix, "/") {
		etcdPrefix += "/"
	}

	acf.ChannelSize = channelSize
	acf.KafkaConn = kafkaConn
	acf.KafkaParallel = kafkaParallel
	acf.EtcdConns = etcdConns
	acf.EtcdPrefix = etcdPrefix

	acf.CollectChannel = make(chan []Collects, 1)

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		logs.Error("get local ip err:", err)
		os.Exit(1)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				acf.ips = append(acf.ips, ipnet.IP.String())
			}
		}
	}

	logs.Info("load config file successfull ! ")
	logs.Info(acf)

	return
}

func (acf *AppConfig) GetCollect() {

	fmt.Println(AppConf.EtcdConns)
	// InsertConf(ips)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   AppConf.EtcdConns,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error("New clientv3 error : ", err)
		return
	}
	defer cli.Close()
	for _, ip := range acf.ips {
		prefix := fmt.Sprintf("%s%s", AppConf.EtcdPrefix, ip)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := cli.Get(ctx, prefix)
		cancel()
		if err != nil {
			logs.Error("get failed, err:", err)
			return
		}
		for _, ev := range resp.Kvs {
			var collects []Collects
			logs.Debug("get from etc conf", string(ev.Value))
			err = json.Unmarshal(ev.Value, &collects)
			if err != nil {
				logs.Error("unformat log collect error : ", err)
				continue
			}
			acf.Collects = collects
		}
	}

}

func (acf *AppConfig) InsertConf() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   AppConf.EtcdConns,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error("New clientv3 error : ", err)
		return
	}
	defer cli.Close()

	for i := 0; i < 10000; i++ {
		// var c1 Collects = Collects{Logpath: "/aaa", Topic: fmt.Sprintf("aaa%d", i)}
		// var c2 Collects = Collects{Logpath: "/bbb", Topic: fmt.Sprintf("bbb%d", i)}

		var c1 Collects = Collects{Logpath: "/root/go/a.log", Topic: "aaa"}
		var c2 Collects = Collects{Logpath: "/root/go/b.log", Topic: "bbb"}

		var collects []Collects
		collects = append(collects, c1)
		collects = append(collects, c2)

		data, err := json.Marshal(collects)
		if err != nil {
			fmt.Println("json format error : ", err)
		}

		fmt.Println("data", string(data))
		for _, v := range acf.ips {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			prefix := fmt.Sprintf("%s%s", AppConf.EtcdPrefix, v)
			resp, err := cli.Put(ctx, prefix, string(data))
			cancel()
			if err != nil {
				logs.Error("put values in etcd error : ", err)
				return
			}
			logs.Info(resp)

		}
		time.Sleep(time.Second * 5)

	}

}

func (acf *AppConfig) WatchEtcd() {

	// InsertConf(ips)
	logs.Info("begin to monitor etcd config change")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   AppConf.EtcdConns,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error("New clientv3 error : ", err)
		return
	}
	defer cli.Close()

	if err != nil {
		logs.Error("get failed, err:", err)
		return
	}

	for _, ip := range AppConf.ips {
		prefix := fmt.Sprintf("%s%s", AppConf.EtcdPrefix, ip)
		for {
			rch := cli.Watch(context.Background(), prefix)
			for wresp := range rch {
				for _, ev := range wresp.Events {
					fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
					var collect []Collects
					err := json.Unmarshal(ev.Kv.Value, &collect)
					if err != nil {
						fmt.Println("unformat json etcd config failed : ", err)
						continue
					}
					AppConf.CollectChannel <- collect
				}
			}
		}
	}

	return
}
