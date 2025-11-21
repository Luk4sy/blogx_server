package api

import (
	"blogx_server/api/banner_api"
	"blogx_server/api/captcha_api"
	"blogx_server/api/image_api"
	"blogx_server/api/log_api"
	"blogx_server/api/site_api"
	"blogx_server/api/user_api"
)

type Api struct {
	SiteApi    site_api.SiteApi
	LogApi     log_api.LogApi
	ImageApi   image_api.ImageApi
	BannerApi  banner_api.BannerApi
	CaptchaApi captcha_api.CaptchaApi
	UserApi    user_api.UserApi
}

var App = Api{}
