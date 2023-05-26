package handler

import "k8s.io/client-go/kubernetes"

type Handler struct {
	clientSet              kubernetes.Clientset
	originLabelSelectorKey string

	Namespace string
	Image     string
}

type PersistentVolumeClaimLink struct {
	claimName string
	mountPath string
}