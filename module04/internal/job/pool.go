package job

import (
	"fmt"
	dlog2 "self/internal/pkg/dlog"
	"sync"
	"time"
)

/*
数据分发池，对应向一个服务地址发送数据的结构
一个池中维护一个工作队列，包括读取数据的任务和发送数据的任务
*/

// 读取kafka的接口配置

type ReadKafka interface {
	/*
	 * 从kafka中读取数据
	 * topic:要读取的数据的topic
	 * v:数据读取后需要写入的channel
	 * c:任务中断通知channel,当需要中断任务时，c <-
	 */
	ReadData(v chan []byte, SuccessDealCount int64, ErrorDealCount int64, c chan int) error

	/*
		日志格式统一处理
	*/
	DealData(s []byte) []byte

	/*
	 * 设置任务名字，供打印日志使用
	 */
	SetTaskName(name string)
	GetNum() int
	GetTopic() string
}

// 发送kafka的接口配置

type IDispatcher interface {
	/*
	 * 将数据写入分发路径中
	 * topic:要写入的数据的topic
	 * v:数据来源的channel
	 * c:任务中断通知channel,当需要中断任务时，c <-
	 */
	WriteData(v chan []byte, SuccessDealCount int64, ErrorDealCount int64, c chan int) error
	SetTaskName(name string)
	GetHash() string
	GetNum() int
	GetBufCache() chan []byte
	GetTopic() string
}

//数据分发池的定义

type SDispatchPool struct {
	CloseFlag  chan int       //结束标记通知
	Wg         sync.WaitGroup //同步所有任务执行结束
	TaskChan   chan ITask     //池中的所有任务队列
	ReadKafka  ReadKafka      //分发池的读取接口
	Dispatcher IDispatcher    //分发池的分发接口
}

func (this *SDispatchPool) Run() {
	this.CloseFlag = make(chan int, 1)
	this.TaskChan = make(chan ITask, 100)
	bufCache := this.Dispatcher.GetBufCache()
	writeNum := this.Dispatcher.GetNum()
	readNum := this.ReadKafka.GetNum()

	if writeNum > 0 {
		for k := 0; k < this.ReadKafka.GetNum(); k++ {
			task := &ReadChanTask{
				Topic:            this.ReadKafka.GetTopic(),
				BufCache:         bufCache,
				SuccessDealCount: 0,
				ErrorDealCount:   0,
				Fun:              this.ReadKafka.ReadData,
				TaskName:         fmt.Sprintf("%s ReadTask %d/%d", this.ReadKafka.GetTopic(), k+1, readNum),
			}
			this.ReadKafka.SetTaskName(task.TaskName)
			select {
			case this.TaskChan <- task:
			case <-time.After(time.Second * 1):
				dlog2.Infof("Task timeout...", task.Topic)
			}
		}
		for k := 0; k < writeNum; k++ {
			task := &WriteChanTask{
				Topic:            this.Dispatcher.GetTopic(),
				BufCache:         bufCache,
				Fun:              this.Dispatcher.WriteData,
				SuccessDealCount: 0,
				ErrorDealCount:   0,
				TaskName:         fmt.Sprintf("%s WriteTask %d/%d", this.Dispatcher.GetTopic(), k+1, writeNum),
			}
			this.Dispatcher.SetTaskName(task.TaskName)
			select {
			case this.TaskChan <- task:
			case <-time.After(time.Second * 1):
				dlog2.Infof("Task timeout...", task.Topic)
			}
		}
	}

	for task := range this.TaskChan {
		this.Wg.Add(1)
		dlog2.Infof("New Task:%s", task.GetTaskName())
		t := task
		go func() {
			defer this.Wg.Done()
			t.Execute(this.TaskChan, &this.Wg, this.CloseFlag)
		}()
		//go task.Execute(this.TaskChan, &this.Wg, this.CloseFlag)
	}
	this.Wg.Wait()
}

func (this *SDispatchPool) Stop() {
	dlog2.Info("Begin stop")
	close(this.TaskChan)
	close(this.CloseFlag)
}
