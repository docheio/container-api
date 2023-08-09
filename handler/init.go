/* ********************************************************************************************************** */
/*                                                                                                            */
/*                                                     :::::::::  ::::::::   ::::::::   :::    ::: :::::::::: */
/* init.go                                            :+:    :+: :+:    :+: :+:    :+: :+:    :+: :+:         */
/*                                                   +:+    +:+ +:+    +:+ +:+        +:+    +:+ +:+          */
/* By: ES-Yukun <yukun@doche.io>                    +#+    +:+ +#+    +:+ +#+        +#++:++#++ +#++:++#      */
/*                                                 +#+    +#+ +#+    +#+ +#+        +#+    +#+ +#+            */
/* Created: 2023/08/09 23:46:45 by ES-Yukun       #+#    #+# #+#    #+# #+#    #+# #+#    #+# #+#             */
/* Updated: 2023/08/09 23:46:46 by ES-Yukun      #########  ########   ########  ###    ### ##########.io.    */
/*                                                                                                            */
/* ********************************************************************************************************** */

package handler

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func (handler *Handler) Init() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic("handler/init.go:11\n" + err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic("handler/init.go:15\n" + err.Error())
	}
	handler.clientSet = *clientSet
}
