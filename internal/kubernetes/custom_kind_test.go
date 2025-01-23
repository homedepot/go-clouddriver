package kubernetes_test

import (
	"encoding/json"
	"os"

	. "github.com/homedepot/go-clouddriver/internal/kubernetes"
	"github.com/homedepot/go-clouddriver/internal/kubernetes/manifest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Custom Kind", Ordered, func() {
	var (
		customKind           *CustomKind
		fakeCustomKindConfig map[string]CustomKindConfig
	)

	BeforeAll(func() {
		fakeCustomKindConfig = map[string]CustomKindConfig{
			"test": {
				StatusChecks: []StatusCheck{
					{
						FieldName:  "ready",
						FieldValue: "true",
					},
				},
			},
		}
		customKindsConfigPath := "./customKindsConfig.json"
		os.Setenv("CUSTOM_KINDS_CONFIG_PATH", customKindsConfigPath)
		f, err := os.Create(customKindsConfigPath)
		Expect(err).To(BeNil())
		err = json.NewEncoder(f).Encode(fakeCustomKindConfig)
		Expect(err).To(BeNil())
		f.Close()
	})

	AfterAll(func() {
		os.Remove(os.Getenv("CUSTOM_KINDS_CONFIG_PATH"))
	})

	Describe("#Object", func() {
		BeforeEach(func() {
			fakeManifest := map[string]interface{}{"status": "1", "kind": "test"}
			customKind = NewCustomKind("test", fakeManifest)
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(customKind.Object()).ToNot(BeNil())
			})
		})
	})

	Describe("#Status", func() {
		var s manifest.Status

		JustBeforeEach(func() {
			s = customKind.Status()
		})

		When("there is no configuration for the kind", func() {
			BeforeEach(func() {
				status := map[string]interface{}{"ready": "false"}
				fakeManifest := map[string]interface{}{"status": status, "kind": "doesnotexist"}
				customKind = NewCustomKind("doesnotexist", fakeManifest)
			})

			It("returns default status", func() {
				Expect(s).To(Equal(manifest.DefaultStatus))
			})
		})

		When("the values do not match", func() {
			BeforeEach(func() {
				status := map[string]interface{}{"ready": "false"}
				fakeManifest := map[string]interface{}{"status": status, "kind": "test"}
				customKind = NewCustomKind("test", fakeManifest)
			})

			It("returns status unstable", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for ready to be true"))
			})
		})

		When("the expected field does not exist", func() {
			BeforeEach(func() {
				fakeManifest := map[string]interface{}{"randomField": "1", "kind": "test"}
				customKind = NewCustomKind("test", fakeManifest)
			})

			It("returns default status", func() {
				Expect(s).To(Equal(manifest.DefaultStatus))
			})
		})

		When("the expected field exists and values match", func() {
			BeforeEach(func() {
				status := map[string]interface{}{"ready": "true"}
				fakeManifest := map[string]interface{}{"status": status, "kind": "test"}
				customKind = NewCustomKind("test", fakeManifest)
			})

			It("returns the expected status", func() {
				Expect(s.Stable.State).To(BeTrue())
			})
		})
	})
})
