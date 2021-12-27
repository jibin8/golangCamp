package cron

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	config "self/configs"
	"self/internal/model/resource"
	dlog2 "self/internal/pkg/dlog"
	semaphore2 "self/internal/pkg/semaphore"
	"sync"
	"time"
)

func appIdInit() *sync.Map {
	ids, err := (&resource.DomainApp{}).GetsAll()
	if err != nil {
		dlog2.Errorf("%+v", err)
		return nil
	}
	if len(ids) > 0 {
		for i := range ids {
			AppIdData.Store(ids[i].Domain, ids[i].AppId)
		}
	}
	return &AppIdData
}

type AppId struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

func doGetAppId(source, secret, domain, uri string, timeout int) (data *AppId, err error) {
	url := fmt.Sprintf("%s%s", domain, uri)
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

type ApplicationData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		AppId             string   `json:"app_id"`
		ProductTreeName   string   `json:"product_tree_name"`
		ApplicationDomain []string `json:"application_domain"`
		ApplicationPort   []struct {
			PortName  string `json:"port_name"`
			PortValue int    `json:"port_value"`
		} `json:"application_port"`
	} `json:"data"`
}

func doGetDomainByAppId(source, secret, domain, uri, appId string, timeout int) (data *ApplicationData, err error) {
	url := fmt.Sprintf("%s%s?app_id=%s", domain, uri, appId)
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
	return
}

func GetAppIdInfo() {
	s := time.Now()
	cf := config.Config.SelfCloud
	domain := cf.Domain
	source := cf.Source
	secret := cf.Secret
	uri := cf.Apis["getAppId"].Uri
	timeout := cf.Apis["getAppId"].TimeoutSec
	appId, err := doGetAppId(source, secret, domain, uri, timeout)
	if err != nil {
		dlog2.Errorf("%+v", err)
		return
	}
	uri = cf.Apis["getDomain"].Uri
	timeout = cf.Apis["getDomain"].TimeoutSec
	var wg sync.WaitGroup
	var seam = semaphore2.NewSemaphore(5)
	var domainChan = make(chan []string, 100)
	var done = make(chan bool)
	var appIdMaps = make(map[string]string, 0)
	go func() {
		for item := range domainChan {
			if _, ok := appIdMaps[item[0]]; !ok {
				appIdMaps[item[0]] = item[1]
			}
		}
		done <- true
	}()
	for i := range appId.Data {
		ad := appId.Data[i]
		seam.Acquire()
		wg.Add(1)
		go func(appId string) {
			defer func() {
				wg.Done()
				seam.Release()
			}()
			d, err := doGetDomainByAppId(source, secret, domain, uri, appId, timeout)
			if err != nil {
				dlog2.Errorf("%+v", err)
				return
			}
			if len(d.Data.ApplicationDomain) > 0 {
				for k := range d.Data.ApplicationDomain {
					do := d.Data.ApplicationDomain[k]
					domainChan <- []string{do, appId}
					AppIdData.Store(do, appId)
				}
			}
			if len(d.Data.ApplicationPort) > 0 {
				for k := range d.Data.ApplicationPort {
					do := fmt.Sprintf("%d", d.Data.ApplicationPort[k].PortValue)
					domainChan <- []string{do, appId}
				}
			}
		}(ad)
	}
	wg.Wait()
	close(domainChan)
	<-done

	AppIdData.Range(func(key, value interface{}) bool {
		if _, ok := appIdMaps[key.(string)]; !ok {
			AppIdData.Delete(key)
		}
		return true
	})
	if UpdateCheck {
		var appIds = make([]*resource.DomainApp, len(appIdMaps))
		var count = 0
		for k, v := range appIdMaps {
			a := &resource.DomainApp{
				Domain: k,
				AppId:  v,
			}
			appIds[count] = a
			count = count + 1
		}
		err = (&resource.DomainApp{}).RemoveAndInsert(appIds)
		if err != nil {
			dlog2.Errorf("%+v", err)
		}
	}
	dlog2.Infof("appId cache done [%f]", time.Now().Sub(s).Seconds())
}
