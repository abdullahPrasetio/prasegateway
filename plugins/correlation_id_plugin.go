package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MyPlugin adalah antarmuka untuk plugin
type MyPlugin interface {
	Apply(*gin.Engine)
}

// CorrelationIDPlugin adalah implementasi dari MyPlugin
type CorrelationIDPlugin struct{}

// Apply adalah metode yang menambahkan header "Correlation-ID" ke setiap respons
func (p CorrelationIDPlugin) Apply(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		// Menghasilkan UUID baru untuk Correlation-ID
		correlationID := uuid.New().String()

		// Menambahkan header "Correlation-ID" ke respons
		c.Header("Correlation-ID", correlationID)

		// Lanjutkan ke handler berikutnya dalam rantai middleware
		c.Next()
	})
}

// Main adalah fungsi entry point untuk kompilasi plugin
func main() {}

// ExportedPlugin adalah variabel yang diekspor dan akan digunakan oleh aplikasi utama
var ExportedPlugin CorrelationIDPlugin
