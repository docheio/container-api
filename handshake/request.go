package handshake

type RequestStorage struct {
	Mount string
	Size  uint16
}

type RequestPort struct {
	Protocol string
	Number   uint16
}

type Request struct {
	Id      string
	Image   string
	Key     string
	Storage []RequestStorage
	Port    []RequestPort
}
