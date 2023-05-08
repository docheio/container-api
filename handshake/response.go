package handshake

type IResponse struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Ports  []Port
	Pvcs   []Pvc
}

type Response []IResponse
