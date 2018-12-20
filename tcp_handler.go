package mhist

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/codeuniversity/ppp-mhist/models"
	"github.com/codeuniversity/ppp-mhist/tcp"
)

//TCPHandler handles tcp connections
type TCPHandler struct {
	address                     string
	outboundCollection          *tcp.ConnectionCollection
	server                      *Server
	filterPerOutboundConnection map[*tcp.Connection]*models.FilterCollection
	filterMutex                 *sync.RWMutex
	pools                       *models.Pools
}

//NewTCPHandler sets the wrapped handlers callbacks correctly, Run() still has to be called
func NewTCPHandler(server *Server, port int, pools *models.Pools) *TCPHandler {
	return &TCPHandler{
		address:                     fmt.Sprintf("0.0.0.0:%v", port),
		server:                      server,
		outboundCollection:          &tcp.ConnectionCollection{},
		filterMutex:                 &sync.RWMutex{},
		filterPerOutboundConnection: make(map[*tcp.Connection]*models.FilterCollection),
		pools:                       pools,
	}
}

//Notify handler about new message
func (h *TCPHandler) Notify(name string, measurement models.Measurement) {
	m := h.pools.GetMessage()
	defer h.pools.PutMessage(m)

	m.Reset()
	m.Name = name
	m.Value = measurement.ValueInterface()
	m.Timestamp = measurement.Timestamp()

	byteSlice, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return
	}
	h.filterMutex.RLock()
	defer h.filterMutex.RUnlock()
	h.outboundCollection.ForEach(func(conn *tcp.Connection) {
		filter := h.filterPerOutboundConnection[conn]
		if filter != nil {
			if filter.Passes(name, measurement) {
				conn.Write(byteSlice)
			}
		} else {
			fmt.Println("Filter for outbound connection was nil, please investigate!")
		}
	})
}

//Run listens for new connections
func (h *TCPHandler) Run() {
	listener, err := net.Listen("tcp", h.address)
	if err != nil {
		panic("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go h.handleNewConnection(conn)
	}
}

func (h *TCPHandler) onNewMessage(byteSlice []byte, isReplication bool) {
	h.server.handleNewMessage(byteSlice, isReplication, func(err error, _ int) {
		if err != nil {
			fmt.Println(err)
		}
	})
}

func (h *TCPHandler) handleNewConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	byteSlice, err := reader.ReadSlice('\n')
	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}
	m := &models.SubscriptionMessage{}
	err = json.Unmarshal(byteSlice, m)
	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}
	connectionWrapper := &tcp.Connection{
		Socket: conn,
		Reader: reader,
	}
	if m.Publisher {
		connectionWrapper.OnNewMessage(func(byteSlice []byte) {
			h.onNewMessage(byteSlice, m.Replication)
		})
	} else {
		h.addFilterForConnection(m.FilterDefinition, connectionWrapper)
		h.outboundCollection.AddConnection(connectionWrapper)
		connectionWrapper.OnConnectionClose(func() {
			h.outboundCollection.RemoveConnection(connectionWrapper)
			h.removeFilterForConnection(connectionWrapper)
		})
	}
	connectionWrapper.Listen()
}

func (h *TCPHandler) removeFilterForConnection(conn *tcp.Connection) {
	h.filterMutex.Lock()
	defer h.filterMutex.Unlock()

	delete(h.filterPerOutboundConnection, conn)
}

func (h *TCPHandler) addFilterForConnection(filterDefinition models.FilterDefinition, conn *tcp.Connection) {
	h.filterMutex.Lock()
	defer h.filterMutex.Unlock()

	filter := models.NewFilterCollection(filterDefinition)
	h.filterPerOutboundConnection[conn] = filter
}
