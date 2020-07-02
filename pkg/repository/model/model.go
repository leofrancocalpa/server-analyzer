package model

type (
	// Server defines the structure data of a server
	Server struct {
		Address  string `json:"address"`
		SSLGrade string `json:"ssl_grade"`
		Country  string `json:"country"`
		Owner    string `json:"owner"`
	}
	//Analysis define the structure of the data of analyze?host=<hostname> request
	Analysis struct {
		Servers          []Server `json:"servers"`
		ServerChanged    bool     `json:"servers_changed"`
		SSLGradre        string   `json:"ssl_grade"`
		PreviousSSLGrade string   `json:"previous_ssl_grade"`
		Logo             string   `json:"logo"`
		Title            string   `json:"title"`
		IsDown           bool     `json:"is_down"`
	}

	//ServerInfo Struct to fetch from DB
	ServerInfo struct {
		DNS  string `json:"dns"`
		Data []byte `json:"data"`
	}

	//ItemsInfo struct that match with DB query to fetch
	ItemsInfo struct {
		Items []ServerInfo `json:"items"`
	}
	//ServersQueried info
	ServersQueried struct {
		DNS  string   `json:"dns"`
		Data Analysis `json:"data"`
	}

	//ItemsToResponse Structure used to wrap the JSON of the Endpoint that returns the servers that have already been consulted before
	ItemsToResponse struct {
		Items []ServersQueried `json:"items"`
	}
)
