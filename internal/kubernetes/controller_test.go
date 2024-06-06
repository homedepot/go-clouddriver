package kubernetes

import (
	"time"

	"github.com/homedepot/go-clouddriver/internal/kubernetes/cached/memory"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"
)

const (
	rootCACert = `-----BEGIN CERTIFICATE-----
MIIC4DCCAcqgAwIBAgIBATALBgkqhkiG9w0BAQswIzEhMB8GA1UEAwwYMTAuMTMu
MTI5LjEwNkAxNDIxMzU5MDU4MB4XDTE1MDExNTIxNTczN1oXDTE2MDExNTIxNTcz
OFowIzEhMB8GA1UEAwwYMTAuMTMuMTI5LjEwNkAxNDIxMzU5MDU4MIIBIjANBgkq
hkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAunDRXGwsiYWGFDlWH6kjGun+PshDGeZX
xtx9lUnL8pIRWH3wX6f13PO9sktaOWW0T0mlo6k2bMlSLlSZgG9H6og0W6gLS3vq
s4VavZ6DbXIwemZG2vbRwsvR+t4G6Nbwelm6F8RFnA1Fwt428pavmNQ/wgYzo+T1
1eS+HiN4ACnSoDSx3QRWcgBkB1g6VReofVjx63i0J+w8Q/41L9GUuLqquFxu6ZnH
60vTB55lHgFiDLjA1FkEz2dGvGh/wtnFlRvjaPC54JH2K1mPYAUXTreoeJtLJKX0
ycoiyB24+zGCniUmgIsmQWRPaOPircexCp1BOeze82BT1LCZNTVaxQIDAQABoyMw
ITAOBgNVHQ8BAf8EBAMCAKQwDwYDVR0TAQH/BAUwAwEB/zALBgkqhkiG9w0BAQsD
ggEBADMxsUuAFlsYDpF4fRCzXXwrhbtj4oQwcHpbu+rnOPHCZupiafzZpDu+rw4x
YGPnCb594bRTQn4pAu3Ac18NbLD5pV3uioAkv8oPkgr8aUhXqiv7KdDiaWm6sbAL
EHiXVBBAFvQws10HMqMoKtO8f1XDNAUkWduakR/U6yMgvOPwS7xl0eUTqyRB6zGb
K55q2dejiFWaFqB/y78txzvz6UlOZKE44g2JAVoJVM6kGaxh33q8/FmrL4kuN3ut
W+MmJCVDvd4eEqPwbp7146ZWTqpIJ8lvA6wuChtqV8lhAPka2hD/LMqY8iXNmfXD
uml0obOEy+ON91k+SWTJ3ggmF/U=
-----END CERTIFICATE-----`

	certData = `-----BEGIN CERTIFICATE-----
MIIC6jCCAdSgAwIBAgIBCzALBgkqhkiG9w0BAQswIzEhMB8GA1UEAwwYMTAuMTMu
MTI5LjEwNkAxNDIxMzU5MDU4MB4XDTE1MDExNTIyMDEzMVoXDTE2MDExNTIyMDEz
MlowGzEZMBcGA1UEAxMQb3BlbnNoaWZ0LWNsaWVudDCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBAKtdhz0+uCLXw5cSYns9rU/XifFSpb/x24WDdrm72S/v
b9BPYsAStiP148buylr1SOuNi8sTAZmlVDDIpIVwMLff+o2rKYDicn9fjbrTxTOj
lI4pHJBH+JU3AJ0tbajupioh70jwFS0oYpwtneg2zcnE2Z4l6mhrj2okrc5Q1/X2
I2HChtIU4JYTisObtin10QKJX01CLfYXJLa8upWzKZ4/GOcHG+eAV3jXWoXidtjb
1Usw70amoTZ6mIVCkiu1QwCoa8+ycojGfZhvqMsAp1536ZcCul+Na+AbCv4zKS7F
kQQaImVrXdUiFansIoofGlw/JNuoKK6ssVpS5Ic3pgcCAwEAAaM1MDMwDgYDVR0P
AQH/BAQDAgCgMBMGA1UdJQQMMAoGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwCwYJ
KoZIhvcNAQELA4IBAQCKLREH7bXtXtZ+8vI6cjD7W3QikiArGqbl36bAhhWsJLp/
p/ndKz39iFNaiZ3GlwIURWOOKx3y3GA0x9m8FR+Llthf0EQ8sUjnwaknWs0Y6DQ3
jjPFZOpV3KPCFrdMJ3++E3MgwFC/Ih/N2ebFX9EcV9Vcc6oVWMdwT0fsrhu683rq
6GSR/3iVX1G/pmOiuaR0fNUaCyCfYrnI4zHBDgSfnlm3vIvN2lrsR/DQBakNL8DJ
HBgKxMGeUPoneBv+c8DMXIL0EhaFXRlBv9QW45/GiAIOuyFJ0i6hCtGZpJjq4OpQ
BRjCI+izPzFTjsxD4aORE+WOkyWFCGPWKfNejfw0
-----END CERTIFICATE-----`
)

var _ = Describe("Controller", func() {
	var (
		client     Client
		clientset  Clientset
		config     *rest.Config
		controller Controller
		err        error
	)

	Describe("#NewClient", func() {
		BeforeEach(func() {
			memCaches = map[string]*memory.Cache{}
			config = &rest.Config{
				Host:        "https://test-host",
				BearerToken: "some.bearer.token",
				TLSClientConfig: rest.TLSClientConfig{
					CAData: []byte(rootCACert),
				},
			}
			controller = NewController()
		})

		JustBeforeEach(func() {
			client, err = controller.NewClient(config)
		})

		Context("memory cache", func() {
			When("generating the dynamic client returns an error", func() {
				BeforeEach(func() {
					config = &rest.Config{
						Host:        ":::badhost;",
						BearerToken: "some.bearer.token",
						TLSClientConfig: rest.TLSClientConfig{
							CAData: []byte(rootCACert),
						},
					}
				})

				It("returns an error", func() {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("parse \"https://:::badhost;\": invalid port \":badhost;\" after host"))
				})
			})

			When("a call is made for a cached client", func() {
				JustBeforeEach(func() {
					client, err = controller.NewClient(config)
				})

				It("creates a mem cache for the client", func() {
					Expect(err).To(BeNil())
					Expect(client).ToNot(BeNil())
					Expect(memCaches).To(HaveLen(1))
					memCache := memCaches[config.Host]
					Expect(memCache).ToNot(BeNil())
				})
			})

			When("the bearer token for a client changes", func() {
				JustBeforeEach(func() {
					newConfig := &rest.Config{
						Host:        "https://test-host",
						BearerToken: "another.bearer.token",
						TLSClientConfig: rest.TLSClientConfig{
							CAData: []byte(rootCACert),
						},
					}
					client, err = controller.NewClient(newConfig)
				})

				It("references the same cache instance", func() {
					Expect(err).To(BeNil())
					Expect(client).ToNot(BeNil())
					Expect(memCaches).To(HaveLen(1))
					memCache := memCaches[config.Host]
					Expect(memCache).ToNot(BeNil())
				})
			})

			When("the CAData for a client changes", func() {
				JustBeforeEach(func() {
					newConfig := &rest.Config{
						Host:        "https://test-host",
						BearerToken: "some.bearer.token",
						TLSClientConfig: rest.TLSClientConfig{
							CAData: []byte(certData),
						},
					}
					client, err = controller.NewClient(newConfig)
				})

				It("references the same cache instance", func() {
					Expect(err).To(BeNil())
					Expect(client).ToNot(BeNil())
					Expect(memCaches).To(HaveLen(1))
					memCache := memCaches[config.Host]
					Expect(memCache).ToNot(BeNil())
				})
			})

			When("the same host has two defined timeouts", func() {
				JustBeforeEach(func() {
					newConfig := &rest.Config{
						Host:        "https://test-host",
						BearerToken: "some.bearer.token",
						TLSClientConfig: rest.TLSClientConfig{
							CAData: []byte(rootCACert),
						},
						Timeout: 1 * time.Second,
					}
					client, err = controller.NewClient(newConfig)
				})

				It("references the same cache instance", func() {
					Expect(err).To(BeNil())
					Expect(client).ToNot(BeNil())
					Expect(memCaches).To(HaveLen(1))
				})
			})

			It("creates a mem cache and generates a client", func() {
				Expect(err).To(BeNil())
				Expect(client).ToNot(BeNil())
				Expect(memCaches).To(HaveLen(1))
				memCache := memCaches[config.Host]
				Expect(memCache).ToNot(BeNil())
			})
		})
	})

	Describe("#NewClientset", func() {
		BeforeEach(func() {
			config = &rest.Config{
				Host:        "https://test-host",
				BearerToken: "some.bearer.token",
				TLSClientConfig: rest.TLSClientConfig{
					CAData: []byte(rootCACert),
				},
			}
			controller = NewController()
		})

		JustBeforeEach(func() {
			clientset, err = controller.NewClientset(config)
		})

		Context("memory cache", func() {
			When("generating the clientset returns an error", func() {
				BeforeEach(func() {
					config = &rest.Config{
						Host:        ":::badhost;",
						BearerToken: "some.bearer.token",
						TLSClientConfig: rest.TLSClientConfig{
							CAData: []byte(rootCACert),
						},
					}
				})

				It("returns an error", func() {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("parse \"https://:::badhost;\": invalid port \":badhost;\" after host"))
				})
			})

			It("returns the clientset", func() {
				Expect(err).To(BeNil())
				Expect(clientset).ToNot(BeNil())
			})
		})
	})
})
