package routers

import (
	"fmt"
	"log"
	"net/http"
	"plugin"
	"strings"
	"time"

	client "github.com/abdullahPrasetio/prasegateway/clients"
	"github.com/abdullahPrasetio/prasegateway/entity"
	"github.com/abdullahPrasetio/prasegateway/utils"
	"github.com/gin-gonic/gin"
)

type MyPlugin interface {
	Apply(*gin.Engine)
}

func Setup(myConfig entity.MyConfig) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		ifNoneMatchHeader := c.Request.Header.Get("If-None-Match")

		// Lakukan validasi sesuai kebutuhan Anda
		// Misalnya, Anda dapat memeriksa apakah ETag yang diberikan sama dengan versi terbaru dari sumber daya
		// Jika sama, Anda dapat mengirimkan respons StatusNotModified (HTTP 304)
		c.Request.Header.Add("If-None-Match", "")
		// Simulasi validasi ETag
		if ifNoneMatchHeader == "12345" {
			c.Status(http.StatusNotModified)
			c.Abort() // Abort middleware agar endpoint tidak dijalankan
			return
		}

		fmt.Println("masuk sini", ifNoneMatchHeader)
		c.Next()
	})
	r.Use(gin.Recovery())

	// Memuat plugin dari file yang sudah dikompilasi
	p, err := plugin.Open("plugins/correlation_id_plugin.so")
	if err != nil {
		fmt.Println("Error loading plugin:", err)

	}

	// Mencari simbol yang mengimplementasikan MyPlugin
	symbol, err := p.Lookup("ExportedPlugin")
	if err != nil {
		fmt.Println("Error looking up symbol:", err)

	}

	// Mencast simbol menjadi MyPlugin
	myPlugin, ok := symbol.(MyPlugin)
	if !ok {
		fmt.Println("Invalid plugin format")

	}

	// Menerapkan plugin ke router Gin
	myPlugin.Apply(r)

	r.GET("healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, utils.HealthCheckResponse(myConfig.NameApp))
	})

	for _, service := range myConfig.Services {
		myService := service
		for _, route := range myService.Endpoints {
			myRoute := route
			getMethodHandler(myService, myRoute, r)
		}

	}

	return r
}

func getMethodHandler(service entity.Service, endpoint entity.Endpoint, r *gin.Engine) {
	rg := r.Group(service.Prefix)
	switch endpoint.Method {
	case "GET":
		rg.GET(endpoint.Path, func(c *gin.Context) {

			startTime := time.Now() // Waktu awal
			headers := client.GetHeaderRequest(c)
			pathParts := strings.Split(c.Request.URL.Path, service.Prefix)
			path := ""
			// Ambil query parameter dari URL
			query := c.Request.URL.Query()
			// Loop melalui query parameter dan tambahkan ke URI
			for key, values := range query {
				for _, value := range values {
					path += fmt.Sprintf("&%s=%s", key, value)
				}
			}
			uri := service.BaseURL + endpoint.Destination + pathParts[1] + "?" + path
			// if len(endpoint.Destination) > 0 {
			// 	uri = service.BaseURL + endpoint.Destination
			// }
			var body []byte
			responseBody, respHeader, err := client.Client_Req(c, headers, uri, endpoint.Method, body)

			// Menghitung waktu yang diperlukan
			elapsedTime := time.Since(startTime)
			// Log waktu yang diperlukan
			log.Printf("Waktu yang diperlukan untuk panggilan API: %v", elapsedTime)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim permintaan" + err.Error()})
				return
			}
			client.SetHttpHeaderResponse(c, respHeader)

			// bodyResponse := client.MappingWithByteReplace(c, responseBody, endpoint)
			fmt.Println("body before", string(responseBody))
			client.FilterProcessJSON(&responseBody, endpoint.Response.Body.Deny, endpoint.Response.Body.Allow)
			client.MappingNestedRecursive(&responseBody, endpoint)
			fmt.Println("body after", string(responseBody))
			// bodyResponse = client.MappingNestedRecursive(bodyResponse, endpoint)

			c.Data(http.StatusOK, "application/json", responseBody)
		})
	case "POST":
		rg.POST(endpoint.Path, func(c *gin.Context) {
			// Handler Anda untuk metode POST di sini
		})
	// Tambahkan jenis metode lain sesuai dengan kebutuhan Anda
	default:
		log.Fatalf("Metode HTTP tidak valid untuk rute %s: %s", endpoint.Path, endpoint.Method)
	}
}
