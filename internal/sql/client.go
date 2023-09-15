package sql

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"gorm.io/gorm"
)

const (
	maxOpenConns    = 5
	connMaxLifetime = time.Second * 30
)

//go:generate counterfeiter . Client

type Client interface {
	Connect() error
	CreateKubernetesProvider(kubernetes.Provider) error
	CreateKubernetesResource(kubernetes.Resource) error
	DeleteKubernetesProvider(string) error
	DeleteKubernetesResourcesByAccountName(string) error
	GetKubernetesProvider(string) (kubernetes.Provider, error)
	GetKubernetesProviderAndPermissions(string) (kubernetes.Provider, error)
	ListKubernetesAccountsBySpinnakerApp(string) ([]string, error)
	ListKubernetesClustersByApplication(string) ([]kubernetes.Resource, error)
	ListKubernetesClustersByFields(...string) ([]kubernetes.Resource, error)
	ListKubernetesProviders() ([]kubernetes.Provider, error)
	ListKubernetesProvidersAndPermissions() ([]kubernetes.Provider, error)
	ListKubernetesResourcesByFields(...string) ([]kubernetes.Resource, error)
	ListKubernetesResourcesByTaskID(string) ([]kubernetes.Resource, error)
	ListReadGroupsByAccountName(string) ([]string, error)
	ListWriteGroupsByAccountName(string) ([]string, error)
	WithConfig(*gorm.Config)
}

func NewClient(dialector gorm.Dialector) Client {
	return &client{
		config:    &gorm.Config{},
		dialector: dialector,
	}
}

type client struct {
	config    *gorm.Config
	dialector gorm.Dialector
	db        *gorm.DB
}

// Connect sets up the database connection and creates tables.
func (c *client) Connect() error {
	db, err := gorm.Open(c.dialector, c.config)
	if err != nil {
		return fmt.Errorf("error opening connection to DB: %w", err)
	}

	err = db.AutoMigrate(
		&kubernetes.Provider{},
		&kubernetes.Resource{},
		&kubernetes.ProviderNamespaces{},
		&clouddriver.ReadPermission{},
		&clouddriver.WritePermission{},
	)
	if err != nil {
		return fmt.Errorf("error migrating DB: %w", err)
	}

	d, err := db.DB()
	if err != nil {
		return fmt.Errorf("error getting sql.DB: %w", err)
	}

	d.SetMaxOpenConns(maxOpenConns)
	d.SetMaxIdleConns(1)
	d.SetConnMaxLifetime(connMaxLifetime)

	c.db = db

	return nil
}

// CreateKubernetesProvider inserts the provider and permissions into the DB.
func (c *client) CreateKubernetesProvider(p kubernetes.Provider) error {
	err := c.db.Create(&p).Error
	if err != nil {
		return err
	}

	for _, group := range p.Permissions.Read {
		rp := clouddriver.ReadPermission{
			ID:          uuid.New().String(),
			AccountName: p.Name,
			ReadGroup:   group,
		}

		err = c.db.Create(&rp).Error
		if err != nil {
			return err
		}
	}

	for _, group := range p.Permissions.Write {
		wp := clouddriver.WritePermission{
			ID:          uuid.New().String(),
			AccountName: p.Name,
			WriteGroup:  group,
		}

		err = c.db.Create(&wp).Error
		if err != nil {
			return err
		}
	}

	for _, namespace := range p.Namespaces {
		ns := kubernetes.ProviderNamespaces{
			AccountName: p.Name,
			Namespace:   namespace,
		}

		err = c.db.Create(&ns).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateKubernetesResource inserts the resource into the DB.
func (c *client) CreateKubernetesResource(r kubernetes.Resource) error {
	db := c.db.Create(&r)
	return db.Error
}

// DeleteKubernetesProvider deletes the provider, namespaces, and permission from the DB.
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

	err = c.db.Where("account_name = ?", name).Delete(&kubernetes.ProviderNamespaces{}).Error
	if err != nil {
		return err
	}

	err = c.DeleteKubernetesResourcesByAccountName(name)
	if err != nil {
		return err
	}

	return nil
}

// DeleteKubernetesResources deletes all resources for the given provider from the DB.
func (c *client) DeleteKubernetesResourcesByAccountName(account string) error {
	err := c.db.Where("account_name = ?", account).Delete(&kubernetes.Resource{}).Error
	if err != nil {
		return err
	}

	return nil
}

// GetKubernetesProvider reads the provider from the DB.
func (c *client) GetKubernetesProvider(name string) (kubernetes.Provider, error) {
	p := kubernetes.Provider{}
	rows, err := c.db.Table("kubernetes_providers a").
		Select("a.name, a.host, a.ca_data, a.bearer_token, a.token_provider, a.namespace as legacy_namespace, b.namespace").
		Joins("LEFT JOIN "+kubernetes.ProviderNamespaces{}.TableName()+" b ON a.name = b.account_name").
		Where("a.name = ?", name).
		Rows()

	if err != nil {
		return p, err
	}
	defer rows.Close()

	namespaces := []string{}

	for rows.Next() {
		var r struct {
			CAData          string
			Host            string
			Name            string
			BearerToken     string
			LegacyNamespace *string
			Namespace       *string
			TokenProvider   string
		}

		err = rows.Scan(&r.Name, &r.Host, &r.CAData, &r.BearerToken, &r.TokenProvider, &r.LegacyNamespace, &r.Namespace)
		if err != nil {
			return p, err
		}

		p = kubernetes.Provider{
			Name:          r.Name,
			Host:          r.Host,
			CAData:        r.CAData,
			BearerToken:   r.BearerToken,
			TokenProvider: r.TokenProvider,
		}

		if r.LegacyNamespace != nil {
			namespaces = append(namespaces, *r.LegacyNamespace)
		}

		if r.Namespace != nil {
			namespaces = append(namespaces, *r.Namespace)
		}
	}

	p.Namespaces = namespaces

	if p.Name == "" {
		return p, gorm.ErrRecordNotFound
	}

	return p, nil
}

// GetKubernetesProviderAndPermissions reads the provider and permissions from the DB.
//
//				select a.name, a.host, a.ca_data, a.token_provider, a.namespace as legacy_namespace, d.namespace, b.read_group, c.write_group
//	         from kubernetes_providers a
//		  		left join provider_read_permissions b on a.name = b.account_name
//		  		left join provider_write_permissions c on a.name = c.account_name
//				left join kubernetes_providers_namespaces d on a.name = d.account_name
//		  	 	where a.name = ?;
func (c *client) GetKubernetesProviderAndPermissions(name string) (kubernetes.Provider, error) {
	p := kubernetes.Provider{}

	rows, err := c.db.Table("kubernetes_providers a").
		Select("a.name, "+
			"a.host, "+
			"a.ca_data, "+
			"a.token_provider, "+
			"a.namespace as legacy_namespace, "+
			"d.namespace, "+
			"b.read_group, "+
			"c.write_group").
		Joins("LEFT JOIN provider_read_permissions b ON a.name = b.account_name").
		Joins("LEFT JOIN provider_write_permissions c ON a.name = c.account_name").
		Joins("LEFT JOIN "+kubernetes.ProviderNamespaces{}.TableName()+" d ON a.name = d.account_name").
		Where("a.name = ?", name).
		Rows()
	if err != nil {
		return p, err
	}
	defer rows.Close()

	readGroups := map[string][]string{}
	writeGroups := map[string][]string{}
	namespaces := []string{}

	for rows.Next() {
		var r struct {
			CAData          string
			Host            string
			Name            string
			LegacyNamespace *string
			Namespace       *string
			ReadGroup       *string
			WriteGroup      *string
			TokenProvider   string
		}

		err = rows.Scan(&r.Name, &r.Host, &r.CAData, &r.TokenProvider, &r.LegacyNamespace, &r.Namespace, &r.ReadGroup, &r.WriteGroup)
		if err != nil {
			return p, err
		}

		p = kubernetes.Provider{
			Name:          r.Name,
			Host:          r.Host,
			CAData:        r.CAData,
			TokenProvider: r.TokenProvider,
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

		if r.LegacyNamespace != nil {
			namespaces = append(namespaces, *r.LegacyNamespace)
		}

		if r.Namespace != nil {
			namespaces = append(namespaces, *r.Namespace)
		}
	}

	p.Permissions.Read = readGroups[name]
	p.Permissions.Write = writeGroups[name]
	p.Namespaces = namespaces

	if p.Name == "" {
		return p, gorm.ErrRecordNotFound
	}

	return p, nil
}

// ListKubernetesClustersByApplication gets the list of kubernetes clusters
// for a Spinnaker application from the DB.
//
// A Kubernetes cluster is of kind deployment, statefulSet, replicaSet,
// ingress, service, and daemonSet.
func (c *client) ListKubernetesClustersByApplication(spinnakerApp string) ([]kubernetes.Resource, error) {
	var rs []kubernetes.Resource
	db := c.db.Select("account_name, cluster").
		Where("spinnaker_app = ? AND UPPER(kind) in ('DEPLOYMENT', 'STATEFULSET', 'REPLICASET', 'INGRESS', 'SERVICE', 'DAEMONSET')",
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
	db := c.db.Select(list).Where("UPPER(kind) in ('DEPLOYMENT', 'STATEFULSET', 'REPLICASET', 'INGRESS', 'SERVICE', 'DAEMONSET')").Group(list).Find(&rs)

	return rs, db.Error
}

// ListKubernetesProviders gets all the kubernetes providers from the DB.
func (c *client) ListKubernetesProviders() ([]kubernetes.Provider, error) {
	var ps []kubernetes.Provider

	rows, err := c.db.Table("kubernetes_providers a").
		Select("a.name, " +
			"a.host, " +
			"a.ca_data, " +
			"a.token_provider, " +
			"a.namespace as legacy_namespace, " +
			"b.namespace").
		Joins("LEFT JOIN kubernetes_providers_namespaces b ON a.name = b.account_name").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	providers := map[string]kubernetes.Provider{}

	for rows.Next() {
		var r struct {
			CAData          string
			Host            string
			Name            string
			LegacyNamespace *string
			Namespace       *string
			TokenProvider   string
		}

		err = rows.Scan(&r.Name, &r.Host, &r.CAData, &r.TokenProvider, &r.LegacyNamespace, &r.Namespace)
		if err != nil {
			return nil, err
		}

		if _, ok := providers[r.Name]; !ok {
			p := kubernetes.Provider{
				Name:          r.Name,
				Host:          r.Host,
				CAData:        r.CAData,
				TokenProvider: r.TokenProvider,
			}
			providers[r.Name] = p
		}

		if r.LegacyNamespace != nil {
			p := providers[r.Name]
			p.Namespaces = append(p.Namespaces, *r.LegacyNamespace)
			providers[r.Name] = p
		}

		if r.Namespace != nil {
			p := providers[r.Name]
			p.Namespaces = append(p.Namespaces, *r.Namespace)
			providers[r.Name] = p
		}
	}

	for _, provider := range providers {
		ps = append(ps, provider)
	}

	// Sort ascending by name.
	sort.Slice(ps, func(i, j int) bool {
		return ps[i].Name < ps[j].Name
	})

	return ps, nil
}

// ListKubernetesProvidersAndPermissions gets all the kubernetes providers,
// with their read/write permissions, from the DB.
//
//			select a.name, a.host, a.ca_data, a.token_provider, a.namespace, b.read_group, c.write_group from kubernetes_providers a
//	  		left join provider_read_permissions b on a.name = b.account_name
//	  		left join provider_write_permissions c on a.name = c.account_name;
func (c *client) ListKubernetesProvidersAndPermissions() ([]kubernetes.Provider, error) {
	ps := []kubernetes.Provider{}

	rows, err := c.db.Table("kubernetes_providers a").
		Select("a.name, " +
			"a.host, " +
			"a.ca_data, " +
			"a.token_provider, " +
			"a.namespace as legacy_namespace, " +
			"d.namespace, " +
			"b.read_group, " +
			"c.write_group").
		Joins("LEFT JOIN provider_read_permissions b ON a.name = b.account_name").
		Joins("LEFT JOIN provider_write_permissions c ON a.name = c.account_name").
		Joins("LEFT JOIN kubernetes_providers_namespaces d ON a.name = d.account_name").
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
			CAData          string
			Host            string
			Name            string
			LegacyNamespace *string
			Namespace       *string
			ReadGroup       *string
			WriteGroup      *string
			TokenProvider   string
		}

		err = rows.Scan(&r.Name, &r.Host, &r.CAData, &r.TokenProvider, &r.LegacyNamespace, &r.Namespace, &r.ReadGroup, &r.WriteGroup)
		if err != nil {
			return nil, err
		}

		if _, ok := providers[r.Name]; !ok {
			p := kubernetes.Provider{
				Name:          r.Name,
				Host:          r.Host,
				CAData:        r.CAData,
				TokenProvider: r.TokenProvider,
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

		if r.LegacyNamespace != nil {
			p := providers[r.Name]
			p.Namespaces = append(p.Namespaces, *r.LegacyNamespace)
			providers[r.Name] = p
		}

		if r.Namespace != nil {
			p := providers[r.Name]
			p.Namespaces = append(p.Namespaces, *r.Namespace)
			providers[r.Name] = p
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

// ListKubernetesResourcesByTaskID get the list of kubernetes resources
// by task ID from the DB.
func (c *client) ListKubernetesResourcesByTaskID(taskID string) ([]kubernetes.Resource, error) {
	var rs []kubernetes.Resource
	db := c.db.Select("account_name, api_group, kind, name, artifact_name, namespace, resource, task_type, version").
		Where("task_id = ?", taskID).Find(&rs)

	return rs, db.Error
}

// ListKubernetesResourcesByFields gets the list of unique set of
// kubernetes resource attributes from the DB.
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

// ListKubernetesAccountsBySpinnakerApp gets the list of account names
// for a Spinnaker application from the DB.
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

// ListReadGroupsByAccountName gets the list of groups with read permission
// for an account name from the DB.
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

// ListWriteGroupsByAccountName gets the list of groups with write permission
// for an account name from the DB.
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

// WithConfig sets the gorm config to use.
func (c *client) WithConfig(config *gorm.Config) {
	c.config = config
}
