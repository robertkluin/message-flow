package routingtable

import (
	"github.com/robertkluin/message-flow/router"
)

// The memory based RoutingTable implements all core client, server, and
// service registration interfaces in memory.  It is suitable for use in a
// single node message-flow system that does not require persistence.

type clientRecord struct {
	SocketServer router.ServerID
	ServiceMap   map[router.ServiceID]router.ServerID
}

type clientTable map[router.ClientID]clientRecord

type MemoryRoutingTable struct {
	clientTable clientTable
}

func NewMemoryRoutingTable() *MemoryRoutingTable {
	table := new(MemoryRoutingTable)
	table.clientTable = make(clientTable)
	return table
}

// Which message server handles communication for client.
func (table *MemoryRoutingTable) GetClientMessageServer(clientID router.ClientID) (string, error) {
	_, ok := table.clientTable[clientID]

	if !ok {
		return "", router.NewRoutingTableError(router.UnknownClient, "No client routing info found.")
	}

	return "", nil
}

// Which server for service should messages from client be routed to.
func (table *MemoryRoutingTable) GetClientServiceServer(router.ClientID, router.ServiceID) (string, error) {
	return "", nil
}
