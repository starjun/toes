// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"time"
	"toes/internal/sysinfo"
	"toes/internal/utils"
)

var hub *Hub

func GetHub() *Hub {
	return hub
}

func StartWS() {
	hub = newHub()
	go hub.run()

	//go PushMessage(hub)
}

func PushMessage(h *Hub) {
	for {
		h.SendMessage(utils.JsonEncodeIndent(sysinfo.GetCpuInfo()))
		time.Sleep(time.Second * 10)
	}
}
