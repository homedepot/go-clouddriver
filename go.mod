module github.com/homedepot/go-clouddriver

go 1.14

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/billiford/go-clouddriver v0.6.5
	github.com/docker/docker-credential-helpers v0.6.3 // indirect
	github.com/fatih/color v1.9.0
	github.com/gdexlab/go-render v1.0.1 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/go-git/go-git/v5 v5.2.0
	github.com/go-playground/assert/v2 v2.0.1 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/gddo v0.0.0-20200715224205-051695c33a3f // indirect
	github.com/google/go-github v17.0.0+incompatible // indirect
	github.com/google/go-github/v32 v32.1.0
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.5.1
	github.com/gregjones/httpcache v0.0.0-20180305231024-9cad4c3443a7
	github.com/iancoleman/strcase v0.1.2
	github.com/jinzhu/gorm v1.9.16
	github.com/jonboulle/clockwork v0.1.0
	github.com/lib/pq v1.2.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.0
	github.com/maxbrunsfeld/counterfeiter/v6 v6.2.3 // indirect
	github.com/mcuadros/go-gin-prometheus v0.1.0
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/peterbourgon/diskv v2.0.1+incompatible
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/stretchr/testify v1.5.1 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	google.golang.org/api v0.15.0 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1 // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/cli-runtime v0.19.2
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v0.0.0-20181102134211-b9b56d5dfc92 // indirect
	k8s.io/klog/v2 v2.3.0
	k8s.io/kube-openapi v0.0.0-20200805222855-6aeccd4b50c6
	k8s.io/kubectl v0.19.2
	sigs.k8s.io/structured-merge-diff/v3 v3.0.0-20200116222232-67a7b8c61874 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.19.2
