package v1

var (
	// common errors
	ErrSuccess             = newError(0, "ok")
	ErrBadRequest          = newError(400, "Bad Request")
	ErrUnauthorized        = newError(401, "Unauthorized")
	ErrUnauthorizedAP      = newError(401, "账号或者密码不正确")
	ErrNotFound            = newError(404, "Not Found")
	ErrInternalServerError = newError(500, "Internal Server Error")

	// more biz errors
	ErrEmailAlreadyUse = newError(1001, "邮箱地址已被使用")
	ErrRePasswd        = newError(1001, "两次密码输入不正确")
	ErrPhone           = newError(1001, "用户手机号为空或者不正确,请更新手机号")

	ErrCashRequest        = newError(2002, "账户余额不足,提现发起失败")
	ErrCashAccountRequest = newError(2003, "账户状态异常,提现发起失败,请联系客服")
	ErrCashPicRequest     = newError(2004, "账户提现方式为空,提现发起失败,请填写提现方式")
	ErrCashNow            = newError(2005, "有未结束的提现申请,请等待之前申请结束后,重新发起")
)
