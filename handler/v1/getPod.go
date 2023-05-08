package v1

import (
	"github.com/docheio/container-api/handshake"
	corev1 "k8s.io/api/core/v1"
)

func GetPod(pods *corev1.PodList, response *handshake.Response, volumesList *[][]handshake.Volume) {
	for _, pod := range pods.Items {
		if label := pod.Labels["app"]; label != "" {
			volumes := []handshake.Volume{}
			*response = append(*response, handshake.IResponse{
				Id:     label,
				Name:   pod.Name,
				Status: string(pod.Status.Phase),
				Ports:  []handshake.Port{},
				Pvcs:   []handshake.Pvc{},
			})
			for _, volume := range pod.Spec.Volumes {
				for _, container := range pod.Spec.Containers {
					for _, volumeMount := range container.VolumeMounts {
						if volumeMount.Name == volume.Name && volume.PersistentVolumeClaim != nil {
							volumes = append(volumes, handshake.Volume{
								Name:  volume.Name,
								Mount: volumeMount.MountPath,
								Claim: volume.PersistentVolumeClaim.ClaimName,
							})
						}
					}
				}
			}
			*volumesList = append(*volumesList, volumes)
		}
	}
}
