package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/internal/api/core"
	"github.com/homedepot/go-clouddriver/internal/fiat"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
)

const (
	headerSpinnakerUser        = `X-Spinnaker-User`
	headerSpinnakerApplication = `X-Spinnaker-Application`
)

func (cc *Controller) AuthApplication(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetHeader(headerSpinnakerUser)
		app := c.GetHeader(headerSpinnakerApplication)

		if user == "" || app == "" {
			c.Next()
			return
		}

		authResp, err := cc.FiatClient.Authorize(user)
		if err != nil {
			clouddriver.Error(c, http.StatusUnauthorized, err)
			return
		}

		applicationsAuth := authResp.Applications
		for _, auth := range applicationsAuth {
			if auth.Name == app {
				for _, p := range permissions {
					found := find(auth.Authorizations, p)
					if !found {
						clouddriver.Error(c, http.StatusForbidden, fmt.Errorf("Access denied to application %s - required authorization: %s", app, p))
						return
					}
				}
			}
		}

		c.Next()
	}
}

func (cc *Controller) AuthAccount(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetHeader(headerSpinnakerUser)
		account := c.Param("account")

		if user == "" || account == "" {
			c.Next()
			return
		}

		authResp, err := cc.FiatClient.Authorize(user)
		if err != nil {
			clouddriver.Error(c, http.StatusUnauthorized, err)
			return
		}

		accountsAuth := authResp.Accounts

		for _, auth := range accountsAuth {
			if auth.Name == account {
				for _, p := range permissions {
					found := find(auth.Authorizations, p)
					if !found {
						clouddriver.Error(c, http.StatusForbidden, fmt.Errorf("Access denied to account %s - required authorization: %s", account, p))
						return
					}
				}
			}
		}

		c.Next()
	}
}

func (cc *Controller) PostFilterAuthorizedApplications(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			return
		}

		allApps := c.MustGet(core.KeyAllApplications).(core.Applications)

		user := c.GetHeader(headerSpinnakerUser)
		if user == "" {
			c.JSON(http.StatusOK, allApps)
			return
		}

		authResp, err := cc.FiatClient.Authorize(user)
		if err != nil {
			clouddriver.Error(c, http.StatusUnauthorized, err)
			return
		}

		authorizedApps := authResp.Applications
		authorizedAppsMap := map[string]fiat.Application{}

		for _, app := range authorizedApps {
			authorizedAppsMap[app.Name] = app
		}

		filteredApps := FilterAuthorizedApps(authorizedAppsMap, allApps, permissions...)

		c.JSON(http.StatusOK, filteredApps)
	}
}

func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}

	return false
}

func FilterAuthorizedApps(authorizedAppsMap map[string]fiat.Application, allApps core.Applications, permissions ...string) []core.Application {
	filteredApps := []core.Application{}

	for _, app := range allApps {
		if authorizedApp, ok := authorizedAppsMap[app.Name]; ok {
			for _, p := range permissions {
				if ok := find(authorizedApp.Authorizations, p); ok {
					filteredApps = append(filteredApps, app)
				}
			}
		}
	}

	return filteredApps
}
