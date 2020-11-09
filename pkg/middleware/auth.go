package middleware

import (
	"fmt"
	"net/http"

	clouddriver "github.com/billiford/go-clouddriver/pkg"
	"github.com/billiford/go-clouddriver/pkg/fiat"
	"github.com/billiford/go-clouddriver/pkg/http/core"
	"github.com/gin-gonic/gin"
)

const (
	headerSpinnakerUser        = `X-Spinnaker-User`
	headerSpinnakerApplication = `X-Spinnaker-Application`
)

//authApplication takes a list of permissions
//authAccount takes a list of accounts

func AuthApplication(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetHeader(headerSpinnakerUser)
		app := c.GetHeader(headerSpinnakerApplication)

		if user == "" || app == "" {
			c.Next()
			return
		}

		fiatClient := fiat.Instance(c)
		authResp, err := fiatClient.Authorize(user)
		if err != nil {
			clouddriver.WriteError(c, http.StatusUnauthorized, err)
			return
		}

		applicationsAuth := authResp.Applications
		for _, auth := range applicationsAuth {
			if auth.Name == app {
				for _, p := range permissions {
					found := find(auth.Authorizations, p)
					if !found {
						clouddriver.WriteError(c, http.StatusForbidden, fmt.Errorf("Access denied to application %s - required authorization: %s", app, p))
						return
					}
				}
			}
		}
		c.Next()
	}
}

func AuthAccount(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetHeader(headerSpinnakerUser)
		account := c.Param("account")
		fiatClient := fiat.Instance(c)

		if user == "" || account == "" {
			c.Next()
			return
		}

		authResp, err := fiatClient.Authorize(user)
		if err != nil {
			clouddriver.WriteError(c, http.StatusUnauthorized, err)
			return
		}

		accountsAuth := authResp.Accounts

		for _, auth := range accountsAuth {
			if auth.Name == account {
				for _, p := range permissions {
					found := find(auth.Authorizations, p)
					if !found {
						clouddriver.WriteError(c, http.StatusForbidden, fmt.Errorf("Access denied to account %s - required authorization: %s", account, p))
						return
					}
				}
			}
		}
		c.Next()
	}
}

func PostFilterAuthorizedApplications(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			return
		}

		allApps := core.Applications{}
		allApps = c.MustGet(core.KeyAllApplications).(core.Applications)

		user := c.GetHeader(headerSpinnakerUser)
		if user == "" {
			c.JSON(http.StatusOK, allApps)
			return
		}

		fiatClient := fiat.Instance(c)
		authResp, err := fiatClient.Authorize(user)
		if err != nil {
			clouddriver.WriteError(c, http.StatusUnauthorized, err)
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
