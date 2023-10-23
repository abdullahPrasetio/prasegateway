package client

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Fungsi untuk menghapus kunci dari data JSON berdasarkan deny dan allow
// func FilterProcessJSON(jsonData *[]byte, deny []string, allow []string) {
// 	var data interface{}
// 	if err := json.Unmarshal(*jsonData, &data); err != nil {
// 		panic(err)
// 	}

// 	filteredData := make(map[string]interface{})
// 	switch obj := data.(type) {
// 	case map[string]interface{}:
// 		// Iterasi melalui array allow untuk memfilter data
// 		for _, key := range allow {
// 			parts := strings.Split(key, ".")
// 			currentData := obj

// 			for _, part := range parts[:len(parts)-1] {
// 				value, ok := currentData[part]
// 				if !ok {
// 					break
// 				}
// 				currentData = value.(map[string]interface{})
// 			}

// 			// Jika key ditemukan, tambahkan ke hasil filter
// 			lastPart := parts[len(parts)-1]
// 			if value, ok := currentData[lastPart]; ok {
// 				filteredData[lastPart] = value
// 			}
// 		}

// 		// Mengganti data JSON dengan data yang telah diproses
// 		*jsonData = []byte{}
// 		if processedData, err := json.Marshal(filteredData); err != nil {
// 			panic(err)
// 		} else {
// 			*jsonData = processedData
// 		}
// 	default:
// 		return
// 	}
// }

func FilterProcessJSON(jsonData *[]byte, deny []string, allow []string) {
	var data interface{}
	if err := json.Unmarshal(*jsonData, &data); err != nil {
		// panic(err)
		fmt.Println(err)
	}
	filterDataRecursively(data, deny, allow)
	// Mengencode data yang telah diproses kembali ke JSON
	processedData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// Mengganti data JSON dengan data yang telah diproses
	*jsonData = processedData
}

func filterDataRecursively(data interface{}, deny []string, allow []string) {
	switch obj := data.(type) {
	case map[string]interface{}:
		// Memproses allow
		for _, key := range allow {
			keys := strings.Split(key, ".")
			value, ok := getNestedFieldRecursive(obj, keys)
			if ok {
				setNestedFieldRecursive(obj, key, keys[len(keys)-1], value, keys)
			}
		}

		// Memproses deny (menghapus field yang tidak diizinkan)
		for _, key := range deny {
			keys := strings.Split(key, ".")
			deleteNestedFieldRecursive(obj, keys)
		}
	case []interface{}:
		for _, val := range obj {
			filterDataRecursively(val, deny, allow)
		}
	}
}
