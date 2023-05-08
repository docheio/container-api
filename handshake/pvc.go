package handshake

type Pvc struct {
	Id    string `json:"id"`
	Mount string `json:"mount"`
	Size  uint16 `json:"size"`
}
