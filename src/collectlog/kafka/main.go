package kafka

import (
	. "collectlog/conf"
	. "collectlog/tailf"

	"github.com/astaxie/beego/logs"

	"github.com/Shopify/sarama"
)

var (
	KafkaMgr *KafkaManager = &KafkaManager{}
)

type KafkaManager struct {
	clients []sarama.SyncProducer
}

func (kmr *KafkaManager) InitKafkaMgr() (err error) {
	logs.NewLogger()

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	var client sarama.SyncProducer
	for i := 0; i < AppConf.KafkaParallel; i++ {
		client, err = sarama.NewSyncProducer([]string{AppConf.KafkaConn}, config)
		if err != nil {
			logs.Error("producer close, err:", err)
			return
		}

		kmr.clients = append(kmr.clients, client)
	}
	return
}

func (kmr *KafkaManager) processMsg() {
	var count int
	for msg := range TailMgr.MsgChannel {
		go func() {
			if count >= AppConf.KafkaParallel {
				count = 0
			}
			kafkaMsg := &sarama.ProducerMessage{}
			kafkaMsg.Topic = msg.Topic
			kafkaMsg.Value = sarama.StringEncoder(msg.Msg)
			logs.Debug("send msg ", msg.Topic, msg.Msg)
			pid, offset, err := kmr.clients[count].SendMessage(kafkaMsg)
			if err != nil {
				logs.Error("send message failed,", err)
				return
			}

			logs.Debug("pid:%v offset:%v\n", pid, offset)
		}()
		count += 1
	}
}

func InitKafka() {

	KafkaMgr.InitKafkaMgr()

	// defer client.Close()

	logs.Info("kafka  begin to receive msg ")
}

// 	for msg := range TailMgr.MsgChannel {
// 		go func() {
// 			kafkaMsg := &sarama.ProducerMessage{}
// 			kafkaMsg.Topic = msg.Topic
// 			kafkaMsg.Value = sarama.StringEncoder(msg.Msg)
// 			logs.Debug("send msg ", msg.Topic, msg.Msg)
// 			pid, offset, err := client.SendMessage(kafkaMsg)
// 			if err != nil {
// 				logs.Error("send message failed,", err)
// 				return
// 			}

// 			logs.Debug("pid:%v offset:%v\n", pid, offset)
// 		}()

// 	}
// }

// for i := 0; i < 100; i++ {
// 	msg := &sarama.ProducerMessage{}
// 	msg.Topic = "nginx_log"
// 	msg.Value = sarama.StringEncoder(fmt.Sprintf("test %d", i))

// 	client, err := sarama.NewSyncProducer([]string{"10.12.9.195:9092"}, config)
// 	if err != nil {
// 		fmt.Println("producer close, err:", err)
// 		return
// 	}

// 	defer client.Close()

// 	pid, offset, err := client.SendMessage(msg)
// 	if err != nil {
// 		fmt.Println("send message failed,", err)
// 		return
// 	}

// 	fmt.Printf("pid:%v offset:%v\n", pid, offset)
// 	time.Sleep(time.Second)
// }
