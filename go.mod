module github.com/homedepot/go-clouddriver

go 1.14

require (
	cloud.google.com/go/storage v1.17.0
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/fatih/color v1.9.0
	github.com/fsouza/fake-gcs-server v1.30.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/go-github/v32 v32.1.0
	github.com/google/uuid v1.1.2
	github.com/googleapis/gnostic v0.5.1
	github.com/gregjones/httpcache v0.0.0-20180305231024-9cad4c3443a7
	github.com/iancoleman/strcase v0.1.2
	github.com/jonboulle/clockwork v0.1.0
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.2.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.8
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.18.1
	github.com/peterbourgon/diskv v2.0.1+incompatible
	github.com/zsais/go-gin-prometheus v0.1.1-0.20200217150448-2199a42d96c1
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	google.golang.org/api v0.57.0
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.1.2 // indirect
	gorm.io/driver/sqlite v1.1.6 // indirect
	gorm.io/gorm v1.21.16 // indirect
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/cli-runtime v0.19.2
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog/v2 v2.3.0
	k8s.io/kube-openapi v0.0.0-20200805222855-6aeccd4b50c6
	k8s.io/kubectl v0.19.2
)

replace k8s.io/client-go => k8s.io/client-go v0.19.2
