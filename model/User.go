package model

import (
	"bbs-go/utils/errmsg"
	"encoding/base64"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
	"log"
	"time"
)

var Instance *Config

type User struct {
	gorm.Model
	Username      string    `gorm:"not null;comment:用户名"json:"username"`
	Password      string    `gorm:"not null;comment:密码"json:"password"`
	Email         string    `gorm:"not null;comment:邮箱"json:"email"`
	LastLoginTime time.Time `gorm:"comment:最后登录时间" json:"last_login_time"` // 最后登录时间
	Role          int       `gorm:"type:int;comment:权限等级" json:"role"`     //0 代表管理员 1 待变普通用户
	//Nickname      string    `gorm:"comment:昵称" json:"nickname"`            // 昵称
	//Bio           string    `gorm:"comment:用户简介" json:"bio"`               // 用户简介
	//Avatar        string    `gorm:"comment:个人头像" json:"avatar"`            // 个人头像
}

type Config struct {
	Jwt struct {
		SignKey       string `yaml:"SignKey"`
		ExpireSeconds int    `yaml:"ExpireSeconds"`
		Issuer        string `yaml:"Issuer"`
	} `yaml:"jwt`
}
type Message struct { //消息
	gorm.Model
	SenderID   uint   `gorm:"comment:发送者ID" json:"sender_id"`   // 发送者ID
	ReceiverID uint   `gorm:"comment:接收者ID" json:"receiver_id"` // 接收者ID
	Content    string `gorm:"comment:内容" json:"content"`        // 内容

}

// 查询用户是否存在
func CheckUser(username, email string) int {
	var users User
	db.Select("ID").Where("username = ?", username).First(&users)

	if users.ID > 0 {

		return errmsg.ERROR_USERNAME_USED //1001
	}
	db.Select("ID").Where("email = ?", email).First(&users)
	if users.ID > 0 {
		return errmsg.ERROR_EMAIL_USED
	}
	return errmsg.SUCCSE //可用
}
func CheckUserid(username string) (int, uint) {
	var users User
	db.Select("ID").Where("username = ?", username).First(&users)

	if users.ID < 0 {

		return errmsg.ERROR_USERNAME_USED, 0 //1001
	}

	return errmsg.SUCCSE, users.ID //可用
}
func Checkenail(email string) (int, string) {
	var user User
	db.Select("ID").Where("enail = ?", email).First(&user)

	if user.ID < 0 {

		return errmsg.ERROR_USERNAME_USED, "" //1001
	}

	return errmsg.SUCCSE, user.Email //可用
}

// 根据id查询用户是否存在
func CheckUserId(id int) int {
	var users User
	db.Select("id").Where("id = ?", id).First(&users)
	if users.ID > 0 {
		return errmsg.ERROR_USERNAME_USED //1001
	}
	return errmsg.ERROR_USER_NOT_EXIST
}

// 密码的加密
// 钩子函数实现不调用就能使用函数
//
//	func (u *User) BeforeSave() {
//		u.Password = ScryptPw(u.Password)
//	}
func ScryptPw(password string) string {
	const KeyLen = 10
	salt := make([]byte, 8)
	salt = []byte{12, 45, 78, 35, 78, 45, 78, 78}

	HashPw, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KeyLen)
	if err != nil {
		log.Fatal(err)
	}
	fpw := base64.StdEncoding.EncodeToString(HashPw)
	return fpw
}

// 新增用户
func CreateUser(data *User) int {
	//密码写入数据库前进行哈希加密
	data.Password = ScryptPw(data.Password)
	//data.BeforeSave()
	var profile Profile
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	profile.Username = data.Username
	profile.Email = data.Email
	err = db.Create(&profile).Error
	if err != nil && data.ID != profile.ID {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 登陆验证
func CheckLOgin(username string, password string) int {
	var user User
	var un = username
	db.Where("username = ?", un).First(&user)
	db.Where("email = ?", un).First(&user)
	if user.ID == 0 {
		return errmsg.ERROR_USER_NOT_EXIST

	}
	if ScryptPw(password) != user.Password {
		return errmsg.ERROR_PASSWORD_WEONG
	}
	if user.Role != 0 {
		return errmsg.ERROR_USER_NO_RIGHT
	}
	return errmsg.SUCCSE
}
