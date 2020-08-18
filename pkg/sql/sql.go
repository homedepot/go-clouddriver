package sql

import (
	"fmt"
	"time"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/gin-gonic/gin"
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
	GetKubernetesProvider(string) (kubernetes.Provider, error)
	CreateKubernetesResource(kubernetes.Resource) error
	ListKubernetesResources(string) ([]kubernetes.Resource, error)
	ListKubernetesAccountsBySpinnakerApp(string) ([]string, error)
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
		&kubernetes.Provider{},
		&kubernetes.Resource{},
	)

	db.DB().SetMaxOpenConns(maxOpenConns)
	db.DB().SetMaxIdleConns(1)
	db.DB().SetConnMaxLifetime(connMaxLifetime)

	return db, nil
}

func (c *client) CreateKubernetesProvider(p kubernetes.Provider) error {
	db := c.db.Create(&p)
	return db.Error
}

func (c *client) GetKubernetesProvider(name string) (kubernetes.Provider, error) {
	var p kubernetes.Provider
	db := c.db.Select("host, ca_data").Where("name = ?", name).First(&p)

	return p, db.Error
}

func (c *client) CreateKubernetesResource(r kubernetes.Resource) error {
	db := c.db.Create(&r)
	return db.Error
}

func (c *client) ListKubernetesResources(taskID string) ([]kubernetes.Resource, error) {
	var rs []kubernetes.Resource
	db := c.db.Select("account_name, api_group, kind, name, namespace, resource, version").
		Where("task_id = ?", taskID).Find(&rs)

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

func Instance(c *gin.Context) Client {
	return c.MustGet(ClientInstanceKey).(Client)
}

type Config struct {
	User     string
	Password string
	Host     string
	Name     string
}

// Get connection to the DB.
func Connection(c Config) string {
	if c.User == "" || c.Password == "" || c.Host == "" || c.Name == "" {
		return "clouddriver.db"
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=UTC",
		c.User, c.Password, c.Host, c.Name)
}
