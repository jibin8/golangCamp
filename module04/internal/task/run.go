package cron

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/google/wire"
	"github.com/niean/cron"
	"github.com/pkg/errors"
	config "self/configs"
	components2 "self/internal/pkg/components"
	dlog2 "self/internal/pkg/dlog"
	"sync"
	"time"
)

var (
	UpdateCheck bool
	IpData      sync.Map
	AppIdData   sync.Map
	notFound    = errors.New("not found")
)

func Init() {
	DbUpdateCheck()
	ipInit()
	appIdInit()
}

func DbUpdateCheck() bool {
	cronServer := config.Config.Cron.Server
	host, _ := components2.ExternalIP()
	if cronServer == host {
		UpdateCheck = true
	}
	return UpdateCheck
}

func doGetData(source, secret, url string, timeout int) (data []byte, err error) {
	now := fmt.Sprintf("%d", time.Now().Unix())
	s := fmt.Sprintf("%s%s%s", source, now, secret)
	h := md5.New()
	h.Write([]byte(s))
	signature := hex.EncodeToString(h.Sum(nil))
	headers := map[string]string{
		"Content-Type": "application/json",
		"timestamp":    now,
		"source":       source,
		"signature":    signature,
	}

	data, err = components2.Nhttp{}.DoHttpGet(url, headers, timeout)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return
}

func task() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		GetIpRelProductLine()
	}()
	go func() {
		defer wg.Done()
		GetAppIdInfo()
	}()
	wg.Wait()
}

func Run() {
	conf := config.Config.Cron
	c := conf.Cache
	cronSecond := conf.CacheSecond
	if c {
		dlog2.Info("cache update cron run")
		c := cron.New()
		c.AddFuncCC(cronSecond, task, 1)
		c.Start()
	}
}

var ProviderSet = wire.NewSet(Init)
