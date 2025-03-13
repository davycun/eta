package ws

import (
	"sync"
)

var (
	HUB = NewHub()
)

// Hub ws Hub
type Hub struct {
	Clients    *sync.Map //map[string]map[*Client]bool
	Message    chan *WsMessage
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    new(sync.Map),
		Message:    make(chan *WsMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) AddClient(client *Client) {
	userId := client.UserId
	clients, ok := h.Clients.Load(userId)
	if !ok {
		cs := new(sync.Map)
		cs.Store(client, true)
		actual, loaded := h.Clients.LoadOrStore(userId, cs)
		if loaded {
			actual1 := actual.(*sync.Map)
			actual1.Store(client, true)
		}
	} else {
		c := clients.(*sync.Map)
		c.Store(client, true)
	}
}

func (h *Hub) RemoveClient(client *Client) {
	userId := client.UserId
	clients, ok := h.Clients.Load(userId)
	if !ok {
		return
	}
	cs := clients.(*sync.Map)
	cs.Delete(client)

	cnt := 0
	cs.Range(func(k, v interface{}) bool {
		cnt += 1
		return true
	})
	if cnt == 0 {
		h.Clients.Delete(userId)
	}
	if client.Send != nil {
		close(client.Send)
		client.Send = nil
	}
}

// GetClients 返回指定用户的客户端列表
func (h *Hub) GetClients(userId string) []*Client {
	// 返回 h.clients[userId] 里所有的 Client slice
	clients := make([]*Client, 0)
	userClients, ok := h.Clients.Load(userId)
	if !ok {
		return clients
	}
	uc := userClients.(*sync.Map)
	uc.Range(func(k, v interface{}) bool {
		clients = append(clients, k.(*Client))
		return true
	})
	return clients
}

// AllClients 返回客户端列表
func (h *Hub) AllClients() []*Client {
	// 返回 h.clients[userId] 里所有的 Client slice
	clients := make([]*Client, 0)
	h.Clients.Range(func(k, v interface{}) bool {
		v.(*sync.Map).Range(func(k, v interface{}) bool {
			clients = append(clients, k.(*Client))
			return true
		})
		return true
	})
	return clients
}

// Run 运行hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.AddClient(client)
		case client := <-h.Unregister:
			h.RemoveClient(client)
		case message := <-h.Message:
			toSendClients := make([]*Client, 0)
			if len(message.UserId) == 0 {
				toSendClients = h.AllClients()
			} else {
				for _, userId := range message.UserId {
					toSendClients = append(toSendClients, h.GetClients(userId)...)
				}
			}
			for _, client := range toSendClients {
				select {
				case client.Send <- message:
				default:
					h.RemoveClient(client)
				}
			}
		}
	}
}
