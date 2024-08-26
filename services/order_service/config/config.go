package config

type FilePathConfig struct {
	ConfigFile string
	LogFile    string
}

type ServiceConfig struct {
	Name                 string                 `json:"name"`
	Host                 string                 `json:"host"`
	MysqlInfo            MysqlConfig            `json:"mysql"`
	ConsulInfo           ConsulConfig           `json:"consul"`
	RedisInfo            RedisConfig            `json:"redis"`
	JaegerInfo           JaegerConfig           `json:"jaeger"`
	RegisterInfo         RegisterConfig         `json:"register"`
	GoodsServiceInfo     GoodsServiceConfig     `json:"goods_service"`
	InventoryServiceInfo InventoryServiceConfig `json:"inventory_service"`
	RocketMQInfo         RocketMQConfig         `json:"rocketmq"`
}

type RocketMQConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
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

type MysqlConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type ConsulConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	PoolSize int    `json:"pool_size"`
}

type JaegerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type RegisterConfig struct {
	Tags           []string `json:"tags"`
	CheckTimeOut   string   `json:"check_time_out"`
	CheckInterval  string   `json:"check_interval"`
	DeregisterTime string   `json:"deregister_time"`
}

type GoodsServiceConfig struct {
	Name string `json:"name"`
}

type InventoryServiceConfig struct {
	Name string `json:"name"`
}
