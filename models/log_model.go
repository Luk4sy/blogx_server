package models

type LogModel struct {
	Model
	LogType   int8      `json:"logType"`                    // 日志类型
	Title     string    `json:"title"`                      // 日志标题
	Content   string    `json:"content"`                    // 日志内容
	Level     int8      `json:"level"`                      // 日志级别
	UserID    uint      `json:"userID"`                     // 用户ID
	UserModel UserModel `gorm:"foreignKey:UserID" json:"-"` // 关联用户信息
	IP        string    `json:"ip"`                         // 操作IP
	Addr      string    `json:"addr"`                       // IP归属地
	IsRead    bool      `json:"isRead"`                     // 是否读取
}
