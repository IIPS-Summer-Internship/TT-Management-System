package controllers

import (
	"net/http"
	"tms-server/config"
	"tms-server/models"
	"tms-server/utils"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong! TMS-server is up"})
}

func Login(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := utils.AuthenticateUser(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// load user with a role
	if err := config.DB.Preload("Role").First(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load user role"})
		return
	}

	token, err := utils.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.SetCookie(
		"auth_token", token,
		int(utils.TokenExpiry.Seconds()), // expires in 7 days (604800 seconds)
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Login Successful",
		"username": user.Username,
		"role":     user.Role.Name,
	})
}

// vaildate session API
func ValidateSession(c *gin.Context) {
	//check cookie
	tokenString, err := c.Cookie("auth_token")
	if err != nil {
		c.JSON(404, gin.H{"error": "Cookie not present"})
		return
	}
	//validate token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid session"})
		return
	}
	//getting user with role from db
	var user models.User
	if err := config.DB.Preload("Role").Where("username = ?", claims.Username).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}
	/*check the roles - superadmin-1, admin-2, faculty-3*/
	if user.Role.LevelPriority > 3 {
		c.JSON(403, gin.H{"error": "Insufficient role - faculty required"})
		return
	}

	//session is vaild
	c.JSON(200, gin.H{"message": "Session is valid", "role": user.Role.Name})

}

func Logout(c *gin.Context) {
	c.SetCookie(
		"auth_token",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
