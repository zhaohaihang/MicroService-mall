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

type FilePathConfig struct {
	ConfigFile string
	LogFile    string
}

type MySqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	Name     string `mapstructure:"name" json:"name"`
}

type ConsulConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
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

type EsConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type ServiceConfig struct {
	Name         string `json:"name"`
	Host         string `json:"host"`
	MySqlInfo    MySqlConfig    `json:"mysql"`
	ConsulInfo   ConsulConfig   `json:"consul"`
	EsInfo       EsConfig       `json:"es"`
	JaegerInfo   JaegerConfig   `json:"jaeger"`
	RegisterInfo RegisterConfig `json:"register"`
}

