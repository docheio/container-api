// This package exports HTTP request handler used by Gin.
package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	handlerv1 "github.com/docheio/container-api/handler/v1"
	"github.com/docheio/container-api/handshake"
	"github.com/docheio/container-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func jsonGinUnmarshal(gc *gin.Context, v any) error {
	var requestBody []byte

	requestBody = make([]byte, gc.Request.ContentLength)
	if err := json.Unmarshal(requestBody, v); err != nil {
		return err
	}
	return nil
}

func errorBadRequest(gc *gin.Context, err error) bool {
	if err != nil {
		gc.Writer.WriteHeader(http.StatusBadRequest)
		gc.Writer.WriteString(err.Error())
		return true
	}
	return false
}

func servicePorts(ports []handshake.Port) []corev1.ServicePort {
	ret := []corev1.ServicePort{}

	for _, port := range ports {
		ret = append(ret, corev1.ServicePort{
			Name:     strings.ToLower(port.Protocol) + strconv.FormatUint(uint64(port.Internal), 10),
			Protocol: corev1.Protocol(port.Protocol),
			TargetPort: intstr.IntOrString{
				IntVal: int32(port.Internal),
			},
			Port: int32(port.Internal),
		})
	}
	return ret
}

func containerPorts(ports []handshake.Port) []corev1.ContainerPort {
	ret := []corev1.ContainerPort{}

	for _, port := range ports {
		ret = append(ret, corev1.ContainerPort{
			Name:          strings.ToLower(port.Protocol) + strconv.FormatUint(uint64(port.Internal), 10),
			Protocol:      corev1.Protocol(port.Protocol),
			ContainerPort: int32(port.Internal),
		})
	}
	return ret
}

func (config *Config) Create(gc *gin.Context) {
	var request handshake.Request
	var id string
	var svcs []*corev1.Service
	var pvcs []*corev1.PersistentVolumeClaim
	var dply *appsv1.Deployment

	if err := jsonGinUnmarshal(gc, &request); err != nil || request.Cpu == 0 || request.Mem == 0 {
		errorBadRequest(gc, errors.New("Invaild Value"))
		return
	}
	id = utils.RFC1123Gen(12)

	svcs = append(svcs, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
			Labels: map[string]string{
				"app":       id,
				"uniquekey": *config.Uniquekey,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": id,
			},
			Type:  "LoadBalancer",
		},
	})
}

// func (config *Config) Create(gc *gin.Context) {
// 	var request handshake.Request
// 	var requestBody []byte
// 	var svcs []*corev1.Service
// 	var pvcs []*corev1.PersistentVolumeClaim
// 	var deploy *appsv1.Deployment
// 	var id string
// 	var cpu string
// 	var mem string
// 	var ctrPort []corev1.ContainerPort
// 	var svcPort []corev1.ServicePort
// 	var volumes []corev1.Volume
// 	var volumeMounts []corev1.VolumeMount

// 	request = handshake.Request{}
// 	svcPort = []corev1.ServicePort{}
// 	ctrPort = []corev1.ContainerPort{}
// 	volumes = []corev1.Volume{}
// 	volumeMounts = []corev1.VolumeMount{}
// 	requestBody = make([]byte, gc.Request.ContentLength)
// 	gc.Request.Body.Read(requestBody)
// 	if err := json.Unmarshal(requestBody, &request); err != nil || request.Cpu == 0 || request.Mem == 0 {
// 		if request.Cpu == 0 || request.Mem == 0 {
// 			config.checkInternalServerError(gc, errors.New("Not allowed value as cpu or memory size"))
// 		} else {
// 			config.checkInternalServerError(gc, err)
// 		}
// 		return
// 	}
// 	id = utils.RFC1123Gen(12)
// 	cpu = strconv.FormatUint(uint64(request.Cpu), 10) + "m"
// 	mem = strconv.FormatUint(uint64(request.Mem), 10) + "M"

// 	{
// 		svcs = append(svcs, &corev1.Service{
// 			ObjectMeta: metav1.ObjectMeta{
// 				Name: id,
// 				Labels: map[string]string{
// 					"app":       id,
// 					"uniquekey": *config.Uniquekey,
// 				},
// 			},
// 			Spec: corev1.ServiceSpec{
// 				Selector: map[string]string{
// 					"app": id,
// 				},
// 				Ports: svcPort,
// 				Type:  "LoadBalancer",
// 			},
// 		})
// 	}
// 	{
// 		paths := []string{}
// 		for num, rpvc := range request.Pvcs {
// 			enable := true
// 			for _, path := range paths {
// 				if path == rpvc.Mount {
// 					enable = false
// 				}
// 			}
// 			if enable {
// 				paths = append(paths, rpvc.Mount)
// 				pvcName := id + "-" + utils.RFC1123Gen(12)
// 				volumes = append(volumes, corev1.Volume{
// 					Name: "vol-" + strconv.Itoa(num),
// 					VolumeSource: corev1.VolumeSource{
// 						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
// 							ClaimName: pvcName,
// 						},
// 					},
// 				})
// 				volumeMounts = append(volumeMounts, corev1.VolumeMount{
// 					Name:      "vol-" + strconv.Itoa(num),
// 					MountPath: rpvc.Mount,
// 				})
// 				pvcs = append(pvcs, &corev1.PersistentVolumeClaim{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name: pvcName,
// 						Labels: map[string]string{
// 							"app":       id,
// 							"uniquekey": *config.Uniquekey,
// 						},
// 					},
// 					Spec: corev1.PersistentVolumeClaimSpec{
// 						AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany},
// 						Resources: corev1.ResourceRequirements{
// 							Requests: corev1.ResourceList{
// 								"storage": resource.MustParse(strconv.FormatUint(uint64(rpvc.Size), 10) + "Gi"),
// 							},
// 						},
// 					},
// 				})
// 			}
// 		}
// 	}
// 	{
// 		deploy = &appsv1.Deployment{
// 			ObjectMeta: metav1.ObjectMeta{
// 				Name: id,
// 				Labels: map[string]string{
// 					"app":       id,
// 					"uniquekey": *config.Uniquekey,
// 				},
// 			},
// 			Spec: appsv1.DeploymentSpec{
// 				Replicas: utils.Int32Ptr(1),
// 				Selector: &metav1.LabelSelector{
// 					MatchLabels: map[string]string{
// 						"app": id,
// 					},
// 				},
// 				Strategy: appsv1.DeploymentStrategy{
// 					Type: "Recreate",
// 				},
// 				Template: corev1.PodTemplateSpec{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Labels: map[string]string{
// 							"app":       id,
// 							"uniquekey": *config.Uniquekey,
// 						},
// 					},
// 					Spec: corev1.PodSpec{
// 						Containers: []corev1.Container{
// 							{
// 								Name:  id,
// 								Image: *config.Image,
// 								Ports: ctrPort,
// 								Resources: corev1.ResourceRequirements{
// 									Limits: corev1.ResourceList{
// 										"cpu":    resource.MustParse(cpu),
// 										"memory": resource.MustParse(mem),
// 									},
// 								},
// 								VolumeMounts: volumeMounts,
// 							},
// 						},
// 						Volumes: volumes,
// 					},
// 				},
// 			},
// 		}
// 	}
// 	{
// 		var errorOccured bool
// 		var errs string

// 		errorOccured = false
// 		errs = ""
// 		{
// 			for _, pvc := range pvcs {
// 				if _, err := config.clientSet.CoreV1().PersistentVolumeClaims(*config.Namespace).Create(context.TODO(), pvc, metav1.CreateOptions{}); err != nil || errorOccured {
// 					errs = errs + err.Error() + "\n"
// 					errorOccured = true
// 				}
// 			}
// 			for _, svc := range svcs {
// 				if _, err := config.clientSet.CoreV1().Services(*config.Namespace).Create(context.TODO(), svc, metav1.CreateOptions{}); err != nil || errorOccured {
// 					errs = errs + err.Error() + "\n"
// 					errorOccured = true
// 				}
// 			}
// 			if _, err := config.clientSet.AppsV1().Deployments(*config.Namespace).Create(context.TODO(), deploy, metav1.CreateOptions{}); err != nil || errorOccured {
// 				errs = errs + err.Error() + "\n"
// 				errorOccured = true
// 			}
// 		}
// 		if errorOccured {
// 			for _, pvc := range pvcs {
// 				if err := config.clientSet.CoreV1().PersistentVolumeClaims(*config.Namespace).Delete(context.TODO(), pvc.Name, metav1.DeleteOptions{}); err != nil {
// 					errs = errs + err.Error() + "\n"
// 				}
// 			}
// 			for _, svc := range svcs {
// 				if err := config.clientSet.CoreV1().Services(*config.Namespace).Delete(context.TODO(), svc.Name, metav1.DeleteOptions{}); err != nil {
// 					errs = errs + err.Error() + "\n"
// 				}
// 			}
// 			if err := config.clientSet.AppsV1().Deployments(*config.Namespace).Delete(context.TODO(), deploy.Name, metav1.DeleteOptions{}); err != nil {
// 				errs = errs + err.Error() + "\n"
// 			}
// 		}
// 		if errs != "" {
// 			for _, str := range strings.Split(errs, "\n") {
// 				if str != "" {
// 					log.Println(str)
// 				}
// 			}
// 		}
// 	}
// 	{
// 		var res handshake.Response
// 		var volumesList [][]handshake.Volume
// 		var option metav1.ListOptions

// 		res = handshake.Response{}
// 		volumesList = [][]handshake.Volume{}
// 		option = metav1.ListOptions{}
// 		option.LabelSelector = "uniquekey=" + *config.Uniquekey + ",app=" + id

// 		pods, err := config.clientSet.CoreV1().Pods(*config.Namespace).List(context.TODO(), option)
// 		config.checkInternalServerError(gc, err)
// 		svcs, err := config.clientSet.CoreV1().Services(*config.Namespace).List(context.TODO(), option)
// 		config.checkInternalServerError(gc, err)
// 		pvcs, err := config.clientSet.CoreV1().PersistentVolumeClaims(*config.Namespace).List(context.TODO(), option)
// 		config.checkInternalServerError(gc, err)

// 		handlerv1.FetchPod(pods, &res, &volumesList)
// 		handlerv1.FetchSvc(svcs, &res, &volumesList)
// 		handlerv1.FetchPvc(pvcs, &res, &volumesList)

// 		gc.Writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
// 		json.NewEncoder(gc.Writer).Encode(res)
// 	}
// }
