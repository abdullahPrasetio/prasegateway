package entity

type MyConfig struct {
	NameApp  string    `json:"name_app"`
	Auth     bool      `json:"auth"`
	Debug    bool      `json:"debug"`
	AppENV   string    `json:"app_env"`
	TTL      int64     `json:"ttl"`
	Services []Service `json:"services"`
}

type Service struct {
	Prefix    string     `json:"prefix"`
	BaseURL   string     `json:"base_url"`
	Endpoints []Endpoint `json:"endpoints"`
	All       bool       `json:"all"`
}

type Endpoint struct {
	Destination string          `json:"destination"`
	Path        string          `json:"path"`
	Method      string          `json:"method"`
	Request     RequestResponse `json:"request"`
	Response    RequestResponse `json:"response"`
}

type RequestResponse struct {
	Headers Header            `json:"headers"`
	Body    Body              `json:"body"`
	Mapping map[string]string `json:"mapping"`
}

type Header struct {
	Allow []string `json:"allow"`
	Deny  []string `json:"deny"`
}

type Body struct {
	Allow []string `json:"allow"`
	Deny  []string `json:"deny"`
}
