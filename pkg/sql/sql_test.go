package sql_test

import (
	"database/sql"
	"io/ioutil"
	"log"

	clouddriver "github.com/billiford/go-clouddriver/pkg"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	. "github.com/billiford/go-clouddriver/pkg/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sql", func() {
	var (
		db   *gorm.DB
		mock sqlmock.Sqlmock
		d    *sql.DB
		c    Client
		err  error
	)

	BeforeEach(func() {
		// Mock DB.
		d, mock, _ = sqlmock.New()
		db, err = Connect("sqlite3", d)
		// Enable DB logging.
		// db.LogMode(true)
		c = NewClient(db)

		log.SetOutput(ioutil.Discard)
	})

	AfterEach(func() {
		db.Close()
	})

	Describe("#Connect", func() {
		When("it fails to connect", func() {
			BeforeEach(func() {
				_, err = Connect("mysql", "mysql")
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("invalid DSN: missing the slash separating the database name"))
			})
		})
	})

	Describe("#CreateKubernetesProvider", func() {
		var provider kubernetes.Provider

		BeforeEach(func() {
			provider = kubernetes.Provider{
				Name:   "test-name",
				Host:   "test-host",
				CAData: "test-ca-data",
			}
		})

		JustBeforeEach(func() {
			err = c.CreateKubernetesProvider(provider)
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				mock.ExpectBegin()
				mock.ExpectExec(`(?i)^INSERT INTO "kubernetes_providers" \(` +
					`"name",` +
					`"host",` +
					`"ca_data",` +
					`"bearer_token"` +
					`\) VALUES \(\?,\?,\?,\?\)$`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("#CreateKubernetesResource", func() {
		var resource kubernetes.Resource

		BeforeEach(func() {
			resource = kubernetes.Resource{
				ID:        "test-id",
				TaskID:    "test-task-id",
				APIGroup:  "test-group",
				Name:      "test-name",
				Namespace: "test-namespace",
				Resource:  "test-resource",
				Version:   "test-version",
			}
		})

		JustBeforeEach(func() {
			err = c.CreateKubernetesResource(resource)
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				mock.ExpectBegin()
				mock.ExpectExec(`(?i)^INSERT INTO "kubernetes_resources" \(` +
					`"account_name",` +
					`"id",` +
					`"task_id",` +
					`"api_group",` +
					`"name",` +
					`"namespace",` +
					`"resource",` +
					`"version",` +
					`"kind",` +
					`"spinnaker_app",` +
					`"cluster"` +
					`\) VALUES \(\?,\?,\?,\?,\?,\?,\?,\?,\?,\?,\?\)$`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("#CreateReadPermission", func() {
		var rp clouddriver.ReadPermission

		BeforeEach(func() {
			rp = clouddriver.ReadPermission{
				ID:          "test-id",
				AccountName: "test-account-name",
				ReadGroup:   "test-write-group",
			}
		})

		JustBeforeEach(func() {
			err = c.CreateReadPermission(rp)
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				mock.ExpectBegin()
				mock.ExpectExec(`(?i)^INSERT INTO "provider_read_permissions" \(` +
					`"id",` +
					`"account_name",` +
					`"read_group"` +
					`\) VALUES \(\?,\?,\?\)$`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("#CreateWritePermission", func() {
		var wp clouddriver.WritePermission

		BeforeEach(func() {
			wp = clouddriver.WritePermission{
				ID:          "test-id",
				AccountName: "test-account-name",
				WriteGroup:  "test-write-group",
			}
		})

		JustBeforeEach(func() {
			err = c.CreateWritePermission(wp)
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				mock.ExpectBegin()
				mock.ExpectExec(`(?i)^INSERT INTO "provider_write_permissions" \(` +
					`"id",` +
					`"account_name",` +
					`"write_group"` +
					`\) VALUES \(\?,\?,\?\)$`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("#DeleteKubernetesProvider", func() {
		var name string

		BeforeEach(func() {
			name = "test-name"
		})

		JustBeforeEach(func() {
			err = c.DeleteKubernetesProvider(name)
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				mock.ExpectBegin()
				mock.ExpectExec(`(?i)^DELETE FROM "kubernetes_providers" WHERE
				"kubernetes_providers"."name" = \?$`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				mock.ExpectBegin()
				mock.ExpectExec(`(?i)^DELETE FROM "provider_read_permissions" WHERE
				\(account_name = \?\)$`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				mock.ExpectBegin()
				mock.ExpectExec(`(?i)^DELETE FROM "provider_write_permissions" WHERE
				\(account_name = \?\)$`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("#GetKubernetesProvider", func() {
		var provider kubernetes.Provider

		BeforeEach(func() {
			provider = kubernetes.Provider{}
		})

		JustBeforeEach(func() {
			provider, err = c.GetKubernetesProvider("test-name")
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				sqlRows := sqlmock.NewRows([]string{"name", "host", "ca_data"}).
					AddRow("test-name", "test-host", "test-ca-data")
				mock.ExpectQuery(`(?i)^SELECT host, ca_data, bearer_token FROM "kubernetes_providers" ` +
					` WHERE \(name = \?\) ORDER BY "kubernetes_providers"."name" ASC LIMIT 1$`).
					WillReturnRows(sqlRows)
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(provider.Name).To(Equal("test-name"))
				Expect(provider.Host).To(Equal("test-host"))
				Expect(provider.CAData).To(Equal("test-ca-data"))
			})
		})
	})

	Describe("#ListKubernetesClustersByApplication", func() {
		var resources []kubernetes.Resource

		JustBeforeEach(func() {
			resources, err = c.ListKubernetesClustersByApplication("test-application")
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				sqlRows := sqlmock.NewRows([]string{"account_name", "cluster"}).
					AddRow("account1", "cluster 1").
					AddRow("account2", "cluster 2")
				mock.ExpectQuery(`(?i)^SELECT ` +
					`account_name, ` +
					`cluster ` +
					`FROM "kubernetes_resources" ` +
					` WHERE \(spinnaker_app = \? AND kind in \('deployment',
						'statefulSet',
						'replicaSet',
						'ingress',
						'service',
						'daemonSet'\)\) GROUP BY
						account_name, cluster$`).
					WillReturnRows(sqlRows)
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(resources).To(HaveLen(2))
			})
		})
	})

	Describe("#ListKubernetesResourceNamesByAccountNameAndKindAndNamespace", func() {
		var names []string

		JustBeforeEach(func() {
			names, err = c.ListKubernetesResourceNamesByAccountNameAndKindAndNamespace("test-account-name",
				"test-kind", "test-namespace")
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				sqlRows := sqlmock.NewRows([]string{"name"}).
					AddRow("name1").
					AddRow("name2")
				mock.ExpectQuery(`(?i)^SELECT ` +
					`name ` +
					`FROM "kubernetes_resources" ` +
					` WHERE \(account_name = \? AND ` +
					`kind = \? AND ` +
					`namespace = \?\) GROUP BY name$`).
					WillReturnRows(sqlRows)
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(names).To(HaveLen(2))
			})
		})
	})

	Describe("#ListKubernetesResourcesByTaskID", func() {
		var resources []kubernetes.Resource

		JustBeforeEach(func() {
			resources, err = c.ListKubernetesResourcesByTaskID("test-task-id")
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				sqlRows := sqlmock.NewRows([]string{"group", "name"}).
					AddRow("group1", "name1").
					AddRow("group2", "name2")
				mock.ExpectQuery(`(?i)^SELECT ` +
					`account_name, ` +
					`api_group, ` +
					`kind, ` +
					`name, ` +
					`namespace, ` +
					`resource, ` +
					`version ` +
					`FROM "kubernetes_resources" ` +
					` WHERE \(task_id = \?\)$`).
					WillReturnRows(sqlRows)
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(resources).To(HaveLen(2))
			})
		})
	})

	Describe("#ListKubernetesProviders", func() {
		var providers []kubernetes.Provider

		JustBeforeEach(func() {
			providers, err = c.ListKubernetesProviders()
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				sqlRows := sqlmock.NewRows([]string{"name", "host", "ca_data"}).
					AddRow("name1", "host1", "ca_data1").
					AddRow("name2", "host2", "ca_data2")
				mock.ExpectQuery(`(?i)^SELECT ` +
					`name, ` +
					`host, ` +
					`ca_data ` +
					`FROM "kubernetes_providers"$`).
					WillReturnRows(sqlRows)
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(providers).To(HaveLen(2))
			})
		})
	})

	Describe("#ListKubernetesProvidersAndPermissions", func() {
		var providers []kubernetes.Provider

		JustBeforeEach(func() {
			providers, err = c.ListKubernetesProvidersAndPermissions()
		})

		When("getting the rows returns an error", func() {
			BeforeEach(func() {
				db.Close()
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("sql: database is closed"))
			})
		})

		When("scanning a row returns an error", func() {
			BeforeEach(func() {
				sqlRows := sqlmock.NewRows([]string{"name", "host", "read_group", "write_group"}).
					AddRow("name1", "host1", "read_group1", "write_group1")
				mock.ExpectQuery(`(?i)^SELECT ` +
					`a.name, ` +
					`a.host, ` +
					`a.ca_data, ` +
					`b.read_group, ` +
					`c.write_group ` +
					`FROM kubernetes_providers a ` +
					`left join provider_read_permissions b on a.name = b.account_name ` +
					`left join provider_write_permissions c on a.name = c.account_name$`).
					WillReturnRows(sqlRows)
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("sql: expected 4 destination arguments in Scan, not 5"))
			})
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				sqlRows := sqlmock.NewRows([]string{"name", "host", "ca_data", "read_group", "write_group"}).
					AddRow("name1", "host1", "ca_data1", "read_group1", "write_group1").
					AddRow("name1", "host1", "ca_data1", "read_group2", "write_group1").
					AddRow("name2", "host2", "ca_data2", "read_group2", "write_group2").
					AddRow("name2", "host2", "ca_data2", "read_group2", "write_group3")
				mock.ExpectQuery(`(?i)^SELECT ` +
					`a.name, ` +
					`a.host, ` +
					`a.ca_data, ` +
					`b.read_group, ` +
					`c.write_group ` +
					`FROM kubernetes_providers a ` +
					`left join provider_read_permissions b on a.name = b.account_name ` +
					`left join provider_write_permissions c on a.name = c.account_name$`).
					WillReturnRows(sqlRows)
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(providers).To(HaveLen(2))
				Expect(providers[0].Permissions.Read).To(HaveLen(2))
				Expect(providers[0].Permissions.Write).To(HaveLen(1))
				Expect(providers[1].Permissions.Read).To(HaveLen(1))
				Expect(providers[1].Permissions.Write).To(HaveLen(2))
			})
		})
	})

	Describe("#ListKubernetesResourcesByFields", func() {
		var resources []kubernetes.Resource
		fields := []string{}

		BeforeEach(func() {
			fields = []string{"field1", "field2"}
		})

		JustBeforeEach(func() {
			resources, err = c.ListKubernetesResourcesByFields(fields...)
		})

		When("no fields are provided", func() {
			BeforeEach(func() {
				fields = []string{}
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("no fields provided"))
			})
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				sqlRows := sqlmock.NewRows([]string{"group", "name"}).
					AddRow("group1", "name1").
					AddRow("group2", "name2")
				mock.ExpectQuery(`(?i)^SELECT ` +
					`field1, ` +
					`field2 ` +
					`FROM "kubernetes_resources" ` +
					` GROUP BY field1, field2$`).
					WillReturnRows(sqlRows)
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(resources).To(HaveLen(2))
			})
		})
	})

	Describe("#ListKubernetesAccountsBySpinnakerApp", func() {
		var accounts []string

		JustBeforeEach(func() {
			accounts, err = c.ListKubernetesAccountsBySpinnakerApp("test-spinnaker-app")
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				sqlRows := sqlmock.NewRows([]string{"account_name"}).
					AddRow("account1").
					AddRow("account2")
				mock.ExpectQuery(`(?i)^SELECT ` +
					`account_name ` +
					`FROM "kubernetes_resources" ` +
					` WHERE \(spinnaker_app = \?\) ` +
					`GROUP BY account_name$`).
					WillReturnRows(sqlRows)
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(accounts).To(HaveLen(2))
			})
		})
	})

	Describe("#ListReadGroupsByAccountName", func() {
		var groups []string

		JustBeforeEach(func() {
			groups, err = c.ListReadGroupsByAccountName("test-account-name")
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				sqlRows := sqlmock.NewRows([]string{"read_group"}).
					AddRow("group1").
					AddRow("group2")
				mock.ExpectQuery(`(?i)^SELECT ` +
					`read_group ` +
					`FROM "provider_read_permissions" ` +
					` WHERE \(account_name = \?\) ` +
					`GROUP BY read_group`).
					WillReturnRows(sqlRows)
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(groups).To(HaveLen(2))
			})
		})
	})

	Describe("#ListWriteGroupsByAccountName", func() {
		var groups []string

		JustBeforeEach(func() {
			groups, err = c.ListWriteGroupsByAccountName("test-account-name")
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				sqlRows := sqlmock.NewRows([]string{"write_group"}).
					AddRow("group1").
					AddRow("group2")
				mock.ExpectQuery(`(?i)^SELECT ` +
					`write_group ` +
					`FROM "provider_write_permissions" ` +
					` WHERE \(account_name = \?\) ` +
					`GROUP BY write_group`).
					WillReturnRows(sqlRows)
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(groups).To(HaveLen(2))
			})
		})
	})

	Describe("#Instance", func() {
		var ctx *gin.Context
		var c2 Client

		BeforeEach(func() {
			ctx = &gin.Context{}
			ctx.Set(ClientInstanceKey, c)
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				c2 = Instance(ctx)
			})

			It("succeeds", func() {
				Expect(c2).ToNot(BeNil())
			})
		})
	})

	Describe("#Connection", func() {
		var driver, connection string
		var c Config

		When("the config is not set", func() {
			BeforeEach(func() {
				driver, connection = Connection(c)
			})

			It("returns a disk db", func() {
				Expect(driver).To(Equal("sqlite3"))
				Expect(connection).To(Equal("clouddriver.db"))
			})
		})

		When("the config is set", func() {
			BeforeEach(func() {
				c = Config{
					User:     "user",
					Password: "password",
					Host:     "10.1.1.1",
					Name:     "go-clouddriver",
				}
				driver, connection = Connection(c)
			})

			It("returns a mysql connection string", func() {
				Expect(driver).To(Equal("mysql"))
				Expect(connection).To(Equal("user:password@tcp(10.1.1.1)/go-clouddriver?charset=utf8&parseTime=True&loc=UTC"))
			})
		})
	})
})
