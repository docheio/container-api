package utils

import (
	"flag"
	"log"
)

type Flags struct {
	Namespace *string
	Address   *string
	Uniquekey *string
	Image     *string
	Debug     *bool
}

func (flags *Flags) Init() {
	flags.Namespace = flag.String("namespace", "", "Namespce of the container to operate.")
	flags.Address = flag.String("address", ":8080", "IP:Port to host api container.")
	flags.Uniquekey = flag.String("key", "", "Used to idnetify witch type of container.")
	flags.Image = flag.String("image", "", "Container image url.")
	flags.Debug = flag.Bool("debug", false, "Debug mode.")
	flag.Parse()
	if *flags.Namespace == "" {
		log.Fatal("Not specified Namespace")
	}
	if *flags.Uniquekey == "" {
		log.Fatal("Not specified Key")
	}
}
