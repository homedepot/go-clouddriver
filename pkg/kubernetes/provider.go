package kubernetes

type Provider struct {
	Name          string              `json:"name" gorm:"primary_key"`
	Host          string              `json:"host"`
	CAData        string              `json:"caData" gorm:"size:8192"`
	BearerToken   string              `json:"bearerToken,omitempty" gorm:"size:2048"`
	TokenProvider string              `json:"tokenProvider,omitempty" gorm:"size:32;not null;default:'google'"`
	Permissions   ProviderPermissions `json:"permissions" gorm:"-"`
}

type ProviderPermissions struct {
	Read  []string `json:"read" gorm:"-"`
	Write []string `json:"write" gorm:"-"`
}

func (Provider) TableName() string {
	return "kubernetes_providers"
}
