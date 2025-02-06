package kubernetes

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/homedepot/go-clouddriver/internal/kubernetes/manifest"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type CustomKindConfig struct {
	StatusChecks []StatusCheck `json:"statusChecks"`
}

type StatusCheck struct {
	FieldPath     string      `json:"fieldPath"`
	ComparedValue interface{} `json:"comparedValue"`
	Operator      string      `json:"operator"`
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
	log.Println(k.manifest.Object)
	log.Println(k.CustomKindConfig)

	s := manifest.DefaultStatus

	unstructuredContent := k.manifest.UnstructuredContent()
	if _, ok := unstructuredContent["status"]; !ok {
		return s
	}

	statusData := unstructuredContent["status"].(map[string]interface{})

	for _, statusCheck := range k.StatusChecks {
		statusValue := getStatusValue(statusData, statusCheck.FieldPath)
		if statusValue == nil {
			continue
		}

		if !evaluatestatusCheck(statusValue, statusCheck.ComparedValue, statusCheck.Operator) {
			s.Stable.State = false
			s.Failed.State = true
			s.Failed.Message = fmt.Sprintf("Field status.%s was %v", statusCheck.FieldPath, statusValue)

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

func getStatusValue(statusMap map[string]interface{}, fieldPath string) interface{} {
	fields := strings.Split(fieldPath, ".")
	if len(fields) == 0 {
		return nil
	}

	currField := fields[0]

	if len(fields) == 1 {
		val, exists := statusMap[currField]
		if !exists {
			return nil
		}

		return val
	}

	remainingFields := fields[1:]

	val, exists := statusMap[currField]
	if !exists {
		return nil
	}

	return getStatusValue(val.(map[string]interface{}), strings.Join(remainingFields, "."))
}

func evaluatestatusCheck(actual, compared interface{}, operator string) bool {
	switch strings.ToLower(operator) {
	case "eq":
		return actual == compared
	case "ne":
		return actual != compared
	default:
		return true
	}
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
