package handler

import (
	"context"
	"encoding/json"

	handlerv1 "github.com/docheio/container-api/handler/v1"
	"github.com/docheio/container-api/handshake"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (config *Config) GetOne(gc *gin.Context) {
	var res handshake.Response
	var volumesList [][]handshake.Volume
	var option metav1.ListOptions
	var id string

	if id = gc.Param("id"); id == "" {
		config.checkInternalServerError(gc, errors.New("Not specified Id"))
		return
	}

	res = handshake.Response{}
	volumesList = [][]handshake.Volume{}
	option = metav1.ListOptions{}
	option.LabelSelector = "uniquekey=" + *config.Uniquekey + ",app=" + id

	pods, err := config.clientSet.CoreV1().Pods(*config.Namepsace).List(context.TODO(), option)
	config.checkInternalServerError(gc, err)
	svcs, err := config.clientSet.CoreV1().Services(*config.Namepsace).List(context.TODO(), option)
	config.checkInternalServerError(gc, err)
	pvcs, err := config.clientSet.CoreV1().PersistentVolumeClaims(*config.Namepsace).List(context.TODO(), option)
	config.checkInternalServerError(gc, err)

	handlerv1.FetchPod(pods, &res, &volumesList)
	handlerv1.FetchSvc(svcs, &res, &volumesList)
	handlerv1.FetchPvc(pvcs, &res, &volumesList)

	gc.Writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(gc.Writer).Encode(res)
}
