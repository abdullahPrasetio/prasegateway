package client

import (
	"encoding/json"
	"fmt"

	"github.com/abdullahPrasetio/prasegateway/entity"
)

func GlobalProcessResponse(jsonData *[]byte, endpoint entity.Endpoint) {
	var data interface{}
	if err := json.Unmarshal(*jsonData, &data); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		panic(err)
	}
	// modifyDataRecursivelyWithFilter()
	// Memanggil fungsi untuk menyimpan kunci berdasarkan allow
	// deny := endpoint.Response.Body.Deny
	// allow := endpoint.Response.Body.Allow
	// Memanggil fungsi untuk menyimpan kunci berdasarkan allow
	// keepKeys(data, deny, allow)

	// Memanggil fungsi untuk menghapus kunci berdasarkan deny
	// removeKeys(data, deny)

	// Jalankan pemrosesan rekursif
	// modifyDataRecursively(data, endpoint.Response.Mapping)

	modifiedJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		panic(err)
	}

	*jsonData = modifiedJSON
}
