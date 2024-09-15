package auth

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Session struct {
	Username  string
	SessionID string
	IP        string
	Device    string
	Expiry    time.Time
}

var Sessions = make(map[string]Session)
var SessionCookieName = "session_id"

func GenerateSessionID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

func ValidateSession(sessionID string) bool {
	session, exists := Sessions[sessionID]
	if !exists || time.Now().After(session.Expiry) {
		delete(Sessions, sessionID)
		return false
	}
	return true
}

func CheckAdminSession(c *gin.Context) bool {
	sessionID, err := c.Cookie(SessionCookieName)
	if err != nil || !ValidateSession(sessionID) || Users[Sessions[sessionID].Username].Role != "admin" {
		return false
	}
	return true
}

func EndSession(sessionID string) {
	delete(Sessions, sessionID)
}

func GetActiveSessions() map[string]Session {
	activeSessions := make(map[string]Session)
	for sessionID, session := range Sessions {
		if time.Now().Before(session.Expiry) {
			activeSessions[sessionID] = session
		} else {
			delete(Sessions, sessionID)
		}
	}
	return activeSessions
}
