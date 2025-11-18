package conf

import "fmt"

type QQ struct {
	AppId    string `yaml:"appId" json:"appId"`
	AppKey   string `yaml:"appKey" json:"appKey"`
	Redirect string `yaml:"redirect" json:"redirect"`
}

func (q QQ) Url() string {
	return fmt.Sprintf("https://graph.qq.com/oauth2.0/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=get_user_info",
		q.AppId, q.Redirect)
}
