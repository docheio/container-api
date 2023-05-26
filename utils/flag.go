package utils

import (
	"flag"
)

type Flags struct {
	Namespace        *string
	Uniquekey        *string
	Address          *string
	LabelSelectorKey *string
	Debug            *bool
}

func (flags *Flags) Init() {
	flags.Namespace = flag.String("namespace", "eureka", "Specify the namespace in which to operate container.")
	flags.Address = flag.String("address", ":8080", "Specify IP:Port for hosting api")
	flags.Debug = flag.Bool("debug", false, "debug mode")
	flags.Uniquekey = flag.String("key", "", "Specify a unique key. This is used to identify to which group the container belongse")
	flag.Parse()
}
