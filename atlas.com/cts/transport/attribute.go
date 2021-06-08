package transport

type dataListContainer struct {
	Data []dataBody `json:"data"`
}

type dataBody struct {
	Id         string    `json:"id"`
	Type       string    `json:"type"`
	Attributes attributes `json:"attributes"`
}

type attributes struct {
	Enabled     bool   `json:"enabled"`
	Source      uint32 `json:"source"`
	Departure   uint32 `json:"departure"`
	Transport   uint32 `json:"transport"`
	Arrival     uint32 `json:"arrival"`
	Destination uint32 `json:"destination"`
}
