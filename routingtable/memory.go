package routingtable

import (
	"github.com/robertkluin/message-flow/router"
)

// The memory based RoutingTable implements all core client, server, and
// service registration interfaces in memory.  It is suitable for use in a
// single node message-flow system that does not require persistence.

type clientRecord struct {
	MessageServer router.ServerID
	ServiceMap    serviceMap
}

type serviceMap map[router.ServiceID]router.ServerID

func newClientRecord(messageServer router.ServerID) *clientRecord {
	record := new(clientRecord)
	record.MessageServer = messageServer
	record.ServiceMap = make(serviceMap)
	return record
}

type clientTable map[router.ClientID]*clientRecord

type MemoryRoutingTable struct {
	clientTable clientTable
}

func NewMemoryRoutingTable() *MemoryRoutingTable {
	table := new(MemoryRoutingTable)
	table.clientTable = make(clientTable)
	return table
}

func (r *clientRecord) getServiceServer(serviceID router.ServiceID) (router.ServerID, error) {
	serverID, ok := r.ServiceMap[serviceID]

	if !ok {
		return "", router.NewRoutingTableError(router.MappingNotFoundError, "No server found for service.")
	}

	return serverID, nil
}

func (table *MemoryRoutingTable) getClientRecord(clientID router.ClientID) (*clientRecord, error) {
	record, ok := table.clientTable[clientID]

	if !ok {
		return nil, router.NewRoutingTableError(router.UnknownClient, "No client routing info found.")
	}

	return record, nil
}

// Which message server handles communication for client.
func (table *MemoryRoutingTable) GetClientMessageServer(clientID router.ClientID) (router.ServerID, error) {
	record, err := table.getClientRecord(clientID)

	if err != nil {
		return "", err
	}

	return record.MessageServer, nil
}

// Which server for service should messages from client be routed to.
func (table *MemoryRoutingTable) GetClientServiceServer(clientID router.ClientID, serviceID router.ServiceID) (router.ServerID, error) {
	record, err := table.getClientRecord(clientID)

	if err != nil {
		return "", err
	}

	serverID, err := record.getServiceServer(serviceID)

	if err != nil {
		return "", err
	}

	return serverID, nil
}
