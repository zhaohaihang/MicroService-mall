package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"github.com/zhaohaihang/goods_api/config"
	"github.com/zhaohaihang/goods_api/global"
	"go.uber.org/zap"
)

func InitConfig() {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(global.FilePath.ConfigFile)
	err := v.ReadInConfig()
	if err != nil {
		zap.S().Errorw("load local nacos config failed", "err", err.Error())
		return
	}
	global.NacosConfig = &config.NacosConfig{}
	err = v.Unmarshal(global.NacosConfig)
	if err != nil {
		zap.S().Errorw("unmarshal local nacos config failed", "err", err.Error())
		return
	}
	zap.S().Infof("nacos config context is : %#v", global.NacosConfig)

	sConfig := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   uint64(global.NacosConfig.Port),
		},
	}
	nacosLogDir := fmt.Sprintf("%s/%s/%s", global.FilePath.LogFile, "nacos", "log")
	nacosCacheDir := fmt.Sprintf("%s/%s/%s", global.FilePath.LogFile, "nacos", "cache")
	cConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              nacosLogDir,
		CacheDir:            nacosCacheDir,
		LogLevel:            "debug",
	}
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sConfig,
		"clientConfig":  cConfig,
	})
	if err != nil {
		zap.S().Errorw("nacos client create failed", "err", err.Error())
		return
	}

	content, err := client.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.Dataid,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		zap.S().Errorw("load config from nacos server failed", "err", err.Error())
		return
	}

	global.ApiConfig = &config.ApiConfig{}
	err = json.Unmarshal([]byte(content), global.ApiConfig)
	if err != nil {
		zap.S().Errorw("load nacos content to global.serviceConfig failed", "err", err.Error())
		return
	}
	zap.S().Infof("pull nacos config success %#v", global.ApiConfig)

	//监听配置修改
	err = client.ListenConfig(vo.ConfigParam{
		DataId: "goods_api.json",
		Group:  "dev",
		OnChange: func(namespace, group, dataId, data string) {
			// TODO 配置变化时，应该重新反序列化，并且重新初始化一些公共资源
		},
	})
	if err != nil {
		zap.S().Fatalw("listen goods_api config from nacos failed: %s", "err", err.Error())
	}
	zap.S().Info("listening nacos config change")
}
