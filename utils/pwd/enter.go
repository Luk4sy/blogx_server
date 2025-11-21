package pwd

import "golang.org/x/crypto/bcrypt"

// / 加密密码
func GenerateFromPassword(password string) (string, error) {
	// 使用 bcrypt 生成哈希密码，第二个参数是成本（越大越安全但越耗时）
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// / 校验密码
func CompareHashAndPassword(hashedPassword string, password string) bool {
	// 对比明文密码和哈希密码，如果一致 err 为 nil
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
