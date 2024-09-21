package routes

import (
	"log"
	"net/http"
	"time"

	"ssl-checker/auth"
	"ssl-checker/cache"
	domains_utils "ssl-checker/domains"
	"ssl-checker/utils"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {

	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		sessionID, err := c.Cookie(auth.SessionCookieName)
		if err != nil || !auth.ValidateSession(sessionID) {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
		c.Redirect(http.StatusSeeOther, "/home")
	})

	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil) // Отображение страницы регистрации
	})

	r.POST("/register", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		if !auth.ValidateUsername(username) {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"Error": "The username must contain only letters, numbers, periods, dashes or underscores.",
			})
			return
		}

		if len(username) < 3 || len(password) < 6 {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"Error": "The username must be at least 3 characters and the password must be at least 6 characters.",
			})
			return
		}

		if auth.UserExists(username) {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"Error": "A user with this name already exists.",
			})
			return
		}

		auth.Users[username] = auth.User{
			Username: username,
			Password: string(auth.HashPassword(password)),
			Role:     "user",
			Status:   "pending",
		}

		auth.SaveUsers()

		c.HTML(http.StatusOK, "login.html", gin.H{
			"Success": "Registration is successful. Wait for administrator's confirmation..",
		})
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		browser, osInfo := utils.ParseUserAgent(c.Request.UserAgent())

		if auth.CheckPassword(username, password) {
			sessionID := auth.GenerateSessionID()
			auth.Sessions[sessionID] = auth.Session{
				Username:  username,
				SessionID: sessionID,
				IP:        c.ClientIP(),
				Device:    browser + " on " + osInfo,
				Expiry:    time.Now().Add(30 * time.Minute)}
			http.SetCookie(c.Writer, &http.Cookie{
				Name:    auth.SessionCookieName,
				Value:   sessionID,
				Expires: time.Now().Add(30 * time.Minute),
			})
			c.Redirect(http.StatusSeeOther, "/home")
		} else {
			log.Println("Failed login attempt for username:", username)
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"Error": "Invalid username or password",
			})
		}
	})

	r.POST("/logout", func(c *gin.Context) {
		sessionID, err := c.Cookie(auth.SessionCookieName)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		delete(auth.Sessions, sessionID)

		http.SetCookie(c.Writer, &http.Cookie{
			Name:    auth.SessionCookieName,
			Value:   "",
			Expires: time.Now().Add(-time.Hour),
			Path:    "/",
		})

		c.Redirect(http.StatusSeeOther, "/login")
	})

	r.GET("/home", func(c *gin.Context) {
		sessionID, err := c.Cookie(auth.SessionCookieName)
		if err != nil || !auth.ValidateSession(sessionID) {
			log.Println("Session invalid or expired, redirecting to login")
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		isAdmin := auth.Users[auth.Sessions[sessionID].Username].Role == "admin"
		userSession := auth.Sessions[sessionID]
		if auth.Users[userSession.Username].Status != "registered" {
			c.HTML(http.StatusForbidden, "login.html", gin.H{
				"Error": "Your account is awaiting confirmation by the administrator.",
			})
			return
		}

		currentUser := auth.Sessions[sessionID].Username
		stats, err := utils.CalculateCertificateStats()
		if err != nil {
			log.Println("message: Error calculating certificate stats.")
			return
		}

		domains, err := cache.LoadCache()
		if err != nil {
			log.Println("Error loading cache:", err)
			c.HTML(http.StatusInternalServerError, "home.html", gin.H{
				"currentUser":      currentUser,
				"isAdmin":          isAdmin,
				"showActionColumn": isAdmin,
				"Error":            "Error loading domain data",
			})
			return
		}

		c.HTML(http.StatusOK, "home.html", gin.H{
			"currentUser":              currentUser,
			"isAdmin":                  isAdmin,
			"showActionColumn":         isAdmin,
			"domains":                  domains,
			"TotalCertificates":        stats.TotalCertificates,
			"ValidCertificates":        stats.ValidCertificates,
			"ExpiringSoonCertificates": stats.ExpiringSoonCertificates,
			"CriticalCertificates":     stats.CriticalCertificates,
			"ErrorCertificates":        stats.ErrorCertificates,
		})
	})

	r.POST("/home/add-domain", func(c *gin.Context) {
		sessionID, err := c.Cookie(auth.SessionCookieName)
		if err != nil || !auth.ValidateSession(sessionID) || auth.Users[auth.Sessions[sessionID].Username].Role != "admin" {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		isAdmin := auth.Users[auth.Sessions[sessionID].Username].Role == "admin"

		domainName := c.PostForm("domain-name")
		if domainName == "" {
			c.Redirect(http.StatusSeeOther, "/home")
			return
		}

		domainName = utils.ExtractDomain(domainName)
		log.Println("Extracted domain:", domainName)

		domains, err := domains_utils.LoadDomains()

		for _, d := range domains {
			if d == domainName {
				c.HTML(http.StatusOK, "home.html", gin.H{
					"isAdmin":          isAdmin,
					"showActionColumn": isAdmin,
					"domains":          domains,
					"Error":            "Domain already exists",
				})
				c.Redirect(http.StatusSeeOther, "/home")
				return
			}

		}

		domains = append(domains, domainName)
		err = domains_utils.SaveDomains(domains)
		if err != nil {
			log.Println("Error saving domains:", err)
			c.Redirect(http.StatusSeeOther, "/home")
			return
		}

		log.Println("Domain added successfully:", domainName)

		domainData, err := domains_utils.UpdateDomainTable(domains)
		if err != nil {
			log.Println("Error updating domain table:", err)
		} else {
			err = cache.SaveCache(domainData)
			if err != nil {
				log.Println("Error saving cache:", err)
			}
		}

		c.Redirect(http.StatusSeeOther, "/home")
	})

	r.POST("/home/del-domain", func(c *gin.Context) {
		sessionID, err := c.Cookie(auth.SessionCookieName)
		if err != nil || !auth.ValidateSession(sessionID) || auth.Users[auth.Sessions[sessionID].Username].Role != "admin" {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		domainName := c.PostForm("domain")
		if domainName == "" {
			c.Redirect(http.StatusSeeOther, "/home")
			return
		}

		domains, err := domains_utils.LoadDomains()
		if err != nil {
			log.Println("Error loading domains:", err)
			c.Redirect(http.StatusSeeOther, "/home")
			return
		}

		index := -1
		for i, d := range domains {
			if d == domainName {
				index = i
				break
			}
		}

		if index != -1 {
			domains = append(domains[:index], domains[index+1:]...)
			err = domains_utils.SaveDomains(domains)
			if err != nil {
				log.Println("Error saving domains:", err)
			}

			domainData, err := domains_utils.UpdateDomainTable(domains)
			if err != nil {
				log.Println("Error updating domain table:", err)
			} else {
				err = cache.SaveCache(domainData)
				if err != nil {
					log.Println("Error saving cache:", err)
				}
			}
		}

		c.Redirect(http.StatusSeeOther, "/home")
	})

	r.GET("/admin", func(c *gin.Context) {
		sessionID, err := c.Cookie(auth.SessionCookieName)
		if err != nil || !auth.ValidateSession(sessionID) || auth.Users[auth.Sessions[sessionID].Username].Role != "admin" {
			log.Println("Unauthorized access to /admin")
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		activeSessions := []auth.Session{}
		for id, session := range auth.Sessions {
			if time.Now().Before(session.Expiry) {
				activeSessions = append(activeSessions, session)
			} else {
				delete(auth.Sessions, id)
			}
		}

		currentUser := auth.Sessions[sessionID].Username
		userList := []auth.User{}
		for _, user := range auth.Users {
			userList = append(userList, user)
		}

		c.HTML(http.StatusOK, "admin.html", gin.H{
			"isAdmin":     true,
			"Sessions":    activeSessions,
			"currentUser": currentUser,
			"users":       userList,
		})
	})

	r.POST("/admin/approve-user", func(c *gin.Context) {
		username := c.PostForm("username")

		user, exists := auth.Users[username]
		if !exists {
			c.HTML(http.StatusBadRequest, "admin.html", gin.H{
				"Error": "User not found.",
			})
			return
		}

		if user.Status != "pending" {
			c.HTML(http.StatusBadRequest, "admin.html", gin.H{
				"Error": "The user has already been confirmed or has an incorrect status.",
			})
			return
		}

		user.Status = "registered"
		auth.Users[username] = user

		auth.SaveUsers()

		c.Redirect(http.StatusSeeOther, "/admin")
	})

	r.POST("/admin/end-session", func(c *gin.Context) {
		sessionID := c.PostForm("session_id")
		if sessionID == "" {
			c.Redirect(http.StatusSeeOther, "/admin")
			return
		}

		auth.EndSession(sessionID)

		c.Redirect(http.StatusSeeOther, "/admin")
	})

	r.POST("/admin/add-user", func(c *gin.Context) {
		var errorMessage string

		if !auth.CheckAdminSession(c) {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		username := c.PostForm("username")
		password := c.PostForm("password")
		role := c.PostForm("role")

		if !auth.ValidateUsername(username) {
			errorMessage = "Invalid username. Only letters, numbers, dashes, underscores, and dots are allowed."
		}

		if auth.UserExists(username) {
			errorMessage = "User already exists."
		}

		if errorMessage != "" {
			c.HTML(http.StatusOK, "admin.html", gin.H{
				"Error": errorMessage,
			})
			return
		}

		auth.Users[username] = auth.User{
			Username: username,
			Password: string(auth.HashPassword(password)),
			Role:     role,
		}
		auth.SaveUsers()
		c.Redirect(http.StatusSeeOther, "/admin")
	})

	r.POST("/admin/delete-user", func(c *gin.Context) {
		if !auth.CheckAdminSession(c) {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		username := c.PostForm("username")
		delete(auth.Users, username)
		auth.SaveUsers()
		c.Redirect(http.StatusSeeOther, "/admin")
	})

	r.POST("/admin/set-role", func(c *gin.Context) {
		if !auth.CheckAdminSession(c) {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		username := c.PostForm("username")
		role := c.PostForm("role")

		user, exists := auth.Users[username]
		if exists {
			user.Role = role
			auth.Users[username] = user
			auth.SaveUsers()
		}
		c.Redirect(http.StatusSeeOther, "/admin")
	})

	r.GET("/change-password", func(c *gin.Context) {
		sessionID, err := c.Cookie(auth.SessionCookieName)
		if err != nil || !auth.ValidateSession(sessionID) {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		currentUser := auth.Sessions[sessionID].Username

		c.HTML(http.StatusOK, "change-password.html", gin.H{
			"isAdmin":     true,
			"currentUser": currentUser,
		})
	})

	r.POST("/change-password", func(c *gin.Context) {
		sessionID, err := c.Cookie(auth.SessionCookieName)
		if err != nil || !auth.ValidateSession(sessionID) {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		username := auth.Sessions[sessionID].Username
		currentPassword := c.PostForm("current-password")
		newPassword := c.PostForm("new-password")

		if !auth.CheckPassword(username, currentPassword) {
			c.HTML(http.StatusBadRequest, "change-password.html", gin.H{
				"Error": "Current password is incorrect",
			})
			return
		}

		auth.Users[username] = auth.User{
			Username: username,
			Password: string(auth.HashPassword(newPassword)),
			Role:     auth.Users[username].Role,
		}
		auth.SaveUsers()
		c.Redirect(http.StatusSeeOther, "/home")
	})

	r.LoadHTMLGlob("templates/*")
}
