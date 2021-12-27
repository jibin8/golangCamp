package kafka

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"time"
)

type GigamonData struct {
	GigamonMdataHttpRespCode            string `json:"GigamonMdataHttpRespCode"`
	GigamonMdataDeltaBytesRcvd          int    `json:"GigamonMdataDeltaBytesRcvd"`
	GigamonMdataHTTPHost                string `json:"GigamonMdataHttpHost"`
	GigamonMdataHTTPMethod              string `json:"GigamonMdataHttpMethod"`
	GigamonMdataIPTTL                   string `json:"GigamonMdataIpTtl"`
	GigamonMdataSslIssuerName           string `json:"GigamonMdataSslIssuerName"`
	GigamonMdataSslServerNameIndication string `json:"GigamonMdataSslServerNameIndication"`
	GigamonMdataTcpFlags                string `json:"GigamonMdataTcpFlags"`
	GigamonMdataUserAgentTxtInd         string `json:"GigamonMdataUserAgentTxtInd"`
	DestinationAddress                  string `json:"destinationAddress"`
	DestinationPort                     string `json:"destinationPort"`
	Host                                string `json:"host"`
	RequestURL                          string `json:"requestUrl"`
	SourceAddress                       string `json:"sourceAddress"`
	SourcePort                          string `json:"sourcePort"`
	Timestamp                           string `json:"timestamp"`
	Syslog                              string `json:"syslog"`
	Name                                string `json:"name"`
}

func (this *GigamonData) Format(s []byte) (data *WriteData, err error) {
	var r GigamonData
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(s, &r)
	if err != nil {
		return nil, errors.Wrap(err, "Giga Log Json Unmarshal Err")
	}
	// 一一对应
	d := &WriteData{
		SrcIp:        r.SourceAddress,
		DstIp:        r.DestinationAddress,
		SrcPort:      r.SourcePort,
		DstPort:      r.DestinationPort,
		HttpsDomain:  r.GigamonMdataSslServerNameIndication,
		HttpDomain:   r.GigamonMdataHTTPHost,
		InBytes:      r.GigamonMdataDeltaBytesRcvd,
		AgentIp:      r.Host,
		LogType:      r.Name,
		TcpFlags:     r.GigamonMdataTcpFlags,
		HttpRespCode: r.GigamonMdataHttpRespCode,
		Ttl:          r.GigamonMdataIPTTL,
		SslName:      r.GigamonMdataSslIssuerName,
		RequestUrl:   r.RequestURL,
		AgentType:    r.GigamonMdataUserAgentTxtInd,
		HttpMethod:   r.GigamonMdataHTTPMethod,
	}
	// Syslog 生成日志时间
	nt, err := time.ParseInLocation("Mon Jan 02 15:04:05 2006  1/1/e1", r.Syslog, time.Local)
	if err != nil {
		return nil, errors.Wrap(err, "Giga Log Format Time Err")
	}
	d.TimeSet = nt.Format("2006-01-02 15:04:05")
	return d, nil
}
