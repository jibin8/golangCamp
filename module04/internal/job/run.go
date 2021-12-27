package job

import (
	"os"
	dlog2 "self/internal/pkg/dlog"
)

var KafkaConf = []map[string]interface{}{
	{
		"read": map[string]interface{}{
			"dispatch_type": "gigamon",
			"server":        "127.0.0.1:9092",
			"topic":         "sys_test",
			"group_id":      "read_kafka_01",
			"read_num":      3,
		},
		"write": map[string]interface{}{
			"server":    "127.0.0.1:9093",
			"topic":     "device-syslog",
			"write_num": 1,
			"buf_cache": make(chan []byte, 1000),
		},
	},
	{"read": map[string]interface{}{
		"dispatch_type": "sflow",
		"server":        "127.0.0.1:9094",
		"topic":         "device-sflow",
		"group_id":      "read_sflow_01",
		"read_num":      3,
	}, "write": map[string]interface{}{
		"server":    "127.0.0.1:9095",
		"topic":     "device-syslog",
		"write_num": 1,
		"buf_cache": make(chan []byte, 1000),
	}},
	{"read": map[string]interface{}{
		"dispatch_type": "netflow",
		"server":        "127.0.0.1:9096",
		"topic":         "device-netflow",
		"group_id":      "netflow_01",
		"read_num":      3,
	}, "write": map[string]interface{}{
		"server":    "127.0.0.1:9092",
		"topic":     "device-syslog",
		"write_num": 1,
		"buf_cache": make(chan []byte, 1000),
	}},
}

var (
	DispatchPools map[string]*SDispatchPool
	sig           chan os.Signal
)

func init() {
	dlog2.Infof("job init")
	DispatchPools = make(map[string]*SDispatchPool)
}

func Run() {
	ConfigDispatchers := make(map[string]*SDispatchPool)
	for i := range KafkaConf {
		//workingPool := new(SDispatchPool)
		conf := KafkaConf[i]
		read := conf["read"].(map[string]interface{})
		write := conf["write"].(map[string]interface{})
		dispatchType := read["dispatch_type"].(string)
		workingPool := selectDispatch(dispatchType, read, write)
		if workingPool == nil {
			continue
		}
		hash := workingPool.Dispatcher.GetHash()
		ConfigDispatchers[hash] = workingPool
		if _, ok := DispatchPools[hash]; ok {
			continue
		}
		DispatchPools[hash] = workingPool
		go workingPool.Run()
	}
	for key := range DispatchPools {
		if _, ok := ConfigDispatchers[key]; !ok {
			DispatchPools[key].Stop()
			dlog2.Infof("Stop Pool:%s", key)
			delete(DispatchPools, key)
		}
	}

}
