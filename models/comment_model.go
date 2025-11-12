package models

type CommentModel struct {
	Model
	Content        string          `gorm:"size:256" json:"content"`                   // 评论内容
	UserID         uint            `json:"userID"`                                    // 评论者 ID
	UserModel      UserModel       `gorm:"foreignKey:UserID" json:"-"`                // 评论者信息
	ArticleID      uint            `json:"articleID"`                                 // 所属文章 ID
	ArticleModel   ArticleModel    `gorm:"foreignKey:ArticleID" json:"-"`             // 所属文章信息
	ParentID       *uint           `json:"parentID"`                                  // 父评论 ID
	ParentModel    *CommentModel   `gorm:"foreignKey:ParentID" json:"-"`              // 父评论对象
	SubCommentList []*CommentModel `gorm:"foreignKey:ParentID" json:"subCommentList"` // 子评论列表
	RootParentID   *uint           `json:"rootParentID"`                              // 根评论 ID（顶层评论）
	DiggCount      int             `json:"diggCount"`                                 // 点赞数
}
