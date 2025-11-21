package model

type RegistryRequest struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	AgentPort int    `json:"agent_port"`
	IsLocal   bool   `json:"is_local"`
	Status    string `json:"status"`
}

type RegistryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
