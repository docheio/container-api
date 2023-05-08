package v1

import (
	"regexp"
	"strconv"

	"github.com/docheio/container-api/handshake"
	corev1 "k8s.io/api/core/v1"
)

func GetPvc(pvcs *corev1.PersistentVolumeClaimList, response *handshake.Response, volumesList *[][]handshake.Volume) {
	for _, pvc := range pvcs.Items {
		if label := pvc.Labels["app"]; label != "" {
			for num, volumes := range *volumesList {
				for _, volume := range volumes {
					reg := regexp.MustCompile("[1-9][0-9]*")
					match := reg.FindString(pvc.Status.Capacity.Storage().String())
					if ui64, err := strconv.ParseUint(match, 10, 64); err == nil {
						(*response)[num].Pvcs = append((*response)[num].Pvcs, handshake.Pvc{
							Id:    pvc.Name,
							Mount: volume.Mount,
							Size:  uint16(ui64),
						})
					}
				}
			}
		}
	}
}
