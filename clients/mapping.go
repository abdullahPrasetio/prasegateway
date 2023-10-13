package client

import (
	"bytes"
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
