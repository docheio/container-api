package utils

import (
	"flag"
	"log"
)

type Flag struct {
	Namespace *string
	Address   *string
	Uniquekey *string
	Image     *string
	Debug     *bool
}

func (f *Flag) Init() {
	*f = Flag{}
	f.Namespace = flag.String("namespace", "", "Namespace of the container to operate.")
	f.Address = flag.String("address", ":8080", "IP:Port to host api container.")
	f.Uniquekey = flag.String("key", "", "Used to idenitify witch type of container.")
	f.Image = flag.String("image", "", "Container image uri.")
	f.Debug = flag.Bool("debug", false, "Debug mode.")

	if *f.Namespace == "" {
		log.Fatal("Namespace must be specified.")
	}
	if *f.Uniquekey == "" {
		log.Fatal("Uniquekey must be specified.")
	}
	if *f.Image == "" {
		log.Fatal("Image URI must be specified.")
	}
}
