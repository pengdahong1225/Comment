package settings

// 从viper反射到数据模型，需要设置`mapstructure`反射字段
type AppConfig struct {
	*SandBox        `mapstructure:"sandbox"`
	*MysqlConfig    `mapstructure:"mysql"`
	*RedisConfig    `mapstructure:"redis"`
	*RegistryConfig `mapstructure:"registry"`
	*MqConfig       `mapstructure:"rabbitmq"`
	*JwtConfig      `mapstructure:"jwt"`
	*SmsConfig      `mapstructure:"sms"`
}

type SandBox struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type MysqlConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Db   string `mapstructure:"db"`
	User string `mapstructure:"user"`
	Pwd  string `mapstructure:"password"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"key"`
}

type SmsConfig struct {
	AccessKeyId     string `mapstructure:"accessKeyId"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
	Endpoint        string `mapstructure:"endpoint"`
	SignName        string `mapstructure:"signName"`
	TemplateCode    string `mapstructure:"templateCode"`
}

type RegistryConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type MqConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	PassWord string `mapstructure:"password"`
	VHost    string `mapstructure:"vhost"`
}
