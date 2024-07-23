package config

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Dataid    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}

type Register struct {
	Tags           []string `json:"tags"`
	CheckTimeOut   string   `json:"check_time_out"`
	CheckInterval  string   `json:"check_interval"`
	DeregisterTime string   `json:"deregister_time"`
}

type MysqlConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type ConsulConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

type JaegerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type FilePathConfig struct {
	ConfigFile string
	LogFile    string
}

type ServiceConfig struct {
	Name        string `json:"name"`
	Host        string
	ServiceInfo Register     `json:"register"`
	MysqlInfo   MysqlConfig  `json:"mysql"`
	ConsulInfo  ConsulConfig `json:"consul"`
	RedisInfo   RedisConfig  `json:"redis"`
	JaegerInfo  JaegerConfig `json:"jaeger"`
}