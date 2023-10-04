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
			headers := []client.Headers{}
			for key, values := range c.Request.Header {
				for _, value := range values {
					headers = append(headers, client.Headers{
						Key:   key,
						Value: value,
					})
				}
			}
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
			for key, values := range respHeader {
				for _, value := range values {
					c.Header(key, value)
				}
			}

			c.Header("Cache-Control", "no-cache")
			// c.Header("If-None-Match")

			c.Data(http.StatusOK, "application/json", responseBody)
			// Handler Anda untuk metode GET di sini
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
