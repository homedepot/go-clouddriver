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
						FieldPath:     "phase.type",
						ComparedValue: "Error",
						Operator:      "NE",
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
				status := map[string]interface{}{"ready": false}
				fakeManifest := map[string]interface{}{"status": status, "kind": "doesnotexist"}
				customKind = NewCustomKind("doesnotexist", fakeManifest)
			})

			It("returns default status", func() {
				Expect(s).To(Equal(manifest.DefaultStatus))
			})
		})

		When("the status check fails", func() {
			BeforeEach(func() {
				status := map[string]interface{}{"phase": map[string]interface{}{"type": "Error", "reason": "some reason"}}
				fakeManifest := map[string]interface{}{"status": status, "kind": "test"}
				customKind = NewCustomKind("test", fakeManifest)
			})

			It("returns status failed", func() {
				Expect(s.Failed.State).To(BeTrue())
				Expect(s.Failed.Message).To(Equal("Field status.phase.type was Error"))
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

		When("the status check succeeds", func() {
			BeforeEach(func() {
				status := map[string]interface{}{"phase": map[string]interface{}{"type": "some type"}}
				fakeManifest := map[string]interface{}{"status": status, "kind": "test"}
				customKind = NewCustomKind("test", fakeManifest)
			})

			It("returns the expected status", func() {
				Expect(s.Stable.State).To(BeTrue())
				Expect(s.Failed.State).To(BeFalse())
			})
		})
	})
})
