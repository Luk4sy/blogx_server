package site

type SiteInfo struct {
	Title string `yaml:"title" json:"title"`
	Logo  string `yaml:"logo" json:"logo"`
	IcpNo string `yaml:"icpNo" json:"icpNo"`
	Mode  int8   `yaml:"mode" json:"mode" binding:"oneof=1 2"` // 1 社区模式 2 博客模式
}

type Project struct {
	Title   string `yaml:"title" json:"title"`
	Icon    string `yaml:"icon" json:"icon"`
	WebPath string `yaml:"webPath" json:"webPath"`
}

type Seo struct {
	Keywords    string `yaml:"keywords" json:"keywords"`
	Description string `yaml:"description" json:"description"`
}

type About struct {
	SiteDate       string `yaml:"siteDate" json:"siteDate"` // 年月日
	QQQrCode       string `yaml:"qqQrCode" json:"QQQrCode"`
	Version        string `yaml:"-" json:"version"`
	WechatQrCode   string `yaml:"wechatQrCode" json:"wechatQrCode"`
	GiteeQrCode    string `yaml:"giteeQrCode" json:"giteeQrCode"`
	BilibiliQrCode string `yaml:"bilibiliQrCode" json:"bilibiliQrCode"`
	GithubQr       string `yaml:"githubQr" json:"githubQr"`
}

type Login struct {
	QQLogin          bool `yaml:"qqLogin" json:"qqLogin"`
	UsernamePwdLogin bool `yaml:"usernamePwdLogin" json:"usernamePwdLogin"`
	EmailLogin       bool `yaml:"emailLogin" json:"emailLogin"`
	Captcha          bool `yaml:"captcha" json:"captcha"`
}

type ComponentInfo struct {
	Title  string `yaml:"title" json:"title"`
	Enable bool   `yaml:"enable" json:"enable"`
}

type IndexRight struct {
	List []ComponentInfo `json:"list" yaml:"list"`
}

type Article struct {
	IsFreeReview bool `json:"isFreeReview" yaml:"isFreeReview"` // 是否免审核
}
