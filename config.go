package main

import (
	"fmt"
	"goshop/service-member/pkg/utils"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//初始化配置文件
func InitConfig() {
	buf := &utils.Config{}
	UnmarshalYaml(getConfigFile(), buf)
	utils.C = buf

	//app解析
	baseInfo := &utils.Base{}
	UnmarshalYaml("./conf/app.yaml", baseInfo)
	utils.C.Base = baseInfo
}

/**
获取到配置文件路径
*/
func getConfigFile() string {
	if utils.FileExists("./conf/local.app.yaml") {
		return "./conf/local.app.yaml"
	}

	return fmt.Sprintf("./conf/%s.app.yaml", gin.Mode())
}

func UnmarshalYaml(fileName string, data interface{}) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(fileName)
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败, err: %v", err))
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("%s文件发生改动哦 \n", e.Name)
		if err := v.Unmarshal(data); err != nil {
			panic(fmt.Sprintf("解析配置文件出错, err: %v", err))
		}
	})

	if err := v.Unmarshal(data); err != nil {
		panic(fmt.Sprintf("解析配置文件出错, err: %v", err))
	}

}
