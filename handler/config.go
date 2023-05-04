package handler

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Config struct {
	clientSet kubernetes.Clientset

	Uniquekey *string
	Namepsace *string
	Image     *string
}

func (config *Config) Init() {
	clconf, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	clset, err := kubernetes.NewForConfig(clconf)
	if err != nil {
		log.Fatal(err.Error())
	}
	config.clientSet = *clset
}
