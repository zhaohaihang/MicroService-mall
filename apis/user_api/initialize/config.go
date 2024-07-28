package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"github.com/zhaohaihang/user_api/config"
	"github.com/zhaohaihang/user_api/global"
	"go.uber.org/zap"
)

// InitConfig 初始化配置
func InitConfig() {

	// 加载本地文件中nacos的配置
	v := viper.New()
	v.SetConfigFile(global.FileConfig.ConfigFile)
	err := v.ReadInConfig()
	if err != nil {
		zap.S().Errorw("viper.ReadInConfig failed", "err", err.Error())
		return
	}
	global.NacosConfig = &config.NacosConfig{}
	err = v.Unmarshal(global.NacosConfig)
	if err != nil {
		zap.S().Errorw("viper unmarshal failed", "err", err.Error())
		return
	}
	zap.S().Infof("global.NacosConfig : %#v", global.NacosConfig)

	// 创建nacos客户端
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
		LogDir:              fmt.Sprintf("%s/%s/%s", global.FileConfig.LogFile, "nacos", "log"),
		CacheDir:            fmt.Sprintf("%s/%s/%s", global.FileConfig.LogFile, "nacos", "cache"),
		LogLevel:            "debug",
	}
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sConfig,
		"clientConfig":  cConfig,
	})
	if err != nil {
		zap.S().Errorw("nacos client conn failed", "err", err.Error())
		return
	}

	// 从nacos拉取配置
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.Dataid,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		zap.S().Errorw("client.GetConfig load failed", "err", err.Error())
		return
	}

	// 加载配置
	global.ApiConfig = &config.ApiConfig{}
	err = json.Unmarshal([]byte(content), global.ApiConfig)
	if err != nil {
		zap.S().Errorw("load nacos content to global.serviceConfig failed", "err", err.Error())
		return
	}
	zap.S().Infof("pull nacos config success %#v", global.ApiConfig)


		//监听配置修改
		err = client.ListenConfig(vo.ConfigParam{
			DataId: "user_api.json",
			Group:  "dev",
			OnChange: func(namespace, group, dataId, data string) {
				// TODO 配置变化时，应该重新反序列化，并且重新初始化一些公共资源
			},
		})
		if err != nil {
			zap.S().Fatalw("listen user_api config from nacos failed: %s", "err", err.Error())
		}
		zap.S().Info("listening nacos config change")
}
