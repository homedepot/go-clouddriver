package clouddriver

type kubernetes_clusters struct {
	account_name string
	cluster string
	id string `json:"-" gorm:"primary_key"`
	kind  string
	name string
	namespace string
	spinnaker_app string

}

func (kubernetes_clusters) TableName() string {
	return "kubernetes_clusters"
}