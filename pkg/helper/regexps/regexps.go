package regexps

import "regexp"

func ValidatePhoneNumber(phoneNumber string) bool {
	// 定义手机号验证的正则表达式
	re := regexp.MustCompile(`^(?:\+?86)?1[3-9]\d{9}$`)

	// 使用正则表达式匹配手机号
	isValid := re.MatchString(phoneNumber)

	return isValid
}

func ValidURL(input string) bool {
	re := regexp.MustCompile(`^http://download\.pduola\.com/v1/p/u/d/.*`)
	match := re.MatchString(input)
	return match

}
