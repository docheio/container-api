/* ********************************************************************************************************** */
/*                                                                                                            */
/*                                                     :::::::::  ::::::::   ::::::::   :::    ::: :::::::::: */
/* handshake.go                                       :+:    :+: :+:    :+: :+:    :+: :+:    :+: :+:         */
/*                                                   +:+    +:+ +:+    +:+ +:+        +:+    +:+ +:+          */
/* By: yushsato <yukun@team.anylinks.jp>            +#+    +:+ +#+    +:+ +#+        +#++:++#++ +#++:++#      */
/*                                                 +#+    +#+ +#+    +#+ +#+        +#+    +#+ +#+            */
/* Created: 2023/05/27 04:24:56 by yushsato       #+#    #+# #+#    #+# #+#    #+# #+#    #+# #+#             */
/* Updated: 2023/05/27 04:24:56 by yushsato      #########  ########   ########  ###    ### ##########.io.    */
/*                                                                                                            */
/* ********************************************************************************************************** */

package handler

type IPort struct {
	Protocol string `json:"protocol"`
	Internal uint16 `json:"internal"`
	External uint16 `json:"external"`
}
type IPvc struct {
	Id    string `json:"id"`
	Mount string `json:"mount"`
	Size  uint16 `json:"size"`
}
type IResponse struct {
	Id     string  `json:"id"`
	Name   string  `json:"name"`
	Status string  `json:"status"`
	Ports  []IPort `json:"ports"`
	Pvcs   []IPvc  `json:"pvcs"`
}
type Response []IResponse

type OPort struct {
	Protocol string `json:"protocol"`
	Internal uint16 `json:"internal"`
}
type Request struct {
	Cpu   uint16  `json:"cpu"`
	Mem   uint16  `json:"mem"`
	Ports []OPort `json:"ports"`
	Pvcs  []IPvc  `json:"pvcs"`
}
