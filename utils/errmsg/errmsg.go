package errmsg

const (
	SUCCESS = 200
	ERROR   = 500
	//1000到1999用户模块
	ERROR_USERNAME_USED    = 1001
	ERROR_EMAIL_USED       = 1010
	ERROR_PASSWORD_WEONG   = 1002
	ERROR_USER_NOT_EXIST   = 1003
	ERROR_TOKEN_EXIST      = 1004
	ERROR_TOKEN_RUNTIME    = 1005
	ERROR_TOKEN_WRONG      = 1006
	ERROR_TOKEN_TYPE_WRONG = 1007
	ERROR_USER_EDIT        = 1008
	ERROR_USER_NO_RIGHT    = 1009
	ERROR_EMAIL_NIL        = 1011
	ERROR_EMAIL_NO         = 1012

	//2000-2999文章模块错误

	ERROR_ART_NOT_EXIST = 2001
	//3000-3999分类模块出错
	ERROR_CATENAME_USED  = 3001
	ERROR_CATE_NOT_EXIST = 3002
)

var Codemsg = map[int]string{
	SUCCESS: "ok",
	ERROR:   "fail",
	//用户
	ERROR_USERNAME_USED:    "用户已存在",
	ERROR_PASSWORD_WEONG:   "密码错误",
	ERROR_USER_NOT_EXIST:   "用户不存在",
	ERROR_TOKEN_EXIST:      "TOKEN不存在",
	ERROR_TOKEN_RUNTIME:    "TOKEN已过期",
	ERROR_TOKEN_WRONG:      "TOKEN不正确",
	ERROR_TOKEN_TYPE_WRONG: "TOKEN格式错误",
	ERROR_USER_EDIT:        "修改成功",
	ERROR_USER_NO_RIGHT:    "权限不足",
	ERROR_EMAIL_USED:       "邮箱已存在",
	ERROR_EMAIL_NIL:        "邮箱为空",
	ERROR_EMAIL_NO:         "发送验证码失败",
	//文章
	ERROR_ART_NOT_EXIST: "话题不存在",
	//分类

	ERROR_CATENAME_USED:  "分类已存在",
	ERROR_CATE_NOT_EXIST: "分类不存在",
}

func GetErrMsg(code int) string {
	return Codemsg[code]
}
