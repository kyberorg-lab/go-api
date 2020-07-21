package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/kyberorg/go-api/app"
	"net"
	"strings"
)

func GetUniqueUserAgent(context *gin.Context) string {
	return getClientIP(context) + app.UAIPDelimiter + getUserAgent(context)
}

func ExtractToken(context *gin.Context) string {
	bearerToken := context.GetHeader("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func getUserAgent(context *gin.Context) string {
	userAgent := context.GetHeader("User-Agent")
	if userAgent == "" {
		userAgent = app.UAUserAgentUnknown
	}
	return userAgent
}

func getClientIP(context *gin.Context) string {
	ip := context.ClientIP()
	if ip == "" {
		ip = app.UAIPUnknown
	}
	parsedIp := net.ParseIP(ip)
	if parsedIp == nil {
		return app.UAIPUnknown
	}
	ipv4 := parsedIp.To4()
	if ipv4 == nil {
		ipv6 := parsedIp.String()
		return ipv6
	}
	return ipv4.String()
}
