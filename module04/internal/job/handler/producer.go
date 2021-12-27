package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	"runtime"
	dlog2 "self/internal/pkg/dlog"
	"sync/atomic"
)

type WriteKafka struct {
	Type     string      `json:"type"`
	Servers  string      `json:"servers"`
	Topic    string      `json:"topic"`
	Num      int         `json:"num"`
	TaskName string      `json:"task_name"`
	BufCache chan []byte `json:"buf_cache"`
}

func (this *WriteKafka) GetTopic() string {
	return this.Topic
}

func (this *WriteKafka) GetBufCache() chan []byte {
	return this.BufCache
}

func (this *WriteKafka) GetNum() int {
	return this.Num
}

func (this *WriteKafka) SetTaskName(name string) {
	this.TaskName = name
}

func (this *WriteKafka) WriteData(v chan []byte, SuccessDealCount int64, ErrorDealCount int64, c chan int) (err error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": this.Servers})
	if err != nil {
		dlog2.Errorf("kafka [%s] produce Error [%s]", this.Servers, err.Error())
		return errors.Wrapf(err, "kafka [%s] produce Error", this.Servers)
	}
	deliveryChan := make(chan kafka.Event)
	a := make(chan int)
	defer func() {
		p.Close()
		close(deliveryChan)
		if err := recover(); err != nil {
			dlog2.Errorf("WriteKafka Topic Error:%s", this.Topic, err)
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			dlog2.Error("Err Buf:", string(buf))
		}
	}()

	go func(this *WriteKafka, SuccessDealCount, ErrorDealCount int64, a chan int, deliveryChan chan kafka.Event) {
		for {
			select {
			case <-a:
				dlog2.Infof("Kafka Callback Topic:%s Writer Quit", this.Topic)
				return
			case e := <-deliveryChan:
				m := e.(*kafka.Message)
				if m.TopicPartition.Error != nil {
					dlog2.Errorf("Delivery failed: %v", m.TopicPartition.Error)
					atomic.AddInt64(&ErrorDealCount, 1)
				} else {
					atomic.AddInt64(&SuccessDealCount, 1)
					dlog2.Infof("Delivered message to topic %s [%d] at offset %v", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
			}
		}
	}(this, SuccessDealCount, ErrorDealCount, a, deliveryChan)
	for {
		select {
		case <-c:
			dlog2.Infof("topic [%s] Writer Quit....", this.Topic)
			a <- 0
			return
		case data := <-v:
			err = p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &this.Topic, Partition: kafka.PartitionAny},
				Value:          data}, deliveryChan)
			if err != nil {
				dlog2.Errorf("kafka [%s] Topic [%s] produce Error [%s]", this.Servers, this.Topic, err.Error())
				atomic.AddInt64(&ErrorDealCount, 1)
				continue
			}
		}
	}
}

func (this *WriteKafka) GetHash() string {
	return fmt.Sprintf("type:%s|broker:%s|topic:%s", this.Type, this.Servers, this.Topic)
}
