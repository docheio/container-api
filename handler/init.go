/* ********************************************************************************************************** */
/*                                                                                                            */
/*                                                     :::::::::  ::::::::   ::::::::   :::    ::: :::::::::: */
/* init.go                                            :+:    :+: :+:    :+: :+:    :+: :+:    :+: :+:         */
/*                                                   +:+    +:+ +:+    +:+ +:+        +:+    +:+ +:+          */
/* By: yushsato <yukun@team.anylinks.jp>            +#+    +:+ +#+    +:+ +#+        +#++:++#++ +#++:++#      */
/*                                                 +#+    +#+ +#+    +#+ +#+        +#+    +#+ +#+            */
/* Created: 2023/05/27 04:25:07 by yushsato       #+#    #+# #+#    #+# #+#    #+# #+#    #+# #+#             */
/* Updated: 2023/05/27 04:25:08 by yushsato      #########  ########   ########  ###    ### ##########.io.    */
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
