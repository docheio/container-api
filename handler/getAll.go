package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/docheio/container-api/handshake"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (config *Config) GetAll(gc *gin.Context) {
	var res handshake.Response
	var volumeLink [][]handshake.Volume
	var option metav1.ListOptions

	gc.Writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	res = handshake.Response{}
	volumeLink = [][]handshake.Volume{}
	option = metav1.ListOptions{
		LabelSelector: "uniquekey=" + *config.Uniquekey,
	}

	pods, err := config.clientSet.CoreV1().Pods(*config.Namepsace).List(context.TODO(), option)
	if err != nil {
		gc.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(gc.Writer).Encode(res)
	}
	svcs, err := config.clientSet.CoreV1().Services(*config.Namepsace).List(context.TODO(), option)
	if err != nil {
		if err != nil {
			gc.Writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(gc.Writer).Encode(res)
		}
	}
	pvcs, err := config.clientSet.CoreV1().PersistentVolumeClaims(*config.Namepsace).List(context.TODO(), option)
	if err != nil {
		if err != nil {
			gc.Writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(gc.Writer).Encode(res)
		}
	}
	for _, pod := range pods.Items {
		if label := pod.Labels["app"]; label != "" {
			volumes := []handshake.Volume{}
			res = append(res, handshake.IResponse{
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
			volumeLink = append(volumeLink, volumes)
		}
	}
	for _, svc := range svcs.Items {
		if label := svc.Labels["app"]; label != "" {
			for num, res_ := range res {
				if res_.Id == label {
					for _, port := range svc.Spec.Ports {
						res[num].Ports = append(res[num].Ports, handshake.Port{
							Protocol: string(port.Protocol),
							Internal: uint16(port.Port),
							External: uint16(port.NodePort),
						})
					}
				}
			}
		}
	}
	for _, pvc := range pvcs.Items {
		if label := pvc.Labels["app"]; label != "" {
			for num, volumes := range volumeLink {
				for _, volume := range volumes {
					reg := regexp.MustCompile("[1-9][0-9]*")
					match := reg.FindString(pvc.Status.Capacity.Storage().String())
					if ui64, err := strconv.ParseUint(match, 10, 64); err == nil {
						res[num].Pvcs = append(res[num].Pvcs, handshake.Pvc{
							Id:    pvc.Name,
							Mount: volume.Mount,
							Size:  uint16(ui64),
						})
					}
				}
			}
		}
	}
	json.NewEncoder(gc.Writer).Encode(res)
}
