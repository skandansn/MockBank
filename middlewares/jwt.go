package middlewares

import (
	"github.com/skandansn/webDevBankBackend/routes"
	"github.com/skandansn/webDevBankBackend/utils/token"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, userId, userTypeId, err := token.ExtractUserTypeUserTypeIdAndUserId(c)
		currentRoute := c.Request.URL.Path

		if strings.HasPrefix(currentRoute, "/api") {
			currentRoute = currentRoute[4:]
		}

		supportedTiersForCurrentRoute := map[string]bool{}
		isCurrentRoutePublic := false

		for _, route := range routes.Routes {
			if matchRoute(currentRoute, route.Path) && route.Method == c.Request.Method {
				supportedTiersForCurrentRoute = route.Tiers
				isCurrentRoutePublic = route.Tiers["public"]
				break
			}
		}

		log.Println("currentRoute: ", currentRoute)
		log.Println("userType: ", userType)
		log.Println("supportedTiersForCurrentRoute: ", supportedTiersForCurrentRoute)

		if userType == "" {
			userType = "public"
		}

		if err != nil && !isCurrentRoutePublic {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		if isCurrentRoutePublic {
			c.Next()
			return
		}

		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		c.Set("userType", userType)
		c.Set("userId", userId)
		c.Set("userTypeId", userTypeId)

		if supportedTiersForCurrentRoute[userType] {
			c.Next()
			return
		}

		c.String(404, "Not found")
		c.Abort()
	}
}

func matchRoute(currentRoute string, definedRoute string) bool {
	currentRouteSegments := strings.Split(currentRoute, "/")
	definedRouteSegments := strings.Split(definedRoute, "/")

	if len(currentRouteSegments) != len(definedRouteSegments) {
		return false
	}

	for i := range currentRouteSegments {
		if !strings.HasPrefix(definedRouteSegments[i], ":") && currentRouteSegments[i] != definedRouteSegments[i] {
			return false
		}
	}

	return true
}
