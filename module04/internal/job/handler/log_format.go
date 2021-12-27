package kafka

import (
	"runtime"
	"self/internal/model/resource"
	dlog2 "self/internal/pkg/dlog"
	"self/internal/task"
	"strings"
)

type LogFormat interface {
	Format(s []byte) (data *WriteData, err error)
}

type WriteData struct {
	SrcIp        string `json:"src_ip"`
	DstIp        string `json:"dst_ip"`
	SrcPort      string `json:"src_port"`
	DstPort      string `json:"dst_port"`
	HttpsDomain  string `json:"https_domain"`
	HttpDomain   string `json:"http_domain"`
	InBytes      int    `json:"in_bytes"`
	OutputValue  string `json:"output_value"`
	InputValue   string `json:"input_value"`
	AgentIp      string `json:"agent_ip"`
	SamplingRate string `json:"sampling_rate"`
	LogType      string `json:"log_type"`
	FlowSeqNum   string `json:"flow_seq_num"`
	TcpFlags     string `json:"tcp_flags"`
	HttpRespCode string `json:"http_resp_code"`
	Ttl          string `json:"ttl"`
	SslName      string `json:"ssl_name"`
	RequestUrl   string `json:"request_url"`
	AgentType    string `json:"agent_type"`
	HttpMethod   string `json:"http_method"`
	DscpName     string `json:"dscp_name"`
	TimeSet      string `json:"time_set"` // 从这里开始，都是处理生成的字段
	Flow         string `json:"flow"`
	SrcIpDstIp   string `json:"src_ip_dst_ip"`
	NetWorkLine  string `json:"netWorkLine"`
	SourceNode1  string `json:"sourceNode1"`
	SourceNode2  string `json:"sourceNode2"`
	SourceNode3  string `json:"sourceNode3"`
	SourceNode4  string `json:"sourceNode4"`
	SourceNode5  string `json:"sourceNode5"`
	DestNode1    string `json:"destNode1"`
	DestNode2    string `json:"destNode2"`
	DestNode3    string `json:"destNode3"`
	DestNode4    string `json:"destNode4"`
	DestNode5    string `json:"destNode5"`
	DestAPPId    string `json:"destAPPId"`
	SourceAPPId  string `json:"sourceAPPId"`
	SourceSite   string `json:"sourceSite"`
	DestSite     string `json:"destSite"`
}

func (this *WriteData) LogFormat(s []byte, do LogFormat) (data *WriteData, err error) {
	defer func() {
		if err := recover(); err != nil {
			dlog2.Error(err)
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			dlog2.Error("Err Buf:", string(buf))
		}
	}()
	data, err = do.Format(s)
	if err != nil {
		return nil, err
	}
	// IP关联产品线、机房、网段
	sa, ok := cron.IpData.Load(data.SrcIp)
	if ok {
		dd := sa.(*resource.IpNetSiteProductLine)
		pl := dd.ProductLine
		site := dd.Site
		if len(pl) > 0 {
			nodes := make([]string, 5)
			sp := strings.Split(pl, "/")
			sp = sp[1:]
			if len(sp) > 5 {
				sp = sp[:5]
			}
			for i := range sp {
				nodes[i] = sp[i]
			}
			data.SourceNode1 = nodes[0]
			data.SourceNode2 = nodes[1]
			data.SourceNode3 = nodes[2]
			data.SourceNode4 = nodes[3]
			data.SourceNode5 = nodes[4]
			//fmt.Println("source",data.SourceNode1,data.SourceNode2,data.SourceNode3,data.SourceNode4,data.SourceNode5)
		}
		if len(site) > 0 {
			data.SourceSite = site
		}
	}
	da, o := cron.IpData.Load(data.DstIp)
	if o {
		dd := da.(*resource.IpNetSiteProductLine)
		pl := dd.ProductLine
		site := dd.Site
		if len(pl) > 0 {
			nodes := make([]string, 5)
			sp := strings.Split(pl, "/")
			sp = sp[1:]
			if len(sp) > 5 {
				sp = sp[:5]
			}
			for i := range sp {
				nodes[i] = sp[i]
			}
			data.DestNode1 = nodes[0]
			data.DestNode2 = nodes[1]
			data.DestNode3 = nodes[2]
			data.DestNode4 = nodes[3]
			data.DestNode5 = nodes[4]
			//fmt.Println("dest",data.DestNode1,data.DestNode2,data.DestNode3,data.DestNode4,data.DestNode5)
		}
		if len(site) > 0 {
			data.DestSite = site
		}
	}
	// 源端口关联应用ID
	sp, ok := cron.AppIdData.Load(data.SrcPort)
	if ok {
		data.SourceAPPId = sp.(string)
	}
	// 目的端口/域名/https域名关联应用ID
	var dt string
	switch {
	case len(data.HttpDomain) > 0:
		dt = data.HttpDomain
	case len(data.HttpsDomain) > 0:
		dt = data.HttpsDomain
	default:
		dt = data.DstPort
	}
	dp, ok := cron.AppIdData.Load(dt)
	if ok {
		data.DestAPPId = dp.(string)
	}
	// 拼接 flow netWorkLine
	data.Flow = data.SrcIp + ":" + data.SrcPort + "_" + data.DstIp + ":" + data.DstPort
	data.SrcIpDstIp = data.SrcIp + "_" + data.DstIp
	data.NetWorkLine = data.SourceSite + "-" + data.DestSite
	//fmt.Println("########################")
	//fmt.Printf("%+v\n", r)
	return
}
