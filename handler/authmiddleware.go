package handler

import (
	"fmt"
	"time"

	"github.com/tapsilat/iban.im/config"
	"github.com/tapsilat/iban.im/model"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const identityKey = "UserID"

type login struct {
	Handle   string `form:"handle" json:"handle" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// the jwt middleware
func AuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       config.Config.App.Realm,
		Key:         []byte(config.Config.App.Key),
		Timeout:     time.Minute * time.Duration(config.Config.App.Timeout),
		MaxRefresh:  time.Minute * time.Duration(config.Config.App.MaxRefresh),
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			fmt.Println("inside payload func")
			// fmt.Printf("payload data: %+v\n",data)
			if v, ok := data.(*model.User); ok {
				// fmt.Println("inside v, ",v.Handle, v.UserID)

				return jwt.MapClaims{
					identityKey: v.UserID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			fmt.Println("inside identity handler")
			claims := jwt.ExtractClaims(c)
			// user, _ := c.Get(identityKey)
			// fmt.Printf("claims: %+v\n",claims)

			return &model.User{
				Handle: fmt.Sprintf("%f", claims[identityKey].(float64)),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			fmt.Println("inside Authenticator")
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", fmt.Errorf("bind error : %v ", jwt.ErrMissingLoginValues)
			}

			user := model.User{}

			config.DB.Where("email = ?", loginVals.Handle).First(&user)
			if user.UserID == 0 {
				return "", fmt.Errorf("database where error : %v ", jwt.ErrFailedAuthentication)
			}

			if !user.ComparePassword(loginVals.Password) {
				return "", fmt.Errorf("compare password error: %v ", jwt.ErrFailedAuthentication)
			}

			return &model.User{
				// Handle: (string)(user.UserID),
				UserID:    user.UserID,
				LastName:  user.LastName,
				FirstName: user.FirstName,
				// Handle: user.Email,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			fmt.Println("inside Authorizator")
			// fmt.Printf("data: %+v\n",data)

			if _, ok := data.(*model.User); ok {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			// fmt.Println("inside unauthorized")
			c.JSON(code, gin.H{
				"code ":   code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}
