// 舍弃
package middleware

import (
	"bbs-go/utils"
	"bbs-go/utils/errmsg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

var JwKey = []byte(utils.JwKey)

type MyClaims struct {
	Userid uint `json:"userid`

	jwt.StandardClaims
}

var code int

// 生成token
func SetToken(userid uint, password string) (string, int) {
	expireTime := time.Now().Add(18 * time.Hour)
	SetClaims := MyClaims{userid,

		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "ginblog",
		},
	}
	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaims)
	token, err := reqClaim.SignedString(JwKey)
	if err != nil {
		return "", errmsg.ERROR
	}
	return token, errmsg.SUCCESS
}

// 验证token
func CheckToken(token string) (*MyClaims, int) {
	settoken, _ := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwKey, nil
	})
	if key, _ := settoken.Claims.(*MyClaims); settoken.Valid {
		return key, errmsg.SUCCESS
	} else {
		return nil, errmsg.ERROR
	}
}

// jwt中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHerder := c.Request.Header.Get("Authorization")
		code := errmsg.SUCCESS
		if tokenHerder == "" {
			code = errmsg.ERROR_TOKEN_EXIST
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"massage": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		checkToken := strings.SplitN(tokenHerder, " ", 2)
		if len(checkToken) != 2 && checkToken[0] != "Bearer" {
			code = errmsg.ERROR_TOKEN_TYPE_WRONG //不正确
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"massage": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		key, Tcode := CheckToken(checkToken[1])
		if Tcode == errmsg.ERROR {
			code = errmsg.ERROR_TOKEN_WRONG
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"massage": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		if time.Now().Unix() > key.ExpiresAt {
			code = errmsg.ERROR_TOKEN_RUNTIME //过期
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"massage": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		c.Set("userid", key.Userid)
		c.Next()
	}
}
func ParseToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	return userID, nil
}
