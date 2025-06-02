package sse

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type Hub struct {
	clients        sync.Map // clientId -> client
	roomClientsMap sync.Map // roomId -> client_list
}

func NewHub() *Hub {
	return &Hub{
		clients:        sync.Map{},
		roomClientsMap: sync.Map{},
	}
}

func (h *Hub) AddClient(client *Client) {
	h.clients.Store(client.ID, client)
	clientList, ok := h.roomClientsMap.Load(client.RoomID)
	if !ok {
		clientList = make([]string, 0, 50)
		h.roomClientsMap.Store(client.RoomID, clientList)
	}
	clientList = append(clientList.([]string), client.ID)
	h.roomClientsMap.Store(client.RoomID, clientList)
}

func (h *Hub) RemoveClient(id string) {
	c, ok := h.clients.Load(id)
	if !ok {
		logrus.Infof("sse client %s not found", id)
		return
	}
	client := c.(*Client)

	clientList, ok := h.roomClientsMap.Load(client.RoomID)
	if !ok {
		return
	}
	for i, id := range clientList.([]string) {
		if id == client.ID {
			clientList = append(clientList.([]string)[:i], clientList.([]string)[i+1:]...)
			break
		}
	}
	h.roomClientsMap.Store(client.RoomID, clientList)

	// 最后再删除clients
	h.clients.Delete(id)
}

func (h *Hub) PushMessage(roomId int64, data []byte) {

}
