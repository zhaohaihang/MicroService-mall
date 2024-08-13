package config

type NacosConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Namespace string `yaml:"namespace"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Dataid    string `yaml:"dataid"`
	Group     string `yaml:"group"`
}

type ApiConfig struct {
	Name string `json:"name"`
	Host string `json:"host"`

	ServiceInfo   Register            `json:"register"`
	ConsulInfo    ConsulConfig        `json:"consul"`
	JWTInfo       JWTConfig           `json:"jwtConfig"`
	UseropService UserOPServiceConfig `json:"userop_service"`
	GoodsService  GoodsServiceConfig  `json:"goods_service"`
	JaegerInfo    JaegerConfig        `json:"jaeger"`
}

type ConsulConfig struct {
	Host string   `json:"host"`
	Port int      `json:"port"`
	Tags []string `json:"tags"`
}

type JWTConfig struct {
	SigningKey string `json:"key"`
}

type JaegerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type UserOPServiceConfig struct {
	Name string `json:"name"`
}

type GoodsServiceConfig struct {
	Name string `json:"name"`
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
