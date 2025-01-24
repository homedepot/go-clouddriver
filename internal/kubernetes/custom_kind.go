package kubernetes

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/homedepot/go-clouddriver/internal/kubernetes/manifest"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type CustomKindConfig struct {
	StatusChecks []StatusCheck `json:"statusChecks"`
}

type StatusCheck struct {
	FieldName  string      `json:"fieldName"`
	FieldValue interface{} `json:"fieldValue"`
}

type CustomKind struct {
	CustomKindConfig
	manifest *unstructured.Unstructured
}

func NewCustomKind(kind string, m map[string]interface{}) *CustomKind {
	manifest, err := ToUnstructured(m)
	if err != nil {
		clouddriver.Log(fmt.Errorf("error creating unstructured object from manifest: %v", err))
	}

	configData := getCustomKindConfig(kind)

	return &CustomKind{manifest: &manifest, CustomKindConfig: configData}
}

func (k *CustomKind) Object() *unstructured.Unstructured {
	return k.manifest
}

func (k *CustomKind) Status() manifest.Status {
	log.Println(*k)

	s := manifest.DefaultStatus

	unstructuredContent := k.manifest.UnstructuredContent()
	if _, ok := unstructuredContent["status"]; !ok {
		return s
	}

	statusData := unstructuredContent["status"].(map[string]interface{})

	for _, statusCheck := range k.StatusChecks {
		if statusData[statusCheck.FieldName] != statusCheck.FieldValue {
			s.Stable.State = false
			s.Stable.Message = fmt.Sprintf("Waiting for %s to be %s", statusCheck.FieldName, statusCheck.FieldValue)

			return s
		}
	}

	return s
}

func getCustomKindConfig(kind string) CustomKindConfig {
	customKindsConfigPath := os.Getenv("CUSTOM_KINDS_CONFIG_PATH")
	allConfigs := map[string]CustomKindConfig{}

	if customKindsConfigPath == "" {
		return CustomKindConfig{}
	}

	configBytes, err := os.ReadFile(customKindsConfigPath)
	if err != nil {
		clouddriver.Log(fmt.Errorf("error reading custom kinds config file at %s: %v",
			customKindsConfigPath, err))
	}

	if err := json.Unmarshal(configBytes, &allConfigs); err != nil {
		clouddriver.Log(fmt.Errorf("error setting up custom kinds config: %v", err))
	}

	log.Println(allConfigs)

	config, ok := allConfigs[kind]
	if !ok {
		return CustomKindConfig{}
	}

	return config
}

// {
//   "TinyhomeDeployment": {
//     "statusChecks": [
//       {
// 	       "fieldName": "ready",
// 	       "fieldValue": true
//       }
//     ]
//   }
// }
