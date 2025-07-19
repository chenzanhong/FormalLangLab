package api

import (
	"backend/internal/domain"
	r_init "backend/internal/repository"
	"backend/pkg/middleware"

	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var loginRequest domain.LoginRequest
	err := c.BindJSON(&loginRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.LoginResponse{
			Code:    400,
			Message: "登录时数据解析失败",
			Error:   err.Error(),
		})
	}

	// 查询数据库，判断用户是否存在
	var user domain.User
	if err := r_init.DB.Where("name = ?", loginRequest.Name).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, domain.LoginResponse{
			Code:    500,
			Message: "登录时查询用户信息识别",
			Error:   err.Error(),
		})
	}

	// 验证密码
	if user.Password != loginRequest.Password {
		c.JSON(http.StatusUnauthorized, domain.LoginResponse{
			Code:    401,
			Message: "用户名或密码错误",
		})
	}

	// 生成 JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &middleware.Claims{
		Username: user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "生成 token 错误"})
		return
	}

	// 返回信息
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"name":  loginRequest.Name,
		"id":    user.ID,
		"role":  user.Role,
	})
}
