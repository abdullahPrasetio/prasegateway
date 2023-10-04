package utils

import "time"

/********************************************************************************
* Temancode Example Response Package                                            *
*                                                                               *
* Version: 1.0.0                                                                *
* Date:    2023-01-05                                                           *
* Author:  Waluyo Ade Prasetio                                                  *
* Github:  https://github.com/abdullahPrasetio                                  *
********************************************************************************/

type HealthCheckResp struct {
	NameApp   string `json:"name_app"`
	Timestamp string `json:"timestamp"`
}

func HealthCheckResponse(nameApp string) HealthCheckResp {
	return HealthCheckResp{
		NameApp:   nameApp,
		Timestamp: time.Now().Local().Format(time.RFC850),
	}
}
