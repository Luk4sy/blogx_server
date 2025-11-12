package models

import "fmt"

type ImageModel struct {
	Model
	Filename string `gorm:"size:64"  json:"filename"`
	Path     string `gorm:"size:256" json:"path"`
	Size     int64  `json:"size"`
	Hash     string `gorm:"size:32"  json:"hash"`
}

// WebPath 返回图片的 Web 访问路径
func (i ImageModel) WebPath() string {
	return fmt.Sprintf("/")
}
