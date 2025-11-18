package conf

type Ai struct {
	Enable    bool   `yaml:"enable" json:"enable"`
	SecretKey string `yaml:"secretKey" json:"secretKey"`
	NickName  string `yaml:"nickName" json:"nickName"`
	Avatar    string `yaml:"avatar" json:"avatar"`
}
