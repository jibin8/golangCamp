package job

import (
	kafka2 "self/internal/job/handler"
	dlog2 "self/internal/pkg/dlog"
)

func selectDispatch(dispatchType string, read, write map[string]interface{}) (workingPool *SDispatchPool) {
	bufCache := write["buf_cache"].(chan []byte)
	workingPool = new(SDispatchPool)

	dispatcher := &kafka2.WriteKafka{
		Type:     dispatchType,
		Servers:  write["server"].(string),
		Topic:    write["topic"].(string),
		Num:      write["write_num"].(int),
		BufCache: bufCache,
	}
	workingPool.Dispatcher = dispatcher

	readKafka := &kafka2.ReadKafka{
		Servers: read["server"].(string),
		Topic:   read["topic"].(string),
		GroupId: read["group_id"].(string),
		Num:     read["read_num"].(int),
	}
	switch dispatchType {
	case "gigamon":
		readKafka.DealLog = &kafka2.GigamonData{}
		workingPool.ReadKafka = readKafka
	case "sflow":
		readKafka.DealLog = &kafka2.SFlowData{}
		workingPool.ReadKafka = readKafka
	case "netflow":
		readKafka.DealLog = &kafka2.NetFlowData{}
		workingPool.ReadKafka = readKafka
	default:
		dlog2.Errorf("Unknown dispatcher type:%s", dispatchType)
		return nil
	}
	return
}
