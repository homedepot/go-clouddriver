package kubernetes

type Cluster struct {
	AccountName  string
	Cluster      string
	Id           string `json:"-" gorm:"primary_key"`
	Kind         string
	Name         string
	Namespace    string
	SpinnakerApp string
}

func (Cluster) TableName() string {
	return "kubernetes_clusters"
}
