package handshake

type Request struct {
	Cpu   uint16 `json:"cpu"`
	Mem   uint16 `json:"mem"`
	Ports []Port
	Pvcs  []Pvc
}
