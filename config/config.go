package config

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/abdullahPrasetio/prasegateway/entity"
)

var configOnce sync.Once
var DataConfig entity.MyConfig

func InitializeConfig() {
	file, err := os.ReadFile("prase.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(file, &DataConfig); err != nil {
		panic(err)
	}
}

func GetMyConfig() entity.MyConfig {
	configOnce.Do(InitializeConfig)
	return DataConfig
}
