package wallet

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"unicode"
)

// generateSecurePassphrase 生成安全passphrase
func GenerateSecurePassphrase() (string, error) {
	// 生成16字节随机数（128位熵）
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("随机数生成失败: %v", err)
	}

	// 使用URL安全的Base64编码
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

// 带版本控制的生成方式
func GeneratePassphrase() (string, error) {
	const (
		version    = 1 // 密码格式版本
		saltLength = 4 // 校验盐长度
	)

	// 生成随机熵
	entropy := make([]byte, 16)
	if _, err := rand.Read(entropy); err != nil {
		return "", err
	}

	// 生成校验盐
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 组合数据：版本 + 盐 + 熵
	data := append([]byte{byte(version)}, salt...)
	data = append(data, entropy...)

	// 使用Base62编码（更友好的显示格式）
	encoded := base62Encode(data)
	return fmt.Sprintf("v%d_%s", version, encoded), nil
}

// base62Encode 自定义Base62编码
func base62Encode(data []byte) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var result []byte

	// 将字节转换为大整数
	n := new(big.Int).SetBytes(data)
	base := big.NewInt(62)

	for n.Cmp(big.NewInt(0)) > 0 {
		mod := new(big.Int)
		n.DivMod(n, base, mod)
		result = append(result, charset[mod.Int64()])
	}

	// 反转结果
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

// VerifyPassphraseComplexity 密码复杂度验证
func VerifyPassphraseComplexity(passphrase string) error {
	if len(passphrase) < 16 {
		return fmt.Errorf("密码长度需至少16字符")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, c := range passphrase {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	if !(hasUpper && hasLower && hasDigit && hasSpecial) {
		return fmt.Errorf("密码需包含大小写字母、数字和特殊字符")
	}

	return nil
}
