package tailf

import (
	. "collectlog/conf"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/hpcloud/tail"
)

type Message struct {
	Msg   string
	Topic string
}

type TailMgrStruct struct {
	TailObjs   []*TailObj
	MsgChannel chan *Message
}

type TailObj struct {
	tail     *tail.Tail
	topic    string
	exitChan chan int
}

var (
	TailMgr *TailMgrStruct = &TailMgrStruct{}
)

func (tmr *TailMgrStruct) getTailObj(collect Collects) (tailObj TailObj, err error) {
	tails, err := tail.TailFile(collect.Logpath, tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		logs.Error("init %s tail obj failed error : %s", collect.Logpath, err)
		return
	}

	tailObj.topic = collect.Topic
	tailObj.tail = tails
	tailObj.exitChan = make(chan int, 1)
	return
}

func (tmr *TailMgrStruct) initTailMgr() (err error) {
	logs.Info("begin init tail Mgr")
	tmr.MsgChannel = make(chan *Message, AppConf.ChannelSize)
	var TailObjs []*TailObj
	for _, v := range AppConf.Collects {
		tailObj, err := tmr.getTailObj(v)
		if err != nil {
			logs.Error("init tailobj failed :", err)
		}
		TailObjs = append(TailObjs, &tailObj)
	}
	tmr.TailObjs = TailObjs
	logs.Info("finish init tail mgr tailobjs : ", tmr.TailObjs)
	return
}

func (tmr *TailMgrStruct) processTailObj(tailObj *TailObj) (err error) {
	var msg *tail.Line
	var ok bool
	logs.Info("begin to process tailobj : ", tailObj)
	for true {
		logs.Debug("process tail obj", tailObj.topic, tailObj.tail.Filename)
		logs.Debug("all mgr tailobjs : ", tmr.TailObjs)

		select {
		case msg, ok = <-tailObj.tail.Lines:
			if !ok {
				logs.Warning("tail file close reopen, filename:%s\n", tailObj.tail.Filename)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			m := &Message{Msg: msg.Text, Topic: tailObj.topic}

			tmr.MsgChannel <- m
			logs.Debug("send msg to kafka channel  ", m)
		case <-tailObj.exitChan:

			logs.Info("tailobjs thread exit : ", tailObj.topic, tailObj.tail.Filename)
			return
		}

	}
	logs.Info("end process tail obj")
	return
}

func (tmr *TailMgrStruct) processTailObjs() (err error) {
	logs.Info("begin process tail objs ", tmr.TailObjs)
	for _, tailObj := range tmr.TailObjs {
		logs.Info("before process tail obj", tailObj.tail.Filename, tailObj.topic)
		go tmr.processTailObj(tailObj)
	}
	logs.Info("end process tailobjs")
	return
}

func RemoveRepByMap(slc []int) []int {
	result := []int{}
	tempMap := map[int]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

func IfIntInSlice(i int, s []int) (b bool) {
	b = false
	for _, j := range s {
		if j == i {
			b = true
			return
		}
	}
	return
}

func (tmr *TailMgrStruct) checkConfigChange() {
	logs.Info("begin running check config change")
	for v := range AppConf.CollectChannel {
		logs.Info("get new config ", v)
		var indexes []int
		var etcdConfIndex []int

		for i, collect := range v {
			var flag = false
			for j, tails := range tmr.TailObjs {
				if tails.tail.Filename == collect.Logpath && tails.topic == collect.Topic {
					indexes = append(indexes, j)
					flag = true
				}
			}
			if flag == false {
				etcdConfIndex = append(etcdConfIndex, i)
			}

		}

		logs.Debug("all", tmr.TailObjs, indexes, etcdConfIndex, v)

		indexes = RemoveRepByMap(indexes)
		var tailObjs []*TailObj
		for i, tailobj := range tmr.TailObjs {
			if IfIntInSlice(i, indexes) {
				tailObjs = append(tailObjs, tailobj)
			} else {
				tailobj.exitChan <- 1
			}
		}

		tmr.TailObjs = tailObjs
		logs.Debug("etcdIndex Content", etcdConfIndex, v)
		for i, collect := range v {
			if IfIntInSlice(i, etcdConfIndex) {
				tailObj, err := tmr.getTailObj(collect)
				if err != nil {
					logs.Error("get tail obj failed ", err)
					break
				}
				tmr.TailObjs = append(tmr.TailObjs, &tailObj)

				logs.Info("start process new tailobj ", tailObj)
				go tmr.processTailObj(&tailObj)
			}
		}
		logs.Info("update etcd config finish")
	}
	logs.Info("end config check finish!")

}

func InitTail() {
	logs.Info("begin init tail")
	TailMgr.initTailMgr()
	go TailMgr.processTailObjs()
	go TailMgr.checkConfigChange()
	logs.Info("finish init tail ")
}
