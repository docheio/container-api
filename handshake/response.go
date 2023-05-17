package handshake

type ResponseStorage struct {
	Mount string
	Size  uint16
}

type ResponsePort struct {
	Protocol string
	Number   uint16
	Extrenal uint16
}

type Response struct {
	Id      string
	Image   string
	Key     string
	Storage []RequestStorage
	Port    []RequestPort
}
