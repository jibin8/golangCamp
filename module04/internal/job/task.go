package job

import (
	"runtime"
	dlog2 "self/internal/pkg/dlog"
	"sync"
	"time"
)

type ITask interface {
	Execute(taskChan chan ITask, wg *sync.WaitGroup, CloseFlag chan int) (err error)
	GetTaskName() string
}

//读取数据，进行处理

type ReadDealFunc func(v chan []byte, SuccessDealCount int64, ErrorDealCount int64, c chan int) error

type ReadChanTask struct {
	Topic            string       //数据类型
	BufCache         chan []byte  //数据管道，指向ByteChannelConf中的BufCache
	Fun              ReadDealFunc //任务执行函数，ReadData或者WriteData
	TaskName         string       //任务名字，供日志输出使用
	SuccessDealCount int64        //任务成功处理的量
	ErrorDealCount   int64        //任务错误处理的量
}

func (this *ReadChanTask) Execute(taskChan chan ITask, wg *sync.WaitGroup, CloseFlag chan int) (err error) {
	defer func() {
		//wg.Done()
		if err := recover(); err != nil {
			dlog2.Errorf("Task:%s Err:%s", this.TaskName, err)
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			dlog2.Error("Err Buf:", string(buf))
		}
		dlog2.Errorf("Task:%s is Down", this.TaskName)
	}()

	dlog2.Infof("Execute Task:%s bufCache:%p", this.TaskName, this.BufCache)
	err = this.Fun(this.BufCache, this.SuccessDealCount, this.ErrorDealCount, CloseFlag)
	if err != nil {
		dlog2.Errorf("Task:%s Error:%s", this.TaskName, err.Error())
	}
	time.Sleep(1 * time.Second) //如果任务执行失败了，延迟1秒再开始
	taskChan <- this
	return
}

func (this *ReadChanTask) GetTaskName() string {
	return this.TaskName
}

//从chan中读取数据，进行处理

type WriteDealFunc func(v chan []byte, SuccessDealCount int64, ErrorDealCount int64, c chan int) error

type WriteChanTask struct {
	Topic            string        //数据类型
	BufCache         chan []byte   //数据管道，指向ByteChannelConf中的BufCache
	Fun              WriteDealFunc //任务执行函数，ReadData或者WriteData
	TaskName         string        //任务名字，供日志输出使用
	SuccessDealCount int64         //任务成功处理的量
	ErrorDealCount   int64         //任务错误处理的量
}

func (this *WriteChanTask) Execute(taskChan chan ITask, wg *sync.WaitGroup, CloseFlag chan int) (err error) {
	defer func() {
		//wg.Done()
		if err := recover(); err != nil {
			dlog2.Errorf("Task:%s Err:%s", this.TaskName, err)
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			dlog2.Error("Err Buf:", string(buf))
		}
		dlog2.Errorf("Task:%s is Down", this.TaskName)
	}()

	dlog2.Infof("Execute Task:%s bufCache:%p", this.TaskName, this.BufCache)
	err = this.Fun(this.BufCache, this.SuccessDealCount, this.ErrorDealCount, CloseFlag)
	if err != nil {
		dlog2.Errorf("Task:%s Error:%s", this.TaskName, err.Error())
	}
	time.Sleep(1 * time.Second) //如果任务执行失败了，延迟1秒再开始
	taskChan <- this
	return
}

func (this *WriteChanTask) GetTaskName() string {
	return this.TaskName
}
