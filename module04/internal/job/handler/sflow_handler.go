package kafka

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type SFlowData struct {
	Timestamp            string `json:"@timestamp"`
	AgentIP              string `json:"agent_ip"`
	SrcIP                string `json:"src_ip"`
	DstIP                string `json:"dst_ip"`
	DstPort              string `json:"dst_port"`
	SrcPort              string `json:"src_port"`
	FrameLength          string `json:"frame_length"`
	InputInterfaceValue  string `json:"input_interface_value"`
	OutputInterfaceValue string `json:"output_interface_value"`
	SamplingRate         string `json:"sampling_rate"`
	Type                 string `json:"type"`
	SrcPriority          string `json:"src_priority"`
}

func (this *SFlowData) Format(s []byte) (data *WriteData, err error) {
	var r SFlowData
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(s, &r)
	if err != nil {
		return nil, errors.Wrap(err, "SFlow Log Json Unmarshal Err")
	}
	// 一一对应
	d := &WriteData{
		SrcIp:        r.SrcIP,
		DstIp:        r.DstIP,
		SrcPort:      r.SrcPort,
		DstPort:      r.DstPort,
		OutputValue:  r.OutputInterfaceValue,
		InputValue:   r.InputInterfaceValue,
		AgentIp:      r.AgentIP,
		SamplingRate: r.SamplingRate,
		LogType:      r.Type,
		DscpName:     r.SrcPriority,
	}
	d.InBytes, _ = strconv.Atoi(r.FrameLength)
	// Syslog 生成日志时间
	nt, err := time.ParseInLocation(time.RFC3339, r.Timestamp, time.Local)
	if err != nil {
		return nil, errors.Wrap(err, "SFlow Log Format Time Err")
	}
	d.TimeSet = nt.Format("2006-01-02 15:04:05")
	return d, nil
}
