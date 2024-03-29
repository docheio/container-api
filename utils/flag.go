/* ********************************************************************************************************** */
/*                                                                                                            */
/*                                                     :::::::::  ::::::::   ::::::::   :::    ::: :::::::::: */
/* flag.go                                            :+:    :+: :+:    :+: :+:    :+: :+:    :+: :+:         */
/*                                                   +:+    +:+ +:+    +:+ +:+        +:+    +:+ +:+          */
/* By: ES-Yukun <yukun@doche.io>                    +#+    +:+ +#+    +:+ +#+        +#++:++#++ +#++:++#      */
/*                                                 +#+    +#+ +#+    +#+ +#+        +#+    +#+ +#+            */
/* Created: 2023/08/09 23:47:39 by ES-Yukun       #+#    #+# #+#    #+# #+#    #+# #+#    #+# #+#             */
/* Updated: 2023/08/09 23:47:40 by ES-Yukun      #########  ########   ########  ###    ### ##########.io.    */
/*                                                                                                            */
/* ********************************************************************************************************** */

package utils

import (
	"flag"
	"log"
)

type Flags struct {
	Namespace *string
	Uniquekey *string
	Image     *string
	Address   *string
	Debug     *bool
}

func (flags *Flags) Init() {
	flags.Namespace = flag.String("namespace", "", "Namespace to operate the container.")
	flags.Uniquekey = flag.String("key", "", "Unique key to identify to which group the container belongs.")
	flags.Image = flag.String("image", "", "Container image url.")
	flags.Address = flag.String("address", ":8080", "IP:Port for hosting api.")
	flags.Debug = flag.Bool("debug", false, "Debug mode.")
	flag.Parse()
	if *flags.Namespace == "" {
		log.Fatalln("Namespace must be specified.")
	}
	if *flags.Uniquekey == "" {
		log.Fatalln("key must be specified.")
	}
	if *flags.Image == "" {
		log.Fatalln("image must be specified.")
	}
}
