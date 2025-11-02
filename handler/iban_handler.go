package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tapsilat/iban.im/config"
	"github.com/tapsilat/iban.im/model"
)

// GetIbanByHandles retrieves an IBAN by user handle and IBAN handle
func GetIbanByHandles(c *gin.Context) {
	userHandle := c.Param("userHandle")
	ibanHandle := c.Param("ibanHandle")

	// Find user by handle
	var user model.User
	if err := config.DB.Where("handle = ?", userHandle).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Find IBAN by handle and owner
	var iban model.Iban
	if err := config.DB.Where("owner_id = ? AND handle = ? AND is_private = false", user.UserID, ibanHandle).First(&iban).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "IBAN not found or is private",
		})
		return
	}

	// Return IBAN information
	c.JSON(http.StatusOK, gin.H{
		"userHandle": userHandle,
		"ibanHandle": ibanHandle,
		"iban":       iban.Text,
		"description": iban.Description,
	})
}

// RenderIbanPage renders a simple HTML page displaying the IBAN or returns JSON based on Accept header
func RenderIbanPage(c *gin.Context) {
	userHandle := c.Param("userHandle")
	ibanHandle := c.Param("ibanHandle")

	// Find user by handle
	var user model.User
	if err := config.DB.Where("handle = ?", userHandle).First(&user).Error; err != nil {
		// Check if client wants JSON
		if c.GetHeader("Accept") == "application/json" || c.Query("format") == "json" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
		} else {
			c.HTML(http.StatusNotFound, "error.tmpl.html", gin.H{
				"error": "User not found",
			})
		}
		return
	}

	// Find IBAN by handle and owner
	var iban model.Iban
	if err := config.DB.Where("owner_id = ? AND handle = ? AND is_private = false", user.UserID, ibanHandle).First(&iban).Error; err != nil {
		// Check if client wants JSON
		if c.GetHeader("Accept") == "application/json" || c.Query("format") == "json" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "IBAN not found or is private",
			})
		} else {
			c.HTML(http.StatusNotFound, "error.tmpl.html", gin.H{
				"error": "IBAN not found or is private",
			})
		}
		return
	}

	// Check if client wants JSON response
	if c.GetHeader("Accept") == "application/json" || c.Query("format") == "json" {
		c.JSON(http.StatusOK, gin.H{
			"userHandle":  userHandle,
			"ibanHandle":  ibanHandle,
			"iban":        iban.Text,
			"description": iban.Description,
			"firstName":   user.FirstName,
			"lastName":    user.LastName,
		})
		return
	}

	// Render the IBAN page
	c.HTML(http.StatusOK, "iban.tmpl.html", gin.H{
		"userHandle":  userHandle,
		"ibanHandle":  ibanHandle,
		"iban":        iban.Text,
		"description": iban.Description,
		"firstName":   user.FirstName,
		"lastName":    user.LastName,
	})
}

// IsValidRoute checks if the route matches the pattern /:userHandle/:ibanHandle
// This helps distinguish from other routes like /assets, /api, /graph, etc.
func IsValidRoute(path string) bool {
	// Exclude known system routes
	excludedPrefixes := []string{
		"/assets",
		"/api",
		"/auth",
		"/graph",
		"/favicon",
		"/robots",
		"/sitemap",
	}
	
	for _, prefix := range excludedPrefixes {
		if strings.HasPrefix(path, prefix) {
			return false
		}
	}
	
	// Check if path has exactly 2 segments (/:handle/:handle)
	segments := 0
	for i := 1; i < len(path); i++ {
		if path[i] == '/' {
			segments++
		}
	}
	
	return segments == 1 && len(path) > 2
}

// HandleIbanRoute is a middleware that checks if the route matches /:userHandle/:ibanHandle
// and delegates to RenderIbanPage if it does
func HandleIbanRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		
		// Only handle routes that match the pattern
		if IsValidRoute(path) {
			// Try to extract user and iban handles
			userHandle := c.Param("userHandle")
			ibanHandle := c.Param("ibanHandle")
			
			if userHandle != "" && ibanHandle != "" {
				RenderIbanPage(c)
				c.Abort() // Stop processing other handlers
				return
			}
		}
		
		c.Next()
	}
}
