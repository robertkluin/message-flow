package routingtable

import (
	"github.com/robertkluin/message-flow/router"
)

// The memory based RoutingTable implements all core client, server, and
// service registration interfaces in memory.  It is suitable for use in a
// single node message-flow system that does not require persistence.

type MemoryRoutingTable struct {
}

func NewMemoryRoutingTable() *MemoryRoutingTable {
	return new(MemoryRoutingTable)
}

// Which message server handles communication for client.
func GetClientMessageServer(router.ClientID) (string, error) {
	return "", nil
}

// Which server for service should messages from client be routed to.
func GetClientServiceServer(router.ClientID, router.ServiceID) (string, error) {
	return "", nil
}
