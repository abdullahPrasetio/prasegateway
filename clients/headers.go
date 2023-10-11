package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHeaderRequest(c *gin.Context) []Headers {
	headers := []Headers{}
	for key, values := range c.Request.Header {
		for _, value := range values {
			headers = append(headers, Headers{
				Key:   key,
				Value: value,
			})
		}
	}

	return headers
}

func SetHttpHeaderResponse(c *gin.Context, respHeader http.Header) {

	for key, values := range respHeader {
		for _, value := range values {
			c.Header(key, value)
		}
	}
}
