package models

type UserLoginModel struct {
	Model
	UserID    uint      `json:"userID"`
	UserModel UserModel `gorm:"foreignKey:UserID" json:"-"` // 关联用户
	IP        string    `gorm:"size:32" json:"ip"`          // 登录IP
	Addr      string    `gorm:"size:64" json:"addr"`        // IP归属地
	UA        string    `gorm:"size:128" json:"ua"`         // 用户代理（浏览器、设备信息）
}
