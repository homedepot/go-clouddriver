package sql_test

import (
	"database/sql"
	"fmt"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	. "github.com/homedepot/go-clouddriver/internal/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sql", func() {
	var (
		mock sqlmock.Sqlmock
		d    *sql.DB
		c    Client
		err  error
	)

	BeforeEach(func() {
		// Mock DB.
		d, mock, _ = sqlmock.New()
		// Define a new MySQL dialector that uses our mocked DB.
		dialector := mysql.New(mysql.Config{
			Conn:                      d,
			DefaultStringSize:         256,
			SkipInitializeWithVersion: true,
		})
		c = NewClient(dialector)

		// Create a new logger that disables logging.
		newLogger := logger.New(nil, logger.Config{})
		config := &gorm.Config{
			Logger: newLogger,
		}
		c.WithConfig(config)

		mock.ExpectExec("(?i)^CREATE TABLE `kubernetes_providers` " +
			"\\(`name`\\ varchar\\(256\\)," +
			"`host` varchar\\(256\\)," +
			"`ca_data` text," +
			"`bearer_token` varchar\\(2048\\)," +
			"`token_provider` varchar\\(32\\) NOT NULL DEFAULT 'google'," +
			"`namespace` varchar\\(253\\)," +
			"PRIMARY KEY \\(`name`\\)" +
			"\\)$").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("(?i)^CREATE TABLE `kubernetes_resources` " +
			"\\(`account_name`\\ varchar\\(256\\)," +
			"`id` varchar\\(256\\)," +
			"`timestamp` timestamp DEFAULT current_timestamp," +
			"`task_id` varchar\\(256\\)," +
			"`task_type` varchar\\(256\\)," +
			"`api_group` varchar\\(256\\)," +
			"`name` varchar\\(256\\)," +
			"`artifact_name` varchar\\(256\\)," +
			"`namespace` varchar\\(256\\)," +
			"`resource` varchar\\(256\\)," +
			"`version` varchar\\(256\\)," +
			"`kind` varchar\\(256\\)," +
			"`spinnaker_app` varchar\\(256\\)," +
			"`cluster` varchar\\(256\\)," +
			"PRIMARY KEY \\(`id`\\)" +
			"\\)$").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("CREATE TABLE `kubernetes_providers_namespaces` " +
			"\\(`account_name` varchar\\(256\\)," +
			"`namespace` varchar\\(256\\)").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("(?i)^CREATE TABLE `provider_read_permissions` " +
			"\\(`id`\\ varchar\\(256\\)," +
			"`account_name` varchar\\(256\\)," +
			"`read_group` varchar\\(256\\)," +
			"PRIMARY KEY \\(`id`\\)" +
			"\\)$").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("(?i)^CREATE TABLE `provider_write_permissions` " +
			"\\(`id`\\ varchar\\(256\\)," +
			"`account_name` varchar\\(256\\)," +
			"`write_group` varchar\\(256\\)," +
			"PRIMARY KEY \\(`id`\\)" +
			"\\)$").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = c.Connect()
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		d.Close()
	})

	Describe("#Connect", func() {
		When("it fails to connect", func() {
			BeforeEach(func() {
				dialector := mysql.New(mysql.Config{
					DSN: "mysql",
				})
				c = NewClient(dialector)
				newLogger := logger.New(nil, logger.Config{})
				config := &gorm.Config{
					Logger: newLogger,
				}
				c.WithConfig(config)
				err = c.Connect()
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("error opening connection to DB: " +
					"invalid DSN: missing the slash separating the database name"))
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
				Permissions: kubernetes.ProviderPermissions{
					Read:  []string{"test-read-group"},
					Write: []string{"test-write-group"},
				},
			}
		})

		JustBeforeEach(func() {
			err = c.CreateKubernetesProvider(provider)
		})

		When("tokenProvider is set", func() {
			BeforeEach(func() {
				provider = kubernetes.Provider{
					Name:          "test-name",
					Host:          "test-host",
					CAData:        "test-ca-data",
					TokenProvider: "test-token",
				}
				mock.ExpectBegin()
				mock.ExpectExec("(?i)^INSERT INTO `kubernetes_providers` \\(" +
					"`name`" +
					",`host`" +
					",`ca_data`" +
					",`bearer_token`" +
					",`token_provider`" +
					",`namespace`" +
					"\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?\\)$").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
			})
		})

		When("namespaces are set", func() {
			BeforeEach(func() {
				provider = kubernetes.Provider{
					Name:          "test-name",
					Host:          "test-host",
					CAData:        "test-ca-data",
					TokenProvider: "test-token",
					Namespaces:    []string{"n1", "n2", "n3"},
				}
				mock.ExpectBegin()
				mock.ExpectExec("(?i)^INSERT INTO `kubernetes_providers` \\(" +
					"`name`" +
					",`host`" +
					",`ca_data`" +
					",`bearer_token`" +
					",`token_provider`" +
					",`namespace`" +
					"\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?\\)$").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				for i := 1; i <= 3; i++ {
					mock.ExpectBegin()
					mock.ExpectExec("INSERT INTO `kubernetes_providers_namespaces` \\(`account_name`,`namespace`\\) VALUES \\(\\?,\\?\\)").
						WithArgs("test-name", "n"+fmt.Sprint(i)).
						WillReturnResult(sqlmock.NewResult(int64(i), 1))

					mock.ExpectCommit()

					// we make sure that all expectations were met
					if err := mock.ExpectationsWereMet(); err != nil {
						fmt.Errorf("there were unfulfilled expections: %s", err)
					}
				}

			})

			It("adds the namespaces to the DB and succeeds", func() {
				Expect(err).To(BeNil())
			})
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				mock.ExpectBegin()
				mock.ExpectExec("(?i)^INSERT INTO `kubernetes_providers` \\(" +
					"`name`" +
					",`host`" +
					",`ca_data`" +
					",`bearer_token`" +
					",`token_provider`" +
					",`namespace`" +
					"\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?\\)$").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				mock.ExpectBegin()
				mock.ExpectExec("(?i)^INSERT INTO `provider_read_permissions` \\(" +
					"`id`," +
					"`account_name`," +
					"`read_group`" +
					"\\) VALUES \\(\\?,\\?,\\?\\)$").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				mock.ExpectBegin()
				mock.ExpectExec("(?i)^INSERT INTO `provider_write_permissions` \\(" +
					"`id`," +
					"`account_name`," +
					"`write_group`" +
					"\\) VALUES \\(\\?,\\?,\\?\\)$").
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
				ID:           "test-id",
				TaskID:       "test-task-id",
				APIGroup:     "test-group",
				Name:         "test-name",
				ArtifactName: "test-name",
				Namespace:    "test-namespace",
				Resource:     "test-resource",
				Version:      "test-version",
			}
		})

		JustBeforeEach(func() {
			err = c.CreateKubernetesResource(resource)
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				mock.ExpectBegin()
				mock.ExpectExec("(?i)^INSERT INTO `kubernetes_resources` \\(" +
					"`account_name`," +
					"`id`," +
					"`task_id`," +
					"`task_type`," +
					"`api_group`," +
					"`name`," +
					"`artifact_name`," +
					"`namespace`," +
					"`resource`," +
					"`version`," +
					"`kind`," +
					"`spinnaker_app`," +
					"`cluster`" +
					"\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?\\)$").
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
				mock.ExpectExec("(?i)^DELETE FROM `kubernetes_providers` WHERE " +
					"`kubernetes_providers`.`name` = \\?$").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				mock.ExpectBegin()
				mock.ExpectExec("(?i)^DELETE FROM `provider_read_permissions` WHERE " +
					"account_name = \\?$").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				mock.ExpectBegin()
				mock.ExpectExec("(?i)^DELETE FROM `provider_write_permissions` WHERE " +
					"account_name = \\?$").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				mock.ExpectBegin()
				mock.ExpectExec("(?i)^DELETE FROM `" + kubernetes.ProviderNamespaces{}.TableName() + "` WHERE " +
					"account_name = \\?$").
					WillReturnResult(sqlmock.NewResult(1, 3))
				mock.ExpectCommit()

				mock.ExpectBegin()
				mock.ExpectExec("(?i)^DELETE FROM `kubernetes_resources` WHERE " +
					"account_name = \\?$").
					WillReturnResult(sqlmock.NewResult(1, 10))
				mock.ExpectCommit()

			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("#DeleteKubernetesResourcesByAccountName", func() {
		var name string

		BeforeEach(func() {
			name = "test-name"
		})

		JustBeforeEach(func() {
			err = c.DeleteKubernetesResourcesByAccountName(name)
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				mock.ExpectBegin()
				mock.ExpectExec("(?i)^DELETE FROM `kubernetes_resources` WHERE " +
					"account_name = \\?$").
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
				mock.ExpectQuery("(?i)^SELECT a.host, a.ca_data, a.bearer_token, a.token_provider, b.namespace FROM kubernetes_providers a " +
					"LEFT JOIN kubernetes_providers_namespaces b ON a.name = b.account_name " +
					"WHERE name = \\? ORDER BY `kubernetes_providers`.`name` LIMIT 1$").
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

	Describe("#GetKubernetesProviderAndPermissions", func() {
		var provider kubernetes.Provider

		BeforeEach(func() {
			provider = kubernetes.Provider{}
		})

		JustBeforeEach(func() {
			provider, err = c.GetKubernetesProviderAndPermissions("test-name")
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				sqlRows := sqlmock.NewRows([]string{"name", "host", "ca_data", "token_provider", "namespace", "read_group", "write_group"}).
					AddRow("test-name", "test-host", "test-ca-data", "test-token-provider", nil, "test-read-group", "test-write-group")
				mock.ExpectQuery("(?i)^SELECT a.name," +
					" a.host," +
					" a.ca_data," +
					" a.token_provider," +
					" a.namespace," +
					" b.read_group," +
					" c.write_group" +
					" FROM kubernetes_providers a" +
					" LEFT JOIN provider_read_permissions b ON a.name = b.account_name " +
					" LEFT JOIN provider_write_permissions c ON a.name = c.account_name" +
					" WHERE a.name = \\?").
					WillReturnRows(sqlRows)
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(provider.Name).To(Equal("test-name"))
				Expect(provider.Host).To(Equal("test-host"))
				Expect(provider.CAData).To(Equal("test-ca-data"))
				Expect(provider.TokenProvider).To(Equal("test-token-provider"))
				Expect(provider.Namespace).To(BeNil())
				Expect(provider.Permissions.Read[0]).To(Equal("test-read-group"))
				Expect(provider.Permissions.Write[0]).To(Equal("test-write-group"))
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
				mock.ExpectQuery("(?i)^SELECT " +
					"account_name, " +
					"cluster " +
					"FROM `kubernetes_resources` " +
					"WHERE spinnaker_app = \\? AND UPPER\\(kind\\) in \\('DEPLOYMENT', " +
					"'STATEFULSET', " +
					"'REPLICASET', " +
					"'INGRESS', " +
					"'SERVICE', " +
					"'DAEMONSET'\\) GROUP BY " +
					"account_name, cluster$").
					WillReturnRows(sqlRows)
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(resources).To(HaveLen(2))
			})
		})
	})

	Describe("#ListKubernetesClustersByFields", func() {
		var resources []kubernetes.Resource
		fields := []string{}

		BeforeEach(func() {
			fields = []string{"field1", "field2"}
		})

		JustBeforeEach(func() {
			resources, err = c.ListKubernetesClustersByFields(fields...)
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
				mock.ExpectQuery("(?i)^SELECT " +
					"field1, " +
					"field2 " +
					"FROM `kubernetes_resources` " +
					"WHERE UPPER\\(kind\\) in \\('DEPLOYMENT', 'STATEFULSET', 'REPLICASET', 'INGRESS', 'SERVICE', 'DAEMONSET'\\)" +
					" GROUP BY field1, field2$").
					WillReturnRows(sqlRows)
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(resources).To(HaveLen(2))
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
				mock.ExpectQuery("(?i)^SELECT " +
					"account_name, " +
					"api_group, " +
					"kind, " +
					"name, " +
					"artifact_name, " +
					"namespace, " +
					"resource, " +
					"task_type, " +
					"version " +
					"FROM `kubernetes_resources` " +
					" WHERE task_id = \\?$").
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
				sqlRows := sqlmock.NewRows([]string{"name", "host", "ca_data", "token_provider", "namespace"}).
					AddRow("name1", "host1", "ca_data1", "google", nil).
					AddRow("name2", "host2", "ca_data2", "rancher", nil)
				mock.ExpectQuery("(?i)^SELECT " +
					"name, " +
					"host, " +
					"ca_data, " +
					"token_provider, " +
					"namespace " +
					"FROM `kubernetes_providers`$").
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
				d.Close()
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
				mock.ExpectQuery("(?i)^SELECT " +
					"a.name, " +
					"a.host, " +
					"a.ca_data, " +
					"a.token_provider, " +
					"a.namespace, " +
					"b.read_group, " +
					"c.write_group " +
					"FROM kubernetes_providers a " +
					"left join provider_read_permissions b on a.name = b.account_name " +
					"left join provider_write_permissions c on a.name = c.account_name$").
					WillReturnRows(sqlRows)
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("sql: expected 4 destination arguments in Scan, not 7"))
			})
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				sqlRows := sqlmock.NewRows([]string{"name", "host", "ca_data", "google", "namespace", "read_group", "write_group"}).
					AddRow("name1", "host1", "ca_data1", "google", "namespace1", "read_group1", "write_group1").
					AddRow("name1", "host1", "ca_data1", "google", "namespace1", "read_group2", "write_group1").
					AddRow("name2", "host2", "ca_data2", "rancher", nil, "read_group2", "write_group2").
					AddRow("name2", "host2", "ca_data2", "rancher", nil, "read_group2", "write_group3")
				mock.ExpectQuery("(?i)^SELECT " +
					"a.name, " +
					"a.host, " +
					"a.ca_data, " +
					"a.token_provider, " +
					"a.namespace, " +
					"b.read_group, " +
					"c.write_group " +
					"FROM kubernetes_providers a " +
					"left join provider_read_permissions b on a.name = b.account_name " +
					"left join provider_write_permissions c on a.name = c.account_name$").
					WillReturnRows(sqlRows)
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				ns := "namespace1"
				Expect(err).To(BeNil())
				Expect(providers).To(HaveLen(2))
				Expect(providers[0].Namespace).To(Equal(&ns))
				Expect(providers[0].Permissions.Read).To(HaveLen(2))
				Expect(providers[0].Permissions.Write).To(HaveLen(1))
				Expect(providers[1].Namespace).To(BeNil())
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
				mock.ExpectQuery("(?i)^SELECT " +
					"field1, " +
					"field2 " +
					"FROM `kubernetes_resources` " +
					" GROUP BY field1, field2$").
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
				mock.ExpectQuery("(?i)^SELECT " +
					"`account_name` " +
					"FROM `kubernetes_resources` " +
					" WHERE spinnaker_app = \\? " +
					"GROUP BY `account_name`$").
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
				mock.ExpectQuery("(?i)^SELECT " +
					"`read_group` " +
					"FROM `provider_read_permissions` " +
					" WHERE account_name = \\? " +
					"GROUP BY `read_group`").
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
				mock.ExpectQuery("(?i)^SELECT " +
					"`write_group` " +
					"FROM `provider_write_permissions` " +
					" WHERE account_name = \\? " +
					"GROUP BY `write_group`").
					WillReturnRows(sqlRows)
				mock.ExpectCommit()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(groups).To(HaveLen(2))
			})
		})
	})
})
