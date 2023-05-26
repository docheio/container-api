/* ********************************************************************************************************** */
/*                                                                                                            */
/*                                                     :::::::::  ::::::::   ::::::::   :::    ::: :::::::::: */
/* rfc1123.go                                         :+:    :+: :+:    :+: :+:    :+: :+:    :+: :+:         */
/*                                                   +:+    +:+ +:+    +:+ +:+        +:+    +:+ +:+          */
/* By: yushsato <yukun@team.anylinks.jp>            +#+    +:+ +#+    +:+ +#+        +#++:++#++ +#++:++#      */
/*                                                 +#+    +#+ +#+    +#+ +#+        +#+    +#+ +#+            */
/* Created: 2023/05/27 04:26:51 by yushsato       #+#    #+# #+#    #+# #+#    #+# #+#    #+# #+#             */
/* Updated: 2023/05/27 04:26:52 by yushsato      #########  ########   ########  ###    ### ##########.io.    */
/*                                                                                                            */
/* ********************************************************************************************************** */

package utils

import "math/rand"

func RFC1123() string {
	return randomStringSmallENG(1) + randomStringSmallENG09_(10) + randomStringSmallENG09(1)
}
func randomStringSmallENG09(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func randomStringSmallENG09_(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789-")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func randomStringSmallENG(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
