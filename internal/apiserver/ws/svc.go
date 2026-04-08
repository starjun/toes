// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ws 提供 WebSocket 服务。
//
// 该包实现 WebSocket 服务器，支持消息广播、
// 客户端管理、连接池等功能。
//
// 主要组件:
//   - Hub: 消息中心
//   - Client: 客户端连接
//   - Service: WebSocket 服务
//
// 使用示例:
//
//	ws.StartWS()
//	ws.Broadcast("message")
package ws

import (
	"time"

	"toes/internal/apiserver/sysinfo"
	"toes/internal/utils"
)

var hub *Hub

func GetHub() *Hub {
	return hub
}

func StartWS() {
	hub = newHub()
	go hub.run()

	// go PushMessage(hub)
}

func PushMessage(h *Hub) {
	for {
		h.SendMessage(utils.JsonEncodeIndent(sysinfo.GetCpuInfo()))
		time.Sleep(time.Second * 10)
	}
}
