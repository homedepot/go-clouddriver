package clouddriver

type kubernetes_clusters struct {
	id string `json:"-" gorm:"primary_key"`
	account_name string
	kind  string
	name string
	namespace string
	spinnaker_app string
	cluster string
}

func (kubernetes_clusters) TableName() string {
	return "kubernetes_clusters"
}