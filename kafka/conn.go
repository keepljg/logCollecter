package kafka

import (
	"github.com/Shopify/sarama"
	"logserver/configs"
	"time"
)

var (
	Producer sarama.SyncProducer
)

// 初始化kafka异步生产者
func InitKafka() error {
	var (
		config *sarama.Config
		err    error
	)
	config = sarama.NewConfig()
	config.Version = sarama.V0_11_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll          // 等待服务器所以副本都保存成功后响应
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 随机分区 即随机向分区发送消息
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true // 是否等待服务器都成功失败响应（如果config.Producer.RequiredAcks设置为NoResponse 则不可设置)

	//异步生产者
	//if Producer, err = sarama.NewAsyncProducer([]string{configs.AppConfig.KafkaAddr}, config); err != nil{
	//	return err
	//}
	//异步生产者
	if Producer, err = sarama.NewSyncProducer([]string{configs.AppConfig.KafkaAddr}, config); err != nil {
		return err
	}
	return nil
}

// 创建kafka消费者
func NewConsumer() (sarama.Consumer, error) {
	var (
		config   *sarama.Config
		err      error
		consumer sarama.Consumer
	)
	config = sarama.NewConfig()
	config.Version = sarama.V2_0_1_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.CommitInterval = 1 * time.Second // 提交offect时间间隔 1s
	if consumer, err = sarama.NewConsumer([]string{configs.AppConfig.KafkaAddr}, config); err != nil {
		return nil, err
	}
	return consumer, nil
}
