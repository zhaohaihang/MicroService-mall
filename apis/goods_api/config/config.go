package config


type ApiConfig struct {
	Name       string       `json:"name"`
	Host       string       `json:"host"`
	ConsulInfo ConsulConfig `json:"consul"`
	JWTInfo    JWTConfig    `json:"jwt"`
	JaegerInfo JaegerConfig `json:"jaeger_info"`
	ServiceInfo 	Register     `json:"register"`
	GoodsServiceInfo GoodsServiceConfig `json:"goods_service_info"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Dataid    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}

type ConsulConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type JaegerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type FilePathConfig struct {
	ConfigFile string
	LogFile    string
}

type Register struct {
	Tags           []string `json:"tags"`
	CheckTimeOut   string   `json:"check_time_out"`
	CheckInterval  string   `json:"check_interval"`
	DeregisterTime string   `json:"deregister_time"`
}

type GoodsServiceConfig struct {
	Name string `json:"name"`
}