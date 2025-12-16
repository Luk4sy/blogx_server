package models

import "minibar_server/models/enum"

type LogModel struct {
	Model
	LogType     enum.LogType      `json:"logType"`                    // 日志类型
	Title       string            `gorm:"size:64" json:"title"`       // 日志标题
	Content     string            `json:"content"`                    // 日志内容
	Level       enum.LogLevelType `json:"level"`                      // 日志级别
	UserID      uint              `json:"userID"`                     // 用户ID
	UserModel   UserModel         `gorm:"foreignKey:UserID" json:"-"` // 关联用户信息
	IP          string            `gorm:"size:32" json:"ip"`          // 操作IP
	Addr        string            `gorm:"size:64" json:"addr"`        // IP归属地
	IsRead      bool              `json:"isRead"`                     // 是否读取
	LoginStatus bool              `json:"loginStatus"`                // 登陆状态
	Username    string            `json:"username"`                   // 登录日志的用户名
	Pwd         string            `json:"pwd"`                        // 登录日志的密码
	LoginType   enum.LoginType    `json:"loginType"`                  // 登陆的类型
	ServiceName string            `gorm:"size:32" json:"serviceName"`
}
