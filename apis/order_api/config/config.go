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

type ApiConfig struct {
	Name             string `json:"name"`
	Host             string  `json:"host"`
	ServiceInfo 	Register     `json:"register"`
	ConsulInfo       ConsulConfig           `json:"consul"`
	JWTInfo          JWTConfig              `json:"jwtConfig"`
	OrderService 	 OrderServiceConfig	 	`json:"order_service"`
	GoodsService     GoodsServiceConfig     `json:"goods_service"`
	InventoryService InventoryServiceConfig `json:"inventory_service"`
	AlipayInfo       AlipayInfoConfig       `json:"alipay_info"`
	JaegerInfo       JaegerConfig           `json:"jaeger"`
}

type ConsulConfig struct {
	Host string   `json:"host"`
	Port int      `json:"port"`
	Tags []string `json:"tags"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type JaegerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type OrderServiceConfig struct {
	Name string `json:"name"`
}

type GoodsServiceConfig struct {
	Name string `json:"name"`
}

type InventoryServiceConfig struct {
	Name string `json:"name"`
}

type AlipayInfoConfig struct {
	AppID        string `mapstructure:"app_id" json:"app_id"`
	PrivateKey   string `mapstructure:"private_key" json:"private_key"`
	AliPublicKey string `mapstructure:"ali_public_key" json:"ali_public_key"`
	NotifyURL    string `mapstructure:"notify_url" json:"notify_url"`
	ReturnURL    string `mapstructure:"return_url" json:"return_url"`
}

type Register struct {
	Tags           []string `json:"tags"`
	CheckTimeOut   string   `json:"check_time_out"`
	CheckInterval  string   `json:"check_interval"`
	DeregisterTime string   `json:"deregister_time"`
}
