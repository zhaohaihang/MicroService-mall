package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"github.com/zhaohaihang/user_service/config"
	"github.com/zhaohaihang/user_service/global"
	"go.uber.org/zap"
)

// InitConfig 初始化配置
func InitConfig() {
	// 获得配置文件路径
	configFileName := fmt.Sprintf(global.FilePath.ConfigFile)
	v := viper.New()
	v.SetConfigFile(configFileName)
	err := v.ReadInConfig()
	if err != nil {
		zap.S().Fatalw("read local nacos config failed: %s", "err", err.Error())
	}
	global.NacosConfig = &config.NacosConfig{}
	err = v.Unmarshal(global.NacosConfig)
	if err != nil {
		zap.S().Fatalw("unmarshal local nacos config failed: %s", "err", err.Error())
	}
	zap.S().Infof("%v",global.NacosConfig)

	// 连接nacos
	sConfig := []constant.ServerConfig{{
		IpAddr: global.NacosConfig.Host,
		Port:   uint64(global.NacosConfig.Port),
	}}
	cConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              fmt.Sprintf("%s/%s/%s", global.FilePath.LogFile, "nacos", "log"),
		CacheDir:            fmt.Sprintf("%s/%s/%s", global.FilePath.LogFile, "nacos", "cache"),
		LogLevel:            "debug",
	}
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sConfig,
		"clientConfig":  cConfig,
	})
	if err != nil {
		zap.S().Fatalw("create nacos client failed: %s", "err", err.Error())
	}

	// 从nacos拉取配置
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.Dataid,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		zap.S().Fatalw("pull user sercice config from nacos failed", "err", err.Error())
	}

	// 反序列化
	global.ServiceConfig = &config.ServiceConfig{}
	err = json.Unmarshal([]byte(content), global.ServiceConfig)
	if err != nil {
		zap.S().Fatalw("Unmarshal user sercice config failed: %s", "err", err.Error())
	}
	zap.S().Info("load user service  config from nacos success ")

}
