package main // import "github.com/tapsilat/iban.im

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/tapsilat/iban.im/config"
	"github.com/tapsilat/iban.im/handler"
	_ "github.com/tapsilat/iban.im/model"

	"github.com/tapsilat/iban.im/resolvers"
	"github.com/tapsilat/iban.im/schema"

	jwt "github.com/appleboy/gin-jwt/v2"

	"fmt"

	"github.com/gin-gonic/gin"
)

const identityKey = "UserID"

func main() {

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Next()
	})

	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	if sqlDB, err := config.DB.DB(); err == nil {
		defer sqlDB.Close()
	}

	context.Background()

	authMiddleware, err := handler.AuthMiddleware()

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	router.POST("/api/login", func(c *gin.Context) {
		authMiddleware.LoginHandler(c)
	})

	auth := router.Group("/auth")
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	router.GET("/graph", func(c *gin.Context) {
		c.HTML(http.StatusOK, "graph.tmpl.html", nil)
	})

	authMW := authMiddleware.MiddlewareFunc()

	router.POST("/graph", func(c *gin.Context) {
		ctx := c.Request.Context()

		if _, ok := c.Request.Header["Authorization"]; ok {
			authMW(c)

			claims := jwt.ExtractClaims(c)

			currentID, ok := claims[identityKey].(float64)
			if !ok {
				currentID = 0
			}
			ctx = context.WithValue(ctx, handler.ContextKey("UserID"), int(currentID))
		}

		var params struct {
			Query         string                 `json:"query"`
			OperationName string                 `json:"operationName"`
			Variables     map[string]interface{} `json:"variables"`
		}
		if err := json.NewDecoder(c.Request.Body).Decode(&params); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}
		schema := graphql.MustParseSchema(*schema.NewSchema(), &resolvers.Resolvers{}, opts...)

		response := schema.Exec(ctx, params.Query, params.OperationName, params.Variables)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(200, response)
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Config.App.Port), router))

}
