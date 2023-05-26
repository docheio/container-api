/* ********************************************************************************************************** */
/*                                                                                                            */
/*                                                     :::::::::  ::::::::   ::::::::   :::    ::: :::::::::: */
/* create.go                                          :+:    :+: :+:    :+: :+:    :+: :+:    :+: :+:         */
/*                                                   +:+    +:+ +:+    +:+ +:+        +:+    +:+ +:+          */
/* By: yushsato <yukun@team.anylinks.jp>            +#+    +:+ +#+    +:+ +#+        +#++:++#++ +#++:++#      */
/*                                                 +#+    +#+ +#+    +#+ +#+        +#+    +#+ +#+            */
/* Created: 2023/05/27 04:22:51 by yushsato       #+#    #+# #+#    #+# #+#    #+# #+#    #+# #+#             */
/* Updated: 2023/05/27 04:22:52 by yushsato      #########  ########   ########  ###    ### ##########.io.    */
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
	"strings"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/docheio/container-api/utils"
)

// This function handles create vpc request

func (handler *Handler) Create(gc *gin.Context) {
	gc.Writer.Header().Set("Content-Type", "application/json;charset=UTF-8")

	response := Response{}
	request := Request{}

	body := make([]byte, gc.Request.ContentLength)
	gc.Request.Body.Read(body)
	if err := json.Unmarshal(body, &request); err != nil || request.Cpu == 0 || request.Mem == 0 {
		gc.Writer.WriteHeader(501)
		json.NewEncoder(gc.Writer).Encode(response)
		return
	}

	pvcClient := handler.clientSet.CoreV1().PersistentVolumeClaims(handler.Namespace)
	svcClient := handler.clientSet.CoreV1().Services(handler.Namespace)
	deployClient := handler.clientSet.AppsV1().Deployments(handler.Namespace)

	var svcs []*apiv1.Service
	var pvcs []*apiv1.PersistentVolumeClaim
	var deploys []*appsv1.Deployment

	id := utils.RFC1123()
	cpu := strconv.FormatUint(uint64(request.Cpu), 10) + "m"
	mem := strconv.FormatUint(uint64(request.Mem), 10) + "M"
	ctrPort := []apiv1.ContainerPort{}
	svcPort := []apiv1.ServicePort{}
	volumes := []apiv1.Volume{}
	volumeMounts := []apiv1.VolumeMount{}

	{ // Store object Services
		for _, port := range request.Ports {
			if !(port.Protocol == "TCP" || port.Protocol == "UDP") {
				gc.Writer.WriteHeader(501)
				json.NewEncoder(gc.Writer).Encode(response)
				return
			}
			svcPort = append(svcPort, apiv1.ServicePort{
				Name:     strings.ToLower(port.Protocol) + strconv.FormatUint(uint64(port.Internal), 10),
				Protocol: apiv1.Protocol(port.Protocol),
				TargetPort: intstr.IntOrString{
					IntVal: int32(port.Internal),
				},
				Port: int32(port.Internal),
			})
			ctrPort = append(ctrPort, apiv1.ContainerPort{
				Name:          strings.ToLower(port.Protocol) + strconv.FormatUint(uint64(port.Internal), 10),
				Protocol:      apiv1.Protocol(port.Protocol),
				ContainerPort: int32(port.Internal),
			})
		}
		svcs = append(svcs, &apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: id,
				Labels: map[string]string{
					"app":      id,
					"uniqekey": handler.Uniqekey,
				},
			},
			Spec: apiv1.ServiceSpec{
				Selector: map[string]string{
					"app": id,
				},
				Ports: svcPort,
				Type:  "LoadBalancer",
			},
		})
	}
	{ // Store object PersistentVOlumeClaims
		paths := []string{}
		for num, rpvc := range request.Pvcs {
			enable := true
			for _, path := range paths {
				if path == rpvc.Mount {
					enable = false
				}
			}
			if enable {
				paths = append(paths, rpvc.Mount)
				pvcName := id + "-" + utils.RFC1123()
				volumes = append(volumes, apiv1.Volume{
					Name: "vol-" + strconv.Itoa(num),
					VolumeSource: apiv1.VolumeSource{
						PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
							ClaimName: pvcName,
						},
					},
				})
				volumeMounts = append(volumeMounts, apiv1.VolumeMount{
					Name:      "vol-" + strconv.Itoa(num),
					MountPath: rpvc.Mount,
				})
				pvcs = append(pvcs, &apiv1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{
						Name: pvcName,
						Labels: map[string]string{
							"app":      id,
							"uniqekey": handler.Uniqekey,
						},
					},
					Spec: apiv1.PersistentVolumeClaimSpec{
						AccessModes: []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteMany},
						Resources: apiv1.ResourceRequirements{
							Requests: apiv1.ResourceList{
								"storage": resource.MustParse(strconv.FormatUint(uint64(rpvc.Size), 10) + "Gi"),
							},
						},
					},
				})
			}
		}
	}
	{ // Store object Deployments
		deploys = append(deploys, &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: id,
				Labels: map[string]string{
					"gen":      utils.RFC1123(),
					"app":      id,
					"uniqekey": handler.Uniqekey,
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: utils.Int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": id,
					},
				},
				Strategy: appsv1.DeploymentStrategy{
					Type: "Recreate",
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app":      id,
							"uniqekey": handler.Uniqekey,
						},
					},
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:  id,
								Image: handler.Image,
								Ports: ctrPort,
								Resources: apiv1.ResourceRequirements{
									Limits: apiv1.ResourceList{
										"cpu":    resource.MustParse(cpu),
										"memory": resource.MustParse(mem),
									},
								},
								VolumeMounts: volumeMounts,
							},
						},
						Volumes: volumes,
					},
				},
			},
		})
	}

	{ // create
		errorOccured := false
		for _, pvc := range pvcs {
			if _, err := pvcClient.Create(context.TODO(), pvc, metav1.CreateOptions{}); err != nil || errorOccured {
				log.Printf("[Create Error] pvc     id:%v (message: %v)\n", id, err)
				errorOccured = true
				break
			}
		}
		for _, deploy := range deploys {
			if _, err := deployClient.Create(context.TODO(), deploy, metav1.CreateOptions{}); err != nil || errorOccured {
				log.Printf("[Create Error] deploy  id:%v (message: %v)\n", id, err)
				errorOccured = true
				break
			}
		}
		for _, svc := range svcs {
			if _, err := svcClient.Create(context.TODO(), svc, metav1.CreateOptions{}); err != nil || errorOccured {
				log.Printf("[Create Error] service id:%v (message: %v)\n", id, err)
				errorOccured = true
				break
			}
		}
		if errorOccured {
			for _, pvc := range pvcs {
				if err := pvcClient.Delete(context.TODO(), pvc.Name, metav1.DeleteOptions{}); err != nil {
					log.Printf("[Delete Error] pvc     id:%v (message: %v)\n", id, err)
				}
			}
			for _, deploy := range deploys {
				if err := deployClient.Delete(context.TODO(), deploy.Name, metav1.DeleteOptions{}); err != nil {
					log.Printf("[Delete Error] deploy  id:%v (message: %v)\n", id, err)
				}
			}
			for _, svc := range svcs {
				if err := svcClient.Delete(context.TODO(), svc.Name, metav1.DeleteOptions{}); err != nil {
					log.Printf("[Delete Error] service id:%v (message: %v)\n", id, err)
				}
			}

			gc.Writer.WriteHeader(501)
			json.NewEncoder(gc.Writer).Encode(response)
			return
		}
	}
	{ // response
		option := metav1.ListOptions{
			LabelSelector: "uniqekey=" + handler.Uniqekey + ",app=" + id,
		}
		volumeLinks := VolumeLinks{} //　pvc connection info to pod

		// get pods
		pods, err := handler.clientSet.CoreV1().Pods(handler.Namespace).List(context.TODO(), option)
		if err != nil {
			gc.Writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(gc.Writer).Encode(response)
			return
		}
		// get svc
		svcs, err := handler.clientSet.CoreV1().Services(handler.Namespace).List(context.TODO(), option)
		if err != nil {
			gc.Writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(gc.Writer).Encode(response)
			return
		}
		// get pvc
		pvcs, err := handler.clientSet.CoreV1().PersistentVolumeClaims(handler.Namespace).List(context.TODO(), option)
		if err != nil {
			gc.Writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(gc.Writer).Encode(response)
			return
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
