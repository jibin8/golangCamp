package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"runtime"
	dlog2 "self/internal/pkg/dlog"
	"sync/atomic"
)

const (
	method    = "latest"
	timeoutMs = 6000
	family    = "v4"
)

type ReadKafka struct {
	Topic    string `json:"topic"`
	Servers  string `json:"servers"`
	GroupId  string `json:"group_id"`
	TaskName string `json:"task_name"`
	Num      int    `json:"num"`
	DealLog  LogFormat
}

func (this *ReadKafka) GetTopic() string {
	return this.Topic
}

func (this *ReadKafka) GetNum() int {
	return this.Num
}

func (this *ReadKafka) SetTaskName(name string) {
	this.TaskName = name
}

func (this *ReadKafka) ReadData(v chan []byte, SuccessDealCount int64, ErrorDealCount int64, c chan int) (err error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     this.Servers,
		"broker.address.family": family,
		"group.id":              this.GroupId,
		"session.timeout.ms":    timeoutMs,
		"auto.offset.reset":     method})
	if err != nil {
		dlog2.Error(err.Error())
		return errors.New(err.Error())
	}
	err = consumer.SubscribeTopics([]string{this.Topic}, nil)
	if err != nil {
		return errors.New(err.Error())
	}
	defer func() {
		err := consumer.Close()
		if err != nil {
			dlog2.Errorf("Kafka [%s] topic [%s] Close Error [%s]", this.Servers, this.Topic, err)
			return
		}
		if err := recover(); err != nil {
			dlog2.Errorf("Kafka [%s] topic [%s] Error [%s]", this.Servers, this.Topic, err)
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			dlog2.Error("Err Buf:", string(buf))
		}
	}()
	for {
		select {
		case <-c:
			dlog2.Infof("Kafka [%s] topic [%s] Reader Quit", this.Servers, this.Topic)
			return
		default:
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				// 数据重组
				vv := this.DealData(e.Value)
				if vv != nil {
					v <- vv
				}
				atomic.AddInt64(&SuccessDealCount, 1)
			case kafka.Error:
				atomic.AddInt64(&ErrorDealCount, 1)
				if e.Code() == kafka.ErrAllBrokersDown {
					dlog2.Errorf("Kafka [%s] topic [%s] Error", this.Servers, this.Topic)
					err = errors.New("Kafka Error")
					return err
				}
			default:
				//common.Println(this.TaskName, "Ignored......")
			}
		}
	}
}

// 具体的日志处理逻辑，不同的日志来源，不同的处理逻辑，处理成统一的日志格式

func (this *ReadKafka) DealData(s []byte) (data []byte) {
	d, err := (&WriteData{}).LogFormat(s, this.DealLog)
	if err != nil {
		dlog2.Errorf("%+v", err)
		return nil
	}
	data, err = jsoniter.Marshal(d)
	dlog2.Debug(string(data))
	if err != nil {
		dlog2.Errorf("Json Marshal Err:%s", err.Error())
		return nil
	}
	return
}
