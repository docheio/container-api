package v1

import (
	"context"
	"errors"
	"log"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateSet(clientSet kubernetes.Clientset, namespace string, pvcs []*corev1.PersistentVolumeClaim, svcs []*corev1.Service, deploy *appsv1.Deployment) error {

	var errorOccured bool
	var errs string

	errorOccured = false
	errs = ""
	{
		for _, pvc := range pvcs {
			if _, err := clientSet.CoreV1().PersistentVolumeClaims(namespace).Create(context.TODO(), pvc, metav1.CreateOptions{}); err != nil || errorOccured {
				errs = errs + err.Error() + "\n"
				errorOccured = true
			}
		}
		for _, svc := range svcs {
			if _, err := clientSet.CoreV1().Services(namespace).Create(context.TODO(), svc, metav1.CreateOptions{}); err != nil || errorOccured {
				errs = errs + err.Error() + "\n"
				errorOccured = true
			}
		}
		if _, err := clientSet.AppsV1().Deployments(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{}); err != nil || errorOccured {
			errs = errs + err.Error() + "\n"
			errorOccured = true
		}
	}
	if errorOccured {
		for _, pvc := range pvcs {
			if err := clientSet.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), pvc.Name, metav1.DeleteOptions{}); err != nil {
				errs = errs + err.Error() + "\n"
			}
		}
		for _, svc := range svcs {
			if err := clientSet.CoreV1().Services(namespace).Delete(context.TODO(), svc.Name, metav1.DeleteOptions{}); err != nil {
				errs = errs + err.Error() + "\n"
			}
		}
		if err := clientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), deploy.Name, metav1.DeleteOptions{}); err != nil {
			errs = errs + err.Error() + "\n"
		}
	}
	if errs == "" {
		return nil
	} else {
		for _, str := range strings.Split(errs, "\n") {
			if str != "" {
				log.Println(str)
			}
		}
		return errors.New(errs)
	}
}
