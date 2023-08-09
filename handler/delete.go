/* ********************************************************************************************************** */
/*                                                                                                            */
/*                                                     :::::::::  ::::::::   ::::::::   :::    ::: :::::::::: */
/* delete.go                                          :+:    :+: :+:    :+: :+:    :+: :+:    :+: :+:         */
/*                                                   +:+    +:+ +:+    +:+ +:+        +:+    +:+ +:+          */
/* By: ES-Yukun <yukun@doche.io>                    +#+    +:+ +#+    +:+ +#+        +#++:++#++ +#++:++#      */
/*                                                 +#+    +#+ +#+    +#+ +#+        +#+    +#+ +#+            */
/* Created: 2023/08/09 23:46:07 by ES-Yukun       #+#    #+# #+#    #+# #+#    #+# #+#    #+# #+#             */
/* Updated: 2023/08/09 23:46:08 by ES-Yukun      #########  ########   ########  ###    ### ##########.io.    */
/*                                                                                                            */
/* ********************************************************************************************************** */

// This package exports HTTP request handler used by Gin.

package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) Delete(gc *gin.Context) {
	gc.Writer.Header().Set("Content-Type", "application/json;charset=UTF-8")

	response := Response{}

	deployClient := handler.clientSet.AppsV1().Deployments(handler.Namespace)
	svcClient := handler.clientSet.CoreV1().Services(handler.Namespace)
	pvcClient := handler.clientSet.CoreV1().PersistentVolumeClaims(handler.Namespace)

	if id := gc.Param("id"); id == "" {
		gc.Writer.WriteHeader(501)
		json.NewEncoder(gc.Writer).Encode(response)
		return
	} else {
		option := metav1.ListOptions{
			LabelSelector: "uniquekey=" + handler.Uniquekey + ",app=" + id,
		}

		// get pods
		deploys, err := handler.clientSet.AppsV1().Deployments(handler.Namespace).List(context.TODO(), option)
		if err != nil {
			gc.Writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(gc.Writer).Encode(response)
		}
		// get svc
		svcs, err := handler.clientSet.CoreV1().Services(handler.Namespace).List(context.TODO(), option)
		if err != nil {
			gc.Writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(gc.Writer).Encode(response)
		}
		// get pvc
		pvcs, err := handler.clientSet.CoreV1().PersistentVolumeClaims(handler.Namespace).List(context.TODO(), option)
		if err != nil {
			gc.Writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(gc.Writer).Encode(response)
		}

		for _, deploy := range deploys.Items {
			if err := deployClient.Delete(context.TODO(), deploy.Name, metav1.DeleteOptions{}); err != nil {
				log.Printf("[Delete Error] pvc     id:%v (message: %v)\n", id, err)
			}
		}
		for _, svc := range svcs.Items {
			if err := svcClient.Delete(context.TODO(), svc.Name, metav1.DeleteOptions{}); err != nil {
				log.Printf("[Delete Error] deploy  id:%v (message: %v)\n", id, err)
			}
		}
		for _, pvc := range pvcs.Items {
			if err := pvcClient.Delete(context.TODO(), pvc.Name, metav1.DeleteOptions{}); err != nil {
				log.Printf("[Delete Error] service id:%v (message: %v)\n", id, err)
			}
		}

		{ // response
			option := metav1.ListOptions{
				LabelSelector: "uniquekey=" + handler.Uniquekey + ",app=" + id,
			}
			volumeLinks := VolumeLinks{} //　pvc connection info to pod

			// get pods
			pods, err := handler.clientSet.CoreV1().Pods(handler.Namespace).List(context.TODO(), option)
			if err != nil {
				gc.Writer.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(gc.Writer).Encode(response)
			}
			// get svc
			svcs, err := handler.clientSet.CoreV1().Services(handler.Namespace).List(context.TODO(), option)
			if err != nil {
				gc.Writer.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(gc.Writer).Encode(response)
			}
			// get pvc
			pvcs, err := handler.clientSet.CoreV1().PersistentVolumeClaims(handler.Namespace).List(context.TODO(), option)
			if err != nil {
				gc.Writer.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(gc.Writer).Encode(response)
			}

			// Parse items and store it in "response"
			for _, pod := range pods.Items {
				if label := pod.Labels["app"]; label != "" {
					response = append(response, IResponse{
						Id:     label,
						Name:   pod.Name,
						Status: string(pod.Status.Phase),
						Ports:  []IPort{},
						Pvcs:   []IPvc{},
					})

					// List pvc connection information and store in "volumeLinks"
					volumeLinkGroup := VolumeLinkGroup{}
					for _, volume := range pod.Spec.Volumes {
						for _, container := range pod.Spec.Containers {
							for _, volumeMount := range container.VolumeMounts {
								if volumeMount.Name == volume.Name && volume.PersistentVolumeClaim != nil {
									volumeLinkGroup = append(volumeLinkGroup, VolumeLink{
										Name:  volume.Name,
										Mount: volumeMount.MountPath,
										Claim: volume.PersistentVolumeClaim.ClaimName,
									})
								}
							}
						}
					}
					volumeLinks = append(volumeLinks, volumeLinkGroup)
				}
			}
			for _, svc := range svcs.Items {
				if label := svc.Labels["app"]; label != "" {
					for num, responseSplit := range response {
						if responseSplit.Id == label {
							for _, port := range svc.Spec.Ports {
								response[num].Ports = append(response[num].Ports, IPort{
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

					// Find connection info matching pvc in "volumeLinks" and store it with size in "response"
					for num, volumeLinkGroup := range volumeLinks {
						for _, volumeLink := range volumeLinkGroup {
							reg := regexp.MustCompile("[1-9][0-9]*")
							match := reg.FindString(pvc.Status.Capacity.Storage().String())
							if ui64, err := strconv.ParseUint(match, 10, 64); err == nil {
								response[num].Pvcs = append(response[num].Pvcs, IPvc{
									Id:    pvc.Name,
									Mount: volumeLink.Mount,
									Size:  uint16(ui64),
								})
							}
						}
					}
				}
			}

			json.NewEncoder(gc.Writer).Encode(response)
		}
	}
}
