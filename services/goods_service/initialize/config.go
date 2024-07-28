package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"github.com/zhaohaihang/goods_service/config"
	"github.com/zhaohaihang/goods_service/global"
	"go.uber.org/zap"
)

// InitConfig 初始化配置
func InitConfig() {
	// 获得配置文件路径
	configFileName := fmt.Sprintf(global.FilePath.ConfigFile)
	// 生成viper, 读取nacos的配置
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
	zap.S().Infof("local nacos config is :%v",global.NacosConfig)
	// 连接nacos
	sConfig := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   uint64(global.NacosConfig.Port),
		},
	}

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
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.Dataid,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		zap.S().Fatalw("pull goods sercice config from nacos failed", "err", err.Error())
		return
	}
	global.ServiceConfig = &config.ServiceConfig{}
	err = json.Unmarshal([]byte(content), global.ServiceConfig)
	if err != nil {
		zap.S().Fatalw("Unmarshal goods sercice config failed: %s", "err", err.Error())
	}
	zap.S().Info("load goods service  config from nacos success ")

	//监听配置修改
	err = client.ListenConfig(vo.ConfigParam{
		DataId: "goods_service.json",
		Group:  "dev",
		OnChange: func(namespace, group, dataId, data string) {
			// TODO 配置变化时，应该重新反序列化，并且重新初始化一些公共资源
		},
	})
	if err != nil {
		zap.S().Fatalw("listen goods_sercice config from nacos failed: %s", "err", err.Error())
	}
	zap.S().Info("listening nacos config change")
}
