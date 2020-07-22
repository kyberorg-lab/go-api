package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/kyberorg/go-api/global"
	"net"
	"strings"
)

func GetUniqueUserAgent(context *gin.Context) string {
	return getClientIP(context) + global.UAIPDelimiter + getUserAgent(context)
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
		userAgent = global.UAUserAgentUnknown
	}
	return userAgent
}

func getClientIP(context *gin.Context) string {
	ip := context.ClientIP()
	if ip == "" {
		ip = global.UAIPUnknown
	}
	parsedIp := net.ParseIP(ip)
	if parsedIp == nil {
		return global.UAIPUnknown
	}
	ipv4 := parsedIp.To4()
	if ipv4 == nil {
		ipv6 := parsedIp.String()
		return ipv6
	}
	return ipv4.String()
}
