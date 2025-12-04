package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go" // 旧版本 JWT 库
	"github.com/gin-gonic/gin"
)

var JwtSecret = []byte("wunder_minibar_secret_key_2025")

// MyClaims
type MyClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Role     int    `json:"role"`
	jwt.StandardClaims
}

func GenToken(userID int, username string, nickname string, role int) (string, error) {
	claims := MyClaims{
		UserID:   userID,
		Username: username,
		Nickname: nickname,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Issuer:    "wunder_minibar_backend",
			NotBefore: time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if err != nil {
		fmt.Println("Token解析失败原因:", err)
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// ---------------- 业务结构体 ----------------

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserInfoResponse struct {
	ID         int    `json:"id"`
	CreatedAt  string `json:"created_at"`
	NickName   string `json:"nick_name"`
	UserName   string `json:"username"` // 这里对应前端 user_name, 但你之前改成 username 了，保持一致
	Avatar     string `json:"avatar"`
	Email      string `json:"email"`
	Tel        string `json:"tel"`
	Addr       string `json:"addr"`
	Token      string `json:"token"`
	IP         string `json:"ip"`
	Role       int    `json:"role"`
	SignStatus string `json:"sign_status"`
	Integral   int    `json:"integral"`
	Sign       string `json:"sign"`
	Link       string `json:"link"`
}

// PageResult 分页返回结构 (新增)
type PageResult struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func main() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, token")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		// --- 1. 登录接口 ---
		api.POST("/login", func(c *gin.Context) {
			var req LoginReq
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusOK, Response{Code: 400, Msg: "参数错误", Data: nil})
				return
			}

			if req.Username == "admin" && req.Password == "123456" {
				token, err := GenToken(1, "admin", "万妙小吧", 1)
				if err != nil {
					c.JSON(http.StatusOK, Response{Code: 500, Msg: "Token生成失败", Data: nil})
					return
				}

				c.JSON(http.StatusOK, Response{
					Code: 0,
					Msg:  "登录成功",
					Data: token,
				})
			} else {
				c.JSON(http.StatusOK, Response{Code: 1001, Msg: "用户名或密码错误", Data: nil})
			}
		})

		// --- 2. 注销接口 ---
		api.POST("/logout", func(c *gin.Context) {
			c.JSON(http.StatusOK, Response{
				Code: 0,
				Msg:  "注销成功",
				Data: "logout success",
			})
		})

		// --- 3. 获取用户信息接口 ---
		api.GET("/user_info", func(c *gin.Context) {
			tokenStr := c.GetHeader("token")
			if tokenStr == "" {
				c.JSON(http.StatusOK, Response{Code: 401, Msg: "未携带Token", Data: nil})
				return
			}

			claims, err := ParseToken(tokenStr)
			if err != nil {
				c.JSON(http.StatusOK, Response{Code: 401, Msg: "Token无效或已过期", Data: nil})
				return
			}

			userInfo := UserInfoResponse{
				ID:         claims.UserID,
				UserName:   claims.Username,
				CreatedAt:  time.Now().Format("2006-01-02"),
				NickName:   "Luk4sy",
				Avatar:     "https://p3-pc-sign.douyinpic.com/tos-cn-i-0813/oQHf9mAbvFDEALmpAAQAVlxf1xUgyFIDANCCay~tplv-dy-aweme-images:q75.webp?biz_tag=aweme_images&from=327834062&lk3s=138a59ce&s=PackSourceEnum_SEARCH&sc=image&se=false&x-expires=1766448000&x-signature=6nKHU2vkRJgt84Qv%2FVNls%2FhdiOU%3D",
				Email:      "admin@wunderminibar.com",
				Role:       claims.Role,
				Token:      tokenStr,
				SignStatus: "邮箱",
				Integral:   100,
			}

			c.JSON(http.StatusOK, Response{Code: 0, Msg: "获取成功", Data: userInfo})
		})

		// --- 4. 获取用户列表接口 (新增) ---
		api.GET("/users", func(c *gin.Context) {
			// 鉴权 (实际开发建议封装成中间件)
			tokenStr := c.GetHeader("token")
			if tokenStr == "" {
				c.JSON(http.StatusOK, Response{Code: 401, Msg: "未携带Token", Data: nil})
				return
			}
			_, err := ParseToken(tokenStr)
			if err != nil {
				c.JSON(http.StatusOK, Response{Code: 401, Msg: "Token无效或已过期", Data: nil})
				return
			}

			// 模拟 5 个用户数据
			mockList := []UserInfoResponse{
				{
					ID: 1, CreatedAt: "2023-02-19", NickName: "Luk4sy", UserName: "admin",
					Avatar: "https://p3-pc-sign.douyinpic.com/tos-cn-i-0813/oQHf9mAbvFDEALmpAAQAVlxf1xUgyFIDANCCay~tplv-dy-aweme-images:q75.webp",
					Email:  "admin@wunder.com", Role: 1, SignStatus: "邮箱", Integral: 9999, Addr: "天津",
				},
				{
					ID: 2, CreatedAt: "2023-03-01", NickName: "张三", UserName: "zhangsan",
					Avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Felix",
					Email:  "zs@test.com", Role: 2, SignStatus: "手机", Integral: 100, Addr: "北京",
				},
				{
					ID: 3, CreatedAt: "2023-04-15", NickName: "李四", UserName: "lisi_007",
					Avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Aneka",
					Email:  "lisi@test.com", Role: 2, SignStatus: "微信", Integral: 50, Addr: "上海",
				},
				{
					ID: 4, CreatedAt: "2023-05-20", NickName: "王五", UserName: "wangwu",
					Avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Bob",
					Email:  "ww@test.com", Role: 2, SignStatus: "QQ", Integral: 200, Addr: "广州",
				},
				{
					ID: 5, CreatedAt: "2023-06-01", NickName: "赵六", UserName: "zhaoliu",
					Avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Jack",
					Email:  "zl@test.com", Role: 3, SignStatus: "邮箱", Integral: 0, Addr: "深圳",
				},
			}

			c.JSON(http.StatusOK, Response{
				Code: 0,
				Msg:  "获取用户列表成功",
				Data: PageResult{
					List:  mockList,
					Total: int64(len(mockList)),
				},
			})
		})
	}

	r.Run(":8080")
}
