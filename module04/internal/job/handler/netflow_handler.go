package kafka

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"time"
)

type NetFlowData struct {
	Timestamp string `json:"@timestamp"`
	Host      string `json:"host"`
	Netflow   struct {
		FlowSeqNum       int64  `json:"flow_seq_num"`
		InBytes          int    `json:"in_bytes"`
		InputSnmp        int64  `json:"input_snmp"`
		Ipv4DstAddr      string `json:"ipv4_dst_addr"`
		Ipv4SrcAddr      string `json:"ipv4_src_addr"`
		L4DstPort        int64  `json:"l4_dst_port"`
		L4SrcPort        int64  `json:"l4_src_port"`
		OutputSnmp       int64  `json:"output_snmp"`
		TcpFlags         int64  `json:"tcp_flags"`
		SamplingInterval int64  `json:"sampling_interval"`
		SrcTos           int64  `json:"src_tos"`
	} `json:"netflow"`
	Type string `json:"type"`
}

func (this *NetFlowData) Format(s []byte) (data *WriteData, err error) {
	var r NetFlowData
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(s, &r)
	if err != nil {
		return nil, errors.Wrap(err, "SFlow Log Json Unmarshal Err")
	}
	// 一一对应
	d := &WriteData{
		SrcIp:        r.Netflow.Ipv4SrcAddr,
		DstIp:        r.Netflow.Ipv4DstAddr,
		SrcPort:      fmt.Sprintf("%d", r.Netflow.L4SrcPort),
		DstPort:      fmt.Sprintf("%d", r.Netflow.L4DstPort),
		InBytes:      r.Netflow.InBytes,
		OutputValue:  fmt.Sprintf("%d", r.Netflow.OutputSnmp),
		InputValue:   fmt.Sprintf("%d", r.Netflow.InputSnmp),
		AgentIp:      r.Host,
		SamplingRate: fmt.Sprintf("%d", r.Netflow.SamplingInterval),
		LogType:      r.Type,
		FlowSeqNum:   fmt.Sprintf("%d", r.Netflow.FlowSeqNum),
		TcpFlags:     fmt.Sprintf("%d", r.Netflow.TcpFlags),
		DscpName:     fmt.Sprintf("%d", r.Netflow.SrcTos),
	}
	// Syslog 生成日志时间
	nt, err := time.ParseInLocation(time.RFC3339, r.Timestamp, time.Local)
	if err != nil {
		return nil, errors.Wrap(err, "SFlow Log Format Time Err")
	}
	d.TimeSet = nt.Format("2006-01-02 15:04:05")
	return d, nil
}
