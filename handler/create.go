// This package exports HTTP request handler used by Gin.
package handler

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	handlerv1 "github.com/docheio/container-api/handler/v1"
	"github.com/docheio/container-api/handshake"
	"github.com/docheio/container-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// This function handles create vpc request

func (config *Config) Create(gc *gin.Context) {
	var request handshake.Request
	var requestBody []byte
	var svcs []*apiv1.Service
	var pvcs []*apiv1.PersistentVolumeClaim
	var deploys []*appsv1.Deployment
	var id string
	var cpu string
	var mem string
	var ctrPort []apiv1.ContainerPort
	var svcPort []apiv1.ServicePort
	var volumes []apiv1.Volume
	var volumeMounts []apiv1.VolumeMount

	request = handshake.Request{}
	ctrPort = []apiv1.ContainerPort{}
	svcPort = []apiv1.ServicePort{}
	volumes = []apiv1.Volume{}
	volumeMounts = []apiv1.VolumeMount{}
	requestBody = make([]byte, gc.Request.ContentLength)
	gc.Request.Body.Read(requestBody)
	if err := json.Unmarshal(requestBody, &request); err != nil || request.Cpu == 0 || request.Mem == 0 {
		if request.Cpu == 0 || request.Mem == 0 {
			config.checkInternalServerError(gc, errors.New("Not allowed value as cpu or memory size"))
		} else {
			config.checkInternalServerError(gc, err)
		}
		return
	}
	id = utils.RFC1123Gen(12)
	cpu = strconv.FormatUint(uint64(request.Cpu), 10) + "m"
	mem = strconv.FormatUint(uint64(request.Mem), 10) + "M"

	{
		for _, port := range request.Ports {
			if !(port.Protocol == "TCP" || port.Protocol == "UDP") {
				config.checkInternalServerError(gc, errors.New("Not allowed value as protocol name"))
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
					"app":       id,
					"uniquekey": *config.Uniquekey,
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
				pvcName := id + "-" + utils.RFC1123Gen(12)
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
							"app":       id,
							"uniquekey": *config.Uniquekey,
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
					"app":       id,
					"uniquekey": *config.Uniquekey,
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
							"app":       id,
							"uniquekey": *config.Uniquekey,
						},
					},
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:  id,
								Image: *config.Image,
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
			if _, err := config.clientSet.CoreV1().PersistentVolumeClaims(*config.Namespace).Create(context.TODO(), pvc, metav1.CreateOptions{}); err != nil || errorOccured {
				log.Printf("[Create Error] pvc     id:%v (message: %v)\n", id, err)
				errorOccured = true
				break
			}
		}
		for _, deploy := range deploys {
			if _, err := config.clientSet.AppsV1().Deployments(*config.Namespace).Create(context.TODO(), deploy, metav1.CreateOptions{}); err != nil || errorOccured {
				log.Printf("[Create Error] deploy  id:%v (message: %v)\n", id, err)
				errorOccured = true
				break
			}
		}
		for _, svc := range svcs {
			if _, err := config.clientSet.CoreV1().Services(*config.Namespace).Create(context.TODO(), svc, metav1.CreateOptions{}); err != nil || errorOccured {
				log.Printf("[Create Error] service id:%v (message: %v)\n", id, err)
				errorOccured = true
				break
			}
		}
		if errorOccured {
			for _, pvc := range pvcs {
				if err := config.clientSet.CoreV1().PersistentVolumeClaims(*config.Namespace).Delete(context.TODO(), pvc.Name, metav1.DeleteOptions{}); err != nil {
					log.Printf("[Delete Error] pvc     id:%v (message: %v)\n", id, err)
				}
			}
			for _, deploy := range deploys {
				if err := config.clientSet.AppsV1().Deployments(*config.Namespace).Delete(context.TODO(), deploy.Name, metav1.DeleteOptions{}); err != nil {
					log.Printf("[Delete Error] deploy  id:%v (message: %v)\n", id, err)
				}
			}
			for _, svc := range svcs {
				if err := config.clientSet.CoreV1().Services(*config.Namespace).Delete(context.TODO(), svc.Name, metav1.DeleteOptions{}); err != nil {
					log.Printf("[Delete Error] service id:%v (message: %v)\n", id, err)
				}
			}
			return
		}
	}
	{
		var res handshake.Response
		var volumesList [][]handshake.Volume
		var option metav1.ListOptions

		res = handshake.Response{}
		volumesList = [][]handshake.Volume{}
		option = metav1.ListOptions{}
		option.LabelSelector = "uniquekey=" + *config.Uniquekey + ",app=" + id

		pods, err := config.clientSet.CoreV1().Pods(*config.Namespace).List(context.TODO(), option)
		config.checkInternalServerError(gc, err)
		svcs, err := config.clientSet.CoreV1().Services(*config.Namespace).List(context.TODO(), option)
		config.checkInternalServerError(gc, err)
		pvcs, err := config.clientSet.CoreV1().PersistentVolumeClaims(*config.Namespace).List(context.TODO(), option)
		config.checkInternalServerError(gc, err)

		handlerv1.FetchPod(pods, &res, &volumesList)
		handlerv1.FetchSvc(svcs, &res, &volumesList)
		handlerv1.FetchPvc(pvcs, &res, &volumesList)

		gc.Writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
		json.NewEncoder(gc.Writer).Encode(res)
	}
}
