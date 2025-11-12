package models

// GlobalNotificationModel 全局通知表
type GlobalNotificationModel struct {
	Model
	Title   string `gorm:"size:32" json:"title"`   // 通知标题
	Icon    string `gorm:"size:256" json:"icon"`   // 图标链接
	Content string `gorm:"size:64" json:"content"` // 通知内容
	Href    string `gorm:"size:256" json:"href"`   // 用户点击消息跳转链接
}
