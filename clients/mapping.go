package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/abdullahPrasetio/prasegateway/entity"
	"github.com/gin-gonic/gin"
)

func MappingWithByteReplace(c *gin.Context, responseBody []byte, endpoint entity.Endpoint) (body []byte) {
	body = responseBody
	if len(endpoint.Response.Mapping) > 0 {
		modifiedResponse := responseBody
		for oldKey, newKey := range endpoint.Response.Mapping {
			oldKeyBytes := []byte("\"" + oldKey + "\":")
			newKeyBytes := []byte("\"" + newKey + "\":")

			modifiedResponse = bytes.Replace(modifiedResponse, oldKeyBytes, newKeyBytes, -1)
		}

		// c.Data(http.StatusOK, "application/json", modifiedResponse)
		body = modifiedResponse
		return
	}
	return
}

func MappingNested(jsonData []byte, endpoint entity.Endpoint) []byte {
	modifiedJSON := jsonData
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return modifiedJSON
	}

	if len(endpoint.Response.Mapping) > 0 {
		for oldKey, newKey := range endpoint.Response.Mapping {
			if strings.Contains(oldKey, ".") {
				// Ini adalah pemetaan untuk nested field
				if val, ok := getNestedField(data, oldKey); ok {
					// Hapus oldKey
					deleteNestedField(data, oldKey)
					// Tambahkan newKey dengan nilai yang sama
					setNestedField(data, newKey, val)
				}
			} else {
				// Ini adalah pemetaan untuk field biasa
				if val, ok := data[oldKey]; ok {
					delete(data, oldKey)
					data[newKey] = val
				}
			}
		}
	}

	// Marshal map kembali ke JSON
	modifiedJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return modifiedJSON
	}

	return modifiedJSON
}

// Fungsi untuk mengambil nilai dari nested field (contoh: "address.geo.lat")
func getNestedField(data map[string]interface{}, fieldPath string) (interface{}, bool) {
	keys := strings.Split(fieldPath, ".")
	current := data
	for _, key := range keys {
		value, ok := current[key]
		if !ok {
			return nil, false
		}
		if obj, ok := value.(map[string]interface{}); ok {
			current = obj
		} else {
			return value, true
		}
	}
	return nil, false
}

// Fungsi untuk mengatur nilai pada nested field (contoh: "address.geo.lat")
func setNestedField(data map[string]interface{}, fieldPath string, value interface{}) {
	keys := strings.Split(fieldPath, ".")
	current := data
	for i, key := range keys {
		if i == len(keys)-1 {
			current[key] = value
		} else {
			if obj, ok := current[key].(map[string]interface{}); ok {
				current = obj
			} else {
				obj = make(map[string]interface{})
				current[key] = obj
				current = obj
			}
		}
	}
}

// Fungsi untuk menghapus nilai pada nested field (contoh: "address.geo.lat")
func deleteNestedField(data map[string]interface{}, fieldPath string) {
	keys := strings.Split(fieldPath, ".")
	current := data
	for i, key := range keys {
		if i == len(keys)-1 {
			delete(current, key)
		} else {
			if obj, ok := current[key].(map[string]interface{}); ok {
				current = obj
			} else {
				return
			}
		}
	}
}

func MappingNestedRecursive(jsonData []byte, endpoint entity.Endpoint) []byte {
	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return jsonData
	}

	// Jalankan pemrosesan rekursif
	modifyDataRecursively(data, endpoint.Response.Mapping)

	modifiedJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return jsonData
	}

	return modifiedJSON
}

func modifyDataRecursively(data interface{}, mapping map[string]string) {
	switch obj := data.(type) {
	case map[string]interface{}:
		for oldKey, newKey := range mapping {
			if strings.Contains(oldKey, ".") {
				keys := strings.Split(oldKey, ".")
				lastKey := keys[len(keys)-1]

				if nestedData, ok := getNestedFieldRecursive(obj, keys); ok {
					deleteNestedFieldRecursive(obj, keys)
					setNestedFieldRecursive(obj, newKey, lastKey, nestedData, keys)
				}
			} else {
				if val, ok := obj[oldKey]; ok {
					delete(obj, oldKey)
					obj[newKey] = val
				}
			}
		}
		for _, val := range obj {
			modifyDataRecursively(val, mapping)
		}
	case []interface{}:
		for _, val := range obj {
			modifyDataRecursively(val, mapping)
		}
	}
}

func getNestedFieldRecursive(data map[string]interface{}, keys []string) (interface{}, bool) {
	current := data
	for _, key := range keys {
		if val, ok := current[key]; ok {
			if obj, ok := val.(map[string]interface{}); ok {
				current = obj
			} else {
				return val, true
			}
		} else {
			return nil, false
		}
	}
	return nil, false
}

func deleteNestedFieldRecursive(data map[string]interface{}, keys []string) {
	current := data
	lastIdx := len(keys) - 1
	for i, key := range keys {
		if i == lastIdx {
			delete(current, key)
		} else {
			if obj, ok := current[key].(map[string]interface{}); ok {
				current = obj
			} else {
				return
			}
		}
	}
}

func setNestedFieldRecursive(data map[string]interface{}, newKey, lastKey string, value interface{}, oldKey []string) {
	keys := strings.Split(newKey, ".")
	current := data
	for i, key := range keys {
		if i == len(keys)-1 {
			current[key] = value
		} else {
			if obj, ok := current[key].(map[string]interface{}); ok {
				current = obj
			} else {
				obj = make(map[string]interface{})
				current[key] = obj
				current = obj
			}
		}
	}

	deleteNestedFieldRecursive(data, oldKey)
}
