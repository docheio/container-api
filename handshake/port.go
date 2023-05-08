package handshake

type Port struct {
	Protocol string `json:"protocol"`
	Internal uint16 `json:"internal"`
	External uint16 `json:"external"`
}
