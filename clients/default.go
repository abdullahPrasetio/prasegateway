package client

/********************************************************************************
* Temancode Example Request Http Package                                        *
*                                                                               *
* Version: 1.0.0                                                                *
* Date:    2023-01-05                                                           *
* Author:  Waluyo Ade Prasetio                                                  *
* Github:  https://github.com/abdullahPrasetio                                  *
********************************************************************************/

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/abdullahPrasetio/prasegateway/config"
	"github.com/abdullahPrasetio/prasegateway/log"
	"github.com/abdullahPrasetio/prasegateway/utils"
)

type Headers struct {
	Key   string
	Value string
}

func Client_Req(ctx context.Context, headers []Headers, uri string, method string, bodyRequest []byte) ([]byte, http.Header, error) {

	config := config.GetMyConfig()

	client := &http.Client{Timeout: time.Second * time.Duration(config.TTL)}

	// logger := log.Logger

	// logger.Info(string(bodyRequest))

	refnum, _ := ctx.Value("refnum").(string)
	req, err := http.NewRequest(method, uri, bytes.NewReader(bodyRequest))
	if err != nil {
		return nil, nil, err
	}

	// fmt.Println("headers", headers)
	headerDeny := []string{"Accept-Encoding"}
	if headers == nil {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "application/json")
		//header tambahan
		for i := range headers {
			if !utils.IsStringInSlice(headers[i].Key, headerDeny) {
				req.Header[headers[i].Key] = []string{headers[i].Value}
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	//fmt.Println(err)
	if err != nil {
		return nil, nil, err
	}

	//check response status
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode >= http.StatusInternalServerError {
			return nil, nil, errors.New(fmt.Sprintf("%d error: %s", resp.StatusCode, resp.Status))
		}
		return nil, nil, fmt.Errorf("%d error: %s", resp.StatusCode, resp.Status)
	}

	respHeader, _ := json.Marshal(resp.Header)

	log.LogDebug(refnum, "response header = "+string(respHeader))
	log.LogDebug(refnum, "response body = "+string(body))

	return body, resp.Header, nil
}

func readBodyString(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
