package conf

type Email struct {
	Domain       string `yaml:"domain" json:"domain"`
	Port         int    `yaml:"port" json:"port"`
	SendEmail    string `yaml:"sendEmail" json:"sendEmail"`
	AuthCode     string `yaml:"authCode" json:"authCode"` // 授权码
	SendNickname string `yaml:"sendNickname" json:"sendNickname"`
	IsSSL        bool   `yaml:"isSSL" json:"isSSL"` // 是否开启SSL
	IsTLS        bool   `yaml:"isTLS" json:"isTLS"` // 是否开启TLS
}
