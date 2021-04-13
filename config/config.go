package config

import (
	"encoding/json"
	"github.com/Zhenghao-Liu/OAuth_demo/common"
	"io/ioutil"
)

type Config struct {
	OAuthDemoDB    MysqlConfig `json:"oauth_demo_db"`
	OAuthDemoCache RedisConfig `json:"oauth_demo_cache"`
}

type MysqlConfig struct {
	Database    string `json:"database"`
	Settings    string `json:"settings"`
	UserName    string `json:"username"`
	Password    string `json:"password"`
	DefaultHost string `json:"default_host"`
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DBIdx    int    `json:"db_idx"`
}

var ConfigInstance *Config

func NewConfigInstance() error {
	content, err := ioutil.ReadFile(common.ConfFile)
	if err != nil {
		return err
	}
	var conf Config
	err = json.Unmarshal(content, &conf)
	if err != nil {
		return err
	}
	ConfigInstance = &conf
	return nil
}
