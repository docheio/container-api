package handler

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Config struct {
	clientSet kubernetes.Clientset

	Uniquekey *string
	Namespace *string
	Image     *string
}

func (config *Config) Init() {
	if cliconf, err := rest.InClusterConfig(); err != nil {
		log.Fatal("Failed to load config file: %w", err)
	} else {
		if cliset, err := kubernetes.NewForConfig(cliconf); err != nil {
			log.Fatal("Failed to load config file: %w", err)
		} else {
			config.clientSet = *cliset
			return
		}
	}
}
