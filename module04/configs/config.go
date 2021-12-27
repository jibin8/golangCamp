package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	dlog2 "self/internal/pkg/dlog"
)

// 以下配置, 来自服务的标准配置文件

type Conf struct {
	Http      HttpConfig      `toml:"http" json:"http"`
	Proc      ProcConfig      `toml:"proc" json:"proc"`
	Runtime   RuntimeConfig   `toml:"runtime" json:"runtime"`
	Dlog      dlog2.LogConfig `toml:"dlog" json:"dlog"`
	OauthWeb  OauthConfig     `toml:"oauth" json:"oauth"`
	Database  DataBaseConfig  `toml:"database" json:"database"`
	Region    RegionConfig    `toml:"region" json:"region"`
	Cron      CronConfig      `toml:"cron" json:"cron"`
	Ipm       IpmConfig       `toml:"ipm" json:"ipm"`
	SelfCloud SelfCloudConfig `toml:"selfCloud" json:"selfCloud"`
}

type HttpConfig struct {
	Disable         bool   `toml:"disable" json:"disable"`
	Listen          string `toml:"listen" json:"listen"`
	Mode            string `toml:"mode" json:"mode"`
	AllTimeoutSec   int    `toml:"allTimeoutSec" json:"allTimeoutSec"`
	ReadTimeoutSec  int    `toml:"readTimeoutSec" json:"readTimeoutSec"`
	WriteTimeoutSec int    `toml:"writeTimeoutSec" json:"writeTimeoutSec"`
	ExitTimeoutSec  int    `toml:"exitTimeoutSec" json:"exitTimeoutSec"`
}

type OauthConfig struct {
	SystemName string               `toml:"systemName" json:"systemName"`
	Source     string               `toml:"source" json:"source"`
	Signature  string               `toml:"signature" json:"signature"`
	SelfDomain string               `toml:"selfDomain" json:"selfDomain"`
	Apis       map[string]ApiConfig `toml:"apis" json:"apis"`
}

type IpmConfig struct {
	Domain string               `toml:"domain" json:"domain"`
	Apis   map[string]ApiConfig `toml:"apis" json:"apis"`
}

type SelfCloudConfig struct {
	Domain string               `toml:"domain" json:"domain"`
	Source string               `toml:"source" json:"source"`
	Secret string               `toml:"secret" json:"secret"`
	Apis   map[string]ApiConfig `toml:"apis" json:"apis"`
}

type ApiConfig struct {
	Domain     string `toml:"domain" json:"domain"`
	Uri        string `toml:"uri" json:"uri"`
	TimeoutSec int    `toml:"timeoutSec" json:"timeoutSec"`
}

type ProcConfig struct {
	Namespace string `toml:"namespace" json:"namespace"`
}

type WorkerConfig struct {
	Num int `toml:"num" json:"num"`
}

type RuntimeConfig struct {
	MaxProcs int `toml:"maxProcs" json:"maxProcs"`
}

type DataBaseConfig struct {
	Db DbConfig `toml:"db" json:"db"`
}

type DbConfig struct {
	MaxIdle int      `toml:"maxIdle" json:"maxIdle"`
	MaxOpen int      `toml:"maxOpen" json:"maxOpen"`
	Debug   bool     `toml:"debug" json:"debug"`
	Addr    []string `toml:"addr" json:"addr"`
}

type RegionConfig struct {
	Baidu  []string `toml:"baidu" json:"baidu"`
	Qcloud []string `toml:"qcloud" json:"qcloud"`
}

type CronConfig struct {
	Server      string `toml:"server" json:"server"`
	Cache       bool   `toml:"cache" json:"cache"`
	CacheSecond string `toml:"cacheSecond" json:"cacheSecond"`
}

var Config Conf

func MustInit(path string) {
	_, err := toml.DecodeFile(path, &Config)
	if err != nil {
		panic(fmt.Sprintf("load config file error, [file: %s][error: %s]", path, err.Error()))
	}
}

func Configs() map[string]interface{} {
	return map[string]interface{}{
		"config": Config,
	}
}

func (this Conf) GetOauthUrl(name string) (res *ApiConfig, err error) {
	api, ok := Config.OauthWeb.Apis[name]
	if !ok {
		ret := &ApiConfig{
			Uri:        "",
			TimeoutSec: 0,
		}
		return ret, fmt.Errorf("config %s url not exisit", name)
	}
	ret := &api
	ret.Uri = Config.OauthWeb.Apis[name].Domain + ret.Uri
	return ret, nil
}
