package cron

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	config "self/configs"
	"self/internal/model/resource"
	components2 "self/internal/pkg/components"
	dlog2 "self/internal/pkg/dlog"
	"strings"
	"sync"
	"time"
)

const Length = 500

type IpsData struct {
	ErrNo  int                              `json:"errNo"`
	ErrStr string                           `json:"errStr"`
	Data   []*resource.IpNetSiteProductLine `json:"data"`
}

func ipInit() (*sync.Map, error) {
	ips, err := (&resource.IpNetSiteProductLine{}).GetsAll()
	if err != nil {
		dlog2.Errorf("%+v", err)
		return nil, err
	}
	if len(ips) > 0 {
		for i := range ips {
			IpData.Store(ips[i].Ip, ips[i])
		}
	}
	return &IpData, err
}

// 获取所有存活状态的IP
func getAllAliveIp() (data []*resource.IpNetSiteProductLine, err error) {
	cf := config.Config.Ipm
	domain := cf.Domain
	uri := cf.Apis["getIp"].Uri
	timeout := cf.Apis["getIp"].TimeoutSec
	url := fmt.Sprintf("%s%s?status=%s", domain, uri, "alive")
	resp, err := components2.Nhttp{}.DoHttpGet(url, nil, timeout)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	var r *IpsData
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	if r.ErrNo != 0 {
		return nil, errors.New(r.ErrStr)
	}
	if len(r.Data) == 0 {
		return nil, errors.WithMessagef(notFound, "url:%s", url)
	}
	return r.Data, nil
}

// 调用云平台/api/product/server接口

type ProductLine struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		Ip    string `json:"ip"`
		CpsId struct {
			Show string `json:"show"`
		} `json:"cps_id"`
	} `json:"data"`
}

func doGetIpProduct(source, secret, domain, uri, ips string, timeout int) (data *ProductLine, err error) {
	url := fmt.Sprintf("%s%s?ip_list=%s", domain, uri, ips)
	resp, err := doGetData(source, secret, url, timeout)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(resp, &data); err != nil {
		return nil, err
	}
	if data.Message != "ok" {
		return nil, errors.New(data.Message)
	}
	if len(data.Data) == 0 {
		return nil, notFound
	}
	return
}

func GetIpProduct(ips []string) (data map[string]string, err error) {
	cf := config.Config.SelfCloud
	domain := cf.Domain
	source := cf.Source
	secret := cf.Secret
	uri := cf.Apis["getProductLine"].Uri
	timeout := cf.Apis["getProductLine"].TimeoutSec
	data = make(map[string]string, 0)
	for i := 0; i < len(ips); i += Length {
		time.Sleep(500 * time.Microsecond)
		if i+Length > len(ips) {
			d, err := doGetIpProduct(source, secret, domain, uri, strings.Join(ips[i:], ","), timeout)
			if err != nil {
				if err == notFound {
					continue
				}
				return nil, err
			}
			for k := range d.Data {
				data[d.Data[k].Ip] = d.Data[k].CpsId.Show
			}
		} else {
			d, err := doGetIpProduct(source, secret, domain, uri, strings.Join(ips[i:i+Length], ","), timeout)
			if err != nil {
				if err == notFound {
					continue
				}
				return nil, err
			}
			for k := range d.Data {
				data[d.Data[k].Ip] = d.Data[k].CpsId.Show
			}
		}
	}
	return
}

// 获取IPM所有IP对应的产品树节点、网段、机房

func GetIpRelProductLine() {
	s := time.Now()
	ips, err := getAllAliveIp()
	if err != nil {
		dlog2.Errorf("%+v", err)
		return
	}
	var ip = make([]string, 0)
	var ipr = make([]*resource.IpNetSiteProductLine, 0)
	var ipsMap = make(map[string]*resource.IpNetSiteProductLine)
	for k := range ips {
		if _, ok := ipsMap[ips[k].Ip]; !ok {
			ipsMap[ips[k].Ip] = ips[k]
			ip = append(ip, ips[k].Ip)
			ipr = append(ipr, ips[k])
		}
	}
	pro, err := GetIpProduct(ip)
	if err != nil {
		dlog2.Errorf("%+v", err)
		return
	}

	for i := range ipr {
		if line, ok := pro[ipr[i].Ip]; ok {
			ipr[i].ProductLine = line
			IpData.Store(ipr[i].Ip, ipr[i])
			//fmt.Println(ips[i].Ip, "product", ips[i].ProductLine, "subnet", ips[i].Network, "site", ips[i].Site)
		}
	}

	IpData.Range(func(key, value interface{}) bool {
		if _, ok := ipsMap[key.(string)]; !ok {
			IpData.Delete(key)
		}
		return true
	})
	if UpdateCheck {
		err = (&resource.IpNetSiteProductLine{}).RemoveAndInsert(ipr)
		if err != nil {
			dlog2.Errorf("%+v", err)
		}
	}
	dlog2.Infof("ip cache done [%f]", time.Now().Sub(s).Seconds())
}
