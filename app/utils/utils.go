package utils

import (
	"github.com/gin-gonic/gin"
	"go-rest/app"
	"net"
)

type ErrJson struct {
	Err   string `json:"err"`
	Error error  `json:"error" binding: "omitempty"`
}

func ErrorJson(err string) ErrJson {
	return ErrJson{err, nil}
}

func ErrorJsonWithError(message string, err error) ErrJson {
	return ErrJson{
		Err:   message,
		Error: err,
	}
}

func GetUserAgent(context *gin.Context) string {
	userAgent := context.GetHeader("User-Agent")
	if userAgent == "" {
		userAgent = app.UserAgentUnknown
	}
	return userAgent
}

func GetClientIP(context *gin.Context) string {
	ip := context.ClientIP()
	if ip == "" {
		ip = app.IPUnknown
	}
	parsedIp := net.ParseIP(ip)
	if parsedIp == nil {
		return app.IPUnknown
	}
	ipv4 := parsedIp.To4()
	if ipv4 == nil {
		ipv6 := parsedIp.String()
		return ipv6
	}
	return ipv4.String()
}
