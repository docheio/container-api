/* ********************************************************************************************************** */
/*                                                                                                            */
/*                                                     :::::::::  ::::::::   ::::::::   :::    ::: :::::::::: */
/* types.go                                           :+:    :+: :+:    :+: :+:    :+: :+:    :+: :+:         */
/*                                                   +:+    +:+ +:+    +:+ +:+        +:+    +:+ +:+          */
/* By: ES-Yukun <yukun@doche.io>                    +#+    +:+ +#+    +:+ +#+        +#++:++#++ +#++:++#      */
/*                                                 +#+    +#+ +#+    +#+ +#+        +#+    +#+ +#+            */
/* Created: 2023/08/09 23:46:58 by ES-Yukun       #+#    #+# #+#    #+# #+#    #+# #+#    #+# #+#             */
/* Updated: 2023/08/09 23:46:59 by ES-Yukun      #########  ########   ########  ###    ### ##########.io.    */
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
