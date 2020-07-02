package client

type (
	//structs used to wrap the response from API_SSLLAB and WHOIS comand

	// Enpoint defines an struct without "country" and "owner" fields
	Enpoint struct {
		IPAddress string `json:"ipAddress"`
		Grade     string `json:"grade"`
		Country   string `json:"country"`
		Owner     string `json:"owner"`
	}
	//ResultServerInfo defines the queried result
	ResultServerInfo struct {
		Host     string    `json:"host"`
		Status   string    `json:"status"`
		Enpoints []Enpoint `json:"endpoints"`
		SSLGrade string    `json:"ssl_grade"`
		Logo     string    `json:"logo"`
		Title    string    `json:"title"`
		IsDown   bool      `json:"is_down"`
	}
)
