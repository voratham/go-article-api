package middleware

import (
	"article-api/config"
	"article-api/models"
	"errors"
	"log"
	"os"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var indentityKey string = "sub"

type login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func Authenticate() *jwt.GinJWTMiddleware {

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{

		Key: []byte(os.Getenv("SECRET_KEY")),

		IdentityKey: indentityKey,

		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",
		IdentityHandler: func(c *gin.Context) interface{} {
			var user models.User
			claims := jwt.ExtractClaims(c)
			id := claims[indentityKey]

			db := config.GetDB()
			err := db.First(&user, uint(id.(float64))).Error

			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Fatal("error middleware authenticate ", err)
				return nil
			}

			return &user
		},

		Authenticator: func(c *gin.Context) (interface{}, error) {
			var data login
			var user models.User

			if err := c.ShouldBindJSON(&data); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			db := config.GetDB()
			err := db.Where("email = ?", data.Email).First(&user).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, jwt.ErrFailedAuthentication
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return &user, nil

		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				claims := jwt.MapClaims{
					indentityKey: v.ID,
				}
				return claims
			}

			return jwt.MapClaims{}
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"error": message})
		},
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware
}
