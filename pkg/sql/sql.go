package sql

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/jinzhu/gorm"

	// Needed for connection.
	_ "github.com/go-sql-driver/mysql"

	// Needed for connection.
	_ "github.com/mattn/go-sqlite3"
)

const (
	ClientInstanceKey = `SQLClient`
	maxOpenConns      = 5
	connMaxLifetime   = time.Second * 30
)

//go:generate counterfeiter . Client

type Client interface {
	CreateKubernetesProvider(kubernetes.Provider) error
	CreateKubernetesResource(kubernetes.Resource) error
	CreateReadPermission(clouddriver.ReadPermission) error
	CreateWritePermission(clouddriver.WritePermission) error
	DeleteKubernetesProvider(string) error
	GetKubernetesProvider(string) (kubernetes.Provider, error)
	ListKubernetesAccountsBySpinnakerApp(string) ([]string, error)
	ListKubernetesClustersByApplication(string) ([]kubernetes.Resource, error)
	ListKubernetesClustersByFields(...string) ([]kubernetes.Resource, error)
	ListKubernetesProviders() ([]kubernetes.Provider, error)
	ListKubernetesProvidersAndPermissions() ([]kubernetes.Provider, error)
	ListKubernetesResourcesByFields(...string) ([]kubernetes.Resource, error)
	ListKubernetesResourcesByTaskID(string) ([]kubernetes.Resource, error)
	ListKubernetesResourceNamesByAccountNameAndKindAndNamespace(string, string, string) ([]string, error)
	ListReadGroupsByAccountName(string) ([]string, error)
	ListWriteGroupsByAccountName(string) ([]string, error)
}

func NewClient(db *gorm.DB) Client {
	return &client{db: db}
}

type client struct {
	db *gorm.DB
}

// Connect sets up the database connection and creates tables.
//
// Connection is of type interface{} - this allows for tests to
// pass in a sqlmock connection and for main to connect given a
// connection string.
func Connect(driver string, connection interface{}) (*gorm.DB, error) {
	db, err := gorm.Open(driver, connection)
	if err != nil {
		return nil, err
	}

	db.LogMode(false)
	db.AutoMigrate(
		&kubernetes.Cluster{},
		&kubernetes.Provider{},
		&kubernetes.Resource{},
		&clouddriver.ReadPermission{},
		&clouddriver.WritePermission{},
	)

	db.DB().SetMaxOpenConns(maxOpenConns)
	db.DB().SetMaxIdleConns(1)
	db.DB().SetConnMaxLifetime(connMaxLifetime)

	return db, nil
}

func Instance(c *gin.Context) Client {
	return c.MustGet(ClientInstanceKey).(Client)
}

type Config struct {
	User     string
	Password string
	Host     string
	Name     string
}

// Get driver and connection string to the DB.
func Connection(c Config) (string, string) {
	if c.User == "" || c.Password == "" || c.Host == "" || c.Name == "" {
		log.Println("[CLOUDDRIVER] SQL config missing field - defaulting to local sqlite DB.")
		return "sqlite3", "clouddriver.db"
	}

	return "mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=UTC",
		c.User, c.Password, c.Host, c.Name)
}

func (c *client) CreateKubernetesCluster(p kubernetes.Cluster) error {
	db := c.db.Create(&p)
	return db.Error
}

func (c *client) CreateKubernetesProvider(p kubernetes.Provider) error {
	db := c.db.Create(&p)
	return db.Error
}

func (c *client) CreateKubernetesResource(r kubernetes.Resource) error {
	db := c.db.Create(&r)
	return db.Error
}

func (c *client) CreateWritePermission(w clouddriver.WritePermission) error {
	db := c.db.Create(&w)
	return db.Error
}

func (c *client) CreateReadPermission(r clouddriver.ReadPermission) error {
	db := c.db.Create(&r)
	return db.Error
}

func (c *client) DeleteKubernetesProvider(name string) error {
	err := c.db.Delete(&kubernetes.Provider{Name: name}).Error
	if err != nil {
		return err
	}

	err = c.db.Where("account_name = ?", name).Delete(&clouddriver.ReadPermission{}).Error
	if err != nil {
		return err
	}

	err = c.db.Where("account_name = ?", name).Delete(&clouddriver.WritePermission{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *client) GetKubernetesProvider(name string) (kubernetes.Provider, error) {
	var p kubernetes.Provider
	db := c.db.Select("host, ca_data, bearer_token").Where("name = ?", name).First(&p)

	return p, db.Error
}

// A Kubernetes cluster is of kind deployment, statefulSet, replicaSet, ingress, service, and daemonSet.
func (c *client) ListKubernetesClustersByApplication(spinnakerApp string) ([]kubernetes.Resource, error) {
	var rs []kubernetes.Resource
	db := c.db.Select("account_name, cluster").
		Where("spinnaker_app = ? AND kind in ('deployment', 'statefulSet', 'replicaSet', 'ingress', 'service', 'daemonSet')",
			spinnakerApp).
		Group("account_name, cluster").Find(&rs)

	return rs, db.Error
}

func (c *client) ListKubernetesClustersByFields(fields ...string) ([]kubernetes.Resource, error) {
	if len(fields) == 0 {
		return nil, errors.New("no fields provided")
	}

	list := ""
	for i, field := range fields {
		list += field
		if i != len(fields)-1 {
			list += ", "
		}
	}

	var rs []kubernetes.Resource
	db := c.db.Select(list).Where("kind in ('deployment', 'statefulSet', 'replicaSet', 'ingress', 'service', 'daemonSet')").Group(list).Find(&rs)

	return rs, db.Error
}

func (c *client) ListKubernetesResourceNamesByAccountNameAndKindAndNamespace(accountName,
	kind, namespace string) ([]string, error) {
	rs := []kubernetes.Resource{}
	names := []string{}
	db := c.db.Select("name").
		Where("account_name = ? AND kind = ? AND namespace = ?",
			accountName, kind, namespace).
		Group("name").Find(&rs)

	for _, r := range rs {
		if r.Name != "" {
			names = append(names, r.Name)
		}
	}

	return names, db.Error
}

func (c *client) ListKubernetesProviders() ([]kubernetes.Provider, error) {
	var ps []kubernetes.Provider
	db := c.db.Select("name, host, ca_data").Find(&ps)

	return ps, db.Error
}

// select a.name, a.host, a.ca_data, b.read_group, c.write_group from kubernetes_providers a
//   left join provider_read_permissions b on a.name = b.account_name
//   left join provider_write_permissions c on a.name = c.account_name;
func (c *client) ListKubernetesProvidersAndPermissions() ([]kubernetes.Provider, error) {
	ps := []kubernetes.Provider{}

	rows, err := c.db.Table("kubernetes_providers a").
		Select("a.name, " +
			"a.host, " +
			"a.ca_data, " +
			"b.read_group, " +
			"c.write_group").
		Joins("LEFT JOIN provider_read_permissions b ON a.name = b.account_name").
		Joins("LEFT JOIN provider_write_permissions c ON a.name = c.account_name").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	providers := map[string]kubernetes.Provider{}
	readGroups := map[string][]string{}
	writeGroups := map[string][]string{}

	for rows.Next() {
		var r struct {
			CAData     string
			Host       string
			Name       string
			ReadGroup  *string
			WriteGroup *string
		}

		err = rows.Scan(&r.Name, &r.Host, &r.CAData, &r.ReadGroup, &r.WriteGroup)
		if err != nil {
			return nil, err
		}

		if _, ok := providers[r.Name]; !ok {
			p := kubernetes.Provider{
				Name:   r.Name,
				Host:   r.Host,
				CAData: r.CAData,
			}
			providers[r.Name] = p
		}

		if r.ReadGroup != nil {
			if _, ok := readGroups[r.Name]; !ok {
				readGroups[r.Name] = []string{}
			}

			if !contains(readGroups[r.Name], *r.ReadGroup) {
				readGroups[r.Name] = append(readGroups[r.Name], *r.ReadGroup)
			}
		}

		if r.WriteGroup != nil {
			if _, ok := writeGroups[r.Name]; !ok {
				writeGroups[r.Name] = []string{}
			}

			if !contains(writeGroups[r.Name], *r.WriteGroup) {
				writeGroups[r.Name] = append(writeGroups[r.Name], *r.WriteGroup)
			}
		}
	}

	for name, provider := range providers {
		provider.Permissions.Read = readGroups[name]
		provider.Permissions.Write = writeGroups[name]
		ps = append(ps, provider)
	}

	// Sort ascending by name.
	sort.Slice(ps, func(i, j int) bool {
		return ps[i].Name < ps[j].Name
	})

	return ps, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func (c *client) ListKubernetesResourcesByTaskID(taskID string) ([]kubernetes.Resource, error) {
	var rs []kubernetes.Resource
	db := c.db.Select("account_name, api_group, kind, name, namespace, resource, task_type, version").
		Where("task_id = ?", taskID).Find(&rs)

	return rs, db.Error
}

func (c *client) ListKubernetesResourcesByFields(fields ...string) ([]kubernetes.Resource, error) {
	if len(fields) == 0 {
		return nil, errors.New("no fields provided")
	}

	list := ""
	for i, field := range fields {
		list += field
		if i != len(fields)-1 {
			list += ", "
		}
	}

	var rs []kubernetes.Resource
	db := c.db.Select(list).Group(list).Find(&rs)

	return rs, db.Error
}

func (c *client) ListKubernetesAccountsBySpinnakerApp(spinnakerApp string) ([]string, error) {
	var rs []kubernetes.Resource
	db := c.db.Select("account_name").
		Where("spinnaker_app = ?", spinnakerApp).
		Group("account_name").
		Find(&rs)

	accounts := []string{}
	for _, r := range rs {
		accounts = append(accounts, r.AccountName)
	}

	return accounts, db.Error
}

func (c *client) ListReadGroupsByAccountName(accountName string) ([]string, error) {
	r := []clouddriver.ReadPermission{}
	db := c.db.Select("read_group").
		Where("account_name = ?", accountName).
		Group("read_group").
		Find(&r)

	groups := []string{}
	for _, v := range r {
		groups = append(groups, v.ReadGroup)
	}

	return groups, db.Error
}

func (c *client) ListWriteGroupsByAccountName(accountName string) ([]string, error) {
	w := []clouddriver.WritePermission{}
	db := c.db.Select("write_group").
		Where("account_name = ?", accountName).
		Group("write_group").
		Find(&w)

	groups := []string{}
	for _, v := range w {
		groups = append(groups, v.WriteGroup)
	}

	return groups, db.Error
}
