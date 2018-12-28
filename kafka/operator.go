package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"log"
	"logserver/common"
	"logserver/configs"
	"logserver/es"
	"logserver/etcd"

	//"sync"
	"time"
)

func SendToKafka(datas []string, topic string) error {
	var (
		msgs []*sarama.ProducerMessage
		err  error
	)
	msgs = make([]*sarama.ProducerMessage, 0)
	for _, data := range datas {
		msgs = append(msgs, &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(data),
		})
	}
	err = Producer.SendMessages(msgs)
	if err != nil {
		logs.Error(err)
	}
	return err
}



func ConsumerFromKafka4(info *common.JobWorkInfo, lock *etcd.JobLock) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_1_0
	client, err := sarama.NewClient([]string{configs.AppConfig.KafkaAddr}, config)
	if err != nil {
		log.Fatalln(err)
	}

	offsetManager, err := sarama.NewOffsetManagerFromClient("test", client)
	if err != nil {
		log.Fatalln(err)
	}

	pids, err := client.Partitions(info.Job.Topic)
	if err != nil {
		log.Fatalln(err)
	}

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Fatalln(err)
	}

	//defer consumer.Close()

	//wg := &sync.WaitGroup{}

	for _, v := range pids {
		//wg.Add(1)
		go consume(consumer, offsetManager, v, info, lock)
	}

	//wg.Wait()
	fmt.Println("tui chu **************************")
}
func consume(c sarama.Consumer, om sarama.OffsetManager, p int32, info *common.JobWorkInfo, lock *etcd.JobLock) {
	var (
		pom    sarama.PartitionOffsetManager
		pc     sarama.PartitionConsumer
		err    error
		offset int64
	)
	defer c.Close()

	if pom, err = om.ManagePartition(info.Job.Topic, p); err != nil {
		log.Fatalln(err)
	}

	defer pom.Close()

	offset, _ = pom.NextOffset()
	if offset == -1 {
		offset = sarama.OffsetOldest
	}

	if pc, err = c.ConsumePartition(info.Job.Topic, p, offset); err != nil {
		log.Fatalln(err)
	}

	defer pc.Close()
	var (
		logDatas []interface{}
		t        *time.Timer
	)
	logDatas = make([]interface{}, 0)
	t = time.NewTimer(time.Second * 1)
	for {
		select {
		case msg := <-pc.Messages():
			logDatas = append(logDatas, string(msg.Value))
			pom.MarkOffset(msg.Offset+1, "")
			if len(logDatas) >= 100 {

				es.GelasticCli.CreateBulkDocument(common.CreateIndexByType(info.Job.Topic, info.Job.IndexType), logDatas, info.Job.Pipeline)
				common.SliceClear(&logDatas)
			}
			//log.Printf("[%v] Consumed message offset %v content is %s\n", p, msg.Offset, string(msg.Value))
		case <-t.C:
			if len(logDatas) > 0 {
				fmt.Println(logDatas)
				es.GelasticCli.CreateBulkDocument(common.CreateIndexByType(info.Job.Topic, info.Job.IndexType), logDatas, info.Job.Pipeline)
				common.SliceClear(&logDatas)
			}
			t.Reset(time.Second * 1)
		case <-info.ConText.Done():
			lock.Unlock()
			return
		}
	}
	//for msg := range pc.Messages() {
	//	log.Printf("[%v] Consumed message offset %v content is %s\n", p, msg.Offset, string(msg.Value))
	//	pom.MarkOffset(msg.Offset+1, "")
	//}
	fmt.Println("tui ************")
}
