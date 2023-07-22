package model

import (
	"bbs-go/utils/errmsg"
	"fmt"
	"github.com/go-gomail/gomail"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// 生成验证码
func GenerateVerificationCode(email string) int {
	rand.Seed(time.Now().UnixNano())
	var code strings.Builder
	for i := 0; i < 6; i++ {
		digit := rand.Intn(10) // 生成0到9的随机数
		code.WriteString(strconv.Itoa(digit))
	}
	errr := storeVerificationCode(email, code.String())
	if errr != errmsg.SUCCESS {
		return errmsg.ERROR
	}
	errr = sendVerificationEmail(email, code.String())
	if errr != errmsg.SUCCESS {
		return errmsg.ERROR
	}
	verificationCodes[email] = code.String()
	return errmsg.SUCCESS
}
func sendVerificationEmail(email, code string) int {
	// 创建邮件内容
	subject := "邮箱验证"
	body := fmt.Sprintf("您的验证码是：%s，请在验证页面输入该验证码完成邮箱验证。", code)

	// 使用邮件发送库发送邮件
	m := gomail.NewMessage()
	m.SetHeader("From", "m15237569469@163.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.163.com", 25, "m15237569469@163.com", "ITHPYKXRGURLLJEP")

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return errmsg.ERROR

	}

	return errmsg.SUCCESS
}

var verificationCodes = make(map[string]string)

func storeVerificationCode(email, code string) int {
	verificationCodes[email] = code
	if verificationCodes == nil {
		return errmsg.ERROR
	}
	fmt.Println("456", verificationCodes)
	return errmsg.SUCCESS
}

func VerifyCode(email, code string) int {

	storedCode, ok := verificationCodes[email]
	if !ok {
		return errmsg.ERROR
	}
	delete(verificationCodes, email) // 验证码使用后从内存中删除，确保一次性验证

	if code != storedCode {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
