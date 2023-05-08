package v1

import (
	"github.com/docheio/container-api/handshake"
	corev1 "k8s.io/api/core/v1"
)

func GetSvc(svcs *corev1.ServiceList, response *handshake.Response, volumesList *[][]handshake.Volume) {
	for _, svc := range svcs.Items {
		if label := svc.Labels["app"]; label != "" {
			for num, res_ := range *response {
				if res_.Id == label {
					for _, port := range svc.Spec.Ports {
						(*response)[num].Ports = append((*response)[num].Ports, handshake.Port{
							Protocol: string(port.Protocol),
							Internal: uint16(port.Port),
							External: uint16(port.NodePort),
						})
					}
				}
			}
		}
	}
}
