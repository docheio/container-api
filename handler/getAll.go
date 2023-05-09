package handler

import (
	"context"
	"encoding/json"

	handlerv1 "github.com/docheio/container-api/handler/v1"
	"github.com/docheio/container-api/handshake"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (config *Config) GetAll(gc *gin.Context) {
	var res handshake.Response
	var volumesList [][]handshake.Volume
	var option metav1.ListOptions

	res = handshake.Response{}
	volumesList = [][]handshake.Volume{}
	option = metav1.ListOptions{}
	option.LabelSelector = "uniquekey=" + *config.Uniquekey

	pods, err := config.clientSet.CoreV1().Pods(*config.Namespace).List(context.TODO(), option)
	if e := config.checkInternalServerError(gc, err); e {
		return
	}
	svcs, err := config.clientSet.CoreV1().Services(*config.Namespace).List(context.TODO(), option)
	if e := config.checkInternalServerError(gc, err); e {
		return
	}
	pvcs, err := config.clientSet.CoreV1().PersistentVolumeClaims(*config.Namespace).List(context.TODO(), option)
	if e := config.checkInternalServerError(gc, err); e {
		return
	}

	handlerv1.FetchPod(pods, &res, &volumesList)
	handlerv1.FetchSvc(svcs, &res, &volumesList)
	handlerv1.FetchPvc(pvcs, &res, &volumesList)

	gc.Writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(gc.Writer).Encode(res)
}
