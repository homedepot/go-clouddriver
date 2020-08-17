package kubernetes

type Provider struct {
	Name   string `json:"name" gorm:"primary_key"`
	Host   string `json:"host"`
	CAData string `json:"caData" gorm:"size:2048"`
}

func (Provider) TableName() string {
	return "provider_kubernetes"
}
