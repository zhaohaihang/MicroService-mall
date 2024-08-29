package config

type ApiConfig struct {
	Name            string            `json:"name"`
	Host            string            `json:"host"`
	ServiceInfo     Register          `json:"register"`
	JWTInfo         JwtConfig         `json:"jwt"`
	AliSmsInfo      AliSmsConfig      `json:"aliyun_message"`
	RedisInfo       RedisConfig       `json:"redis"`
	ConsulInfo      ConsulConfig      `json:"consul"`
	JaegerInfo      JaegerConfig      `json:"jaeger_info"`
	UserServiceInfo UserServiceConfig `json:"user_service_info"`
}

type JwtConfig struct {
	SigningKey string `json:"key"`
}

type AliSmsConfig struct {
	ApiKey       string `json:"key"`
	ApiSecret    string `json:"secret"`
	SignName     string `json:"signName"`
	TemplateCode string `json:"template_code"`
	RegionId     string `json:"region_id"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
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

type UserServiceConfig struct {
	Name string `json:"name"`
}

type FileConfig struct {
	ConfigFile string
	LogFile    string
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

type Register struct {
	Tags           []string `json:"tags"`
	CheckTimeOut   string   `json:"check_time_out"`
	CheckInterval  string   `json:"check_interval"`
	DeregisterTime string   `json:"deregister_time"`
}
