package routingtable

import (
	"github.com/robertkluin/message-flow/router"
)

// `MemoryRoutingTable` implements all core client, server, and service
// registration interfaces in memory.  It is suitable for use in a single node
// message-flow system that does not require persistence.

type MemoryRoutingTable struct {
	clientTable  clientTable
	serviceTable serviceTable
}

func NewMemoryRoutingTable() *MemoryRoutingTable {
	table := new(MemoryRoutingTable)
	table.clientTable = make(clientTable)
	table.serviceTable = make(serviceTable)
	return table
}

// Which message server handles communication for client.
func (table *MemoryRoutingTable) GetClientMessageServer(clientID router.ClientID) (router.ServerID, error) {
	record, err := table.getClientRecord(clientID)
	if err != nil {
		return "", err
	}

	server, err := record.getMessageServer()
	if err != nil {
		return "", err
	}

	return server, nil
}

// Set the message server that handles communication for the client.
func (table *MemoryRoutingTable) SetClientMessageServer(clientID router.ClientID, messageServer router.ServerID) error {
	record, err := table.getOrCreateClientRecord(clientID)
	if err != nil {
		return err
	}

	err = record.setMessageServer(messageServer)
	if err != nil {
		return err
	}

	return nil
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

// Set server for service responsible for handling messages from client.
func (table *MemoryRoutingTable) SetClientServiceServer(clientID router.ClientID, serviceID router.ServiceID, serverID router.ServerID) error {
	record, err := table.getOrCreateClientRecord(clientID)
	if err != nil {
		return err
	}

	err = record.setServiceServer(serviceID, serverID)
	if err != nil {
		return err
	}

	return nil
}

// Insert new client record in routing table
func (table *MemoryRoutingTable) getOrCreateClientRecord(clientID router.ClientID) (*clientRecord, error) {
	record, ok := table.clientTable[clientID]

	if !ok {
		record = newClientRecord("")
		table.clientTable[clientID] = record
	}

	return record, nil
}

// Lookup client information in routing table
func (table *MemoryRoutingTable) getClientRecord(clientID router.ClientID) (*clientRecord, error) {
	record, ok := table.clientTable[clientID]

	if !ok {
		return nil, router.NewRoutingTableError(router.UnknownClient, "No client routing info found.")
	}

	return record, nil
}

// Routing information tracked per client
type clientRecord struct {
	messageServer router.ServerID
	serviceMap    serviceMap
}

type serviceMap map[router.ServiceID]router.ServerID

func newClientRecord(messageServer router.ServerID) *clientRecord {
	record := new(clientRecord)
	record.messageServer = messageServer
	record.serviceMap = make(serviceMap)
	return record
}

type clientTable map[router.ClientID]*clientRecord

func (r *clientRecord) getMessageServer() (router.ServerID, error) {
	if r.messageServer == "" {
		return "", router.NewRoutingTableError(router.MappingNotFoundError, "No message server found for client.")
	}
	return r.messageServer, nil
}

func (r *clientRecord) setMessageServer(serverID router.ServerID) error {
	r.messageServer = serverID
	return nil
}

func (r *clientRecord) getServiceServer(serviceID router.ServiceID) (router.ServerID, error) {
	serverID, ok := r.serviceMap[serviceID]

	if !ok {
		return "", router.NewRoutingTableError(router.MappingNotFoundError, "No server found for service.")
	}

	return serverID, nil
}

func (r *clientRecord) setServiceServer(serviceID router.ServiceID, serverID router.ServerID) error {
	r.serviceMap[serviceID] = serverID
	return nil
}

// Routing information tracked per service
type serviceRecord struct {
	server     router.ServerID
}

func newServiceRecord(server router.ServerID) *serviceRecord {
	record := new(serviceRecord)
	record.server = server
	return record
}

type serviceTable map[router.ServiceID]*serviceRecord

