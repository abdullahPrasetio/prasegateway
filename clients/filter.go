package client

func ProcessJSON(data map[string]interface{}, deny, allow []string) {
	// Iterasi melalui elemen JSON
	for key, value := range data {
		if inSlice(key, deny) {
			// Jika dalam daftar deny, hapus key tersebut
			delete(data, key)
		} else if len(allow) == 0 || inSlice(key, allow) {
			// Jika dalam daftar allow (atau daftar allow kosong), biarkan key tersebut
			if subdata, ok := value.(map[string]interface{}); ok {
				// Jika value adalah objek map, proses rekursif
				ProcessJSON(subdata, deny, allow)
			}
		}
	}
}

func inSlice(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}
