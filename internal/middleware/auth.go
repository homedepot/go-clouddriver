package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/homedepot/go-clouddriver/internal/api/core"
	"github.com/homedepot/go-clouddriver/internal/api/core/kubernetes"
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
			c.Abort()

			return
		}

		applicationsAuth := authResp.Applications
		for _, auth := range applicationsAuth {
			if auth.Name == app {
				for _, p := range permissions {
					found := find(auth.Authorizations, p)
					if !found {
						clouddriver.Error(c, http.StatusForbidden, fmt.Errorf("Access denied to application %s - required authorization: %s", app, p))
						c.Abort()

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
			c.Abort()

			return
		}

		accountsAuth := authResp.Accounts

		for _, auth := range accountsAuth {
			if auth.Name == account {
				for _, p := range permissions {
					found := find(auth.Authorizations, p)
					if !found {
						clouddriver.Error(c, http.StatusForbidden, fmt.Errorf("Access denied to account %s - required authorization: %s", account, p))
						c.Abort()

						return
					}
				}
			}
		}

		c.Next()
	}
}

func (cc *Controller) AuthOps(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetHeader(headerSpinnakerUser)
		if user == "" {
			c.Next()
			return
		}

		// Get account(s) from payload
		ko := kubernetes.Operations{}

		if err := c.ShouldBindBodyWith(&ko, binding.JSON); err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			c.Abort()

			return
		}

		accounts := []string{}

		// Loop through each request in the kubernetes operations and perform
		// each requested action.
		for _, req := range ko {
			if req.DeployManifest != nil {
				accounts = appendAccount(accounts, req.DeployManifest.Account)
			}

			if req.DeleteManifest != nil {
				accounts = appendAccount(accounts, req.DeleteManifest.Account)
			}

			if req.DisableManifest != nil {
				accounts = appendAccount(accounts, req.DisableManifest.Account)
			}

			if req.ScaleManifest != nil {
				accounts = appendAccount(accounts, req.ScaleManifest.Account)
			}

			if req.CleanupArtifacts != nil {
				accounts = appendAccount(accounts, req.CleanupArtifacts.Account)
			}

			if req.RollingRestartManifest != nil {
				accounts = appendAccount(accounts, req.RollingRestartManifest.Account)
			}

			if req.RunJob != nil {
				accounts = appendAccount(accounts, req.RunJob.Account)
			}

			if req.UndoRolloutManifest != nil {
				accounts = appendAccount(accounts, req.UndoRolloutManifest.Account)
			}

			if req.PatchManifest != nil {
				accounts = appendAccount(accounts, req.PatchManifest.Account)
			}
		}

		if len(accounts) == 0 {
			c.Next()
			return
		}

		authResp, err := cc.FiatClient.Authorize(user)
		if err != nil {
			clouddriver.Error(c, http.StatusUnauthorized, err)
			c.Abort()

			return
		}

		// For each account in the request, verify user have required permissions to it.
		for _, account := range accounts {
			for _, p := range permissions {
				permitted := false
				for _, auth := range authResp.Accounts {
					if auth.Name == account {
						// User doesn't have required permission to account.
						permitted = find(auth.Authorizations, p)
						if !permitted {
							clouddriver.Error(c, http.StatusForbidden, fmt.Errorf("Access denied to account %s - required authorization: %s", account, p))
							c.Abort()

							return
						}
					}
				}
				// User doesn't have any permission to account.
				if !permitted {
					clouddriver.Error(c, http.StatusForbidden, fmt.Errorf("Access denied to account %s - required authorization: %s", account, p))
					c.Abort()

					return
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

func appendAccount(accounts []string, account string) []string {
	if account != "" && !find(accounts, account) {
		accounts = append(accounts, account)
	}

	return accounts
}
