package routingtable

import (
	"github.com/robertkluin/message-flow/router"
	"math/rand"
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

// Get the catch-all server, if defined, for the service.
func (table *MemoryRoutingTable) GetServiceServer(serviceID router.ServiceID) (router.ServerID, error) {
	record, err := table.getServiceRecord(serviceID)
	if err != nil {
		return "", err
	}

	serverID, err := record.getServer()
	if err != nil {
		return "", err
	}

	return serverID, nil
}

//  Set a catch-all server for the service.
func (table *MemoryRoutingTable) SetServiceServer(serviceID router.ServiceID, serverID router.ServerID) error {
	record, err := table.getOrCreateServiceRecord(serviceID)
	if err != nil {
		return err
	}

	err = record.setServer(serverID)
	if err != nil {
		return err
	}

	return nil
}

// Get the registrar, if defined, for the service.
func (table *MemoryRoutingTable) GetServiceRegistrar(serviceID router.ServiceID) (router.ServerID, error) {
	record, err := table.getServiceRecord(serviceID)
	if err != nil {
		return "", err
	}

	serverID, err := record.getRegistrar()
	if err != nil {
		return "", err
	}

	return serverID, nil
}

// Set the registrar for the service.
func (table *MemoryRoutingTable) SetServiceRegistrar(serviceID router.ServiceID, serverID router.ServerID) error {
	record, err := table.getOrCreateServiceRecord(serviceID)
	if err != nil {
		return err
	}

	err = record.setRegistrar(serverID)
	if err != nil {
		return err
	}

	return nil
}

// Get a server from the pool of the service's registered servers
func (table *MemoryRoutingTable) GetServiceRandomServer(serviceID router.ServiceID) (router.ServerID, error) {
	record, err := table.getServiceRecord(serviceID)
	if err != nil {
		return "", err
	}

	serverID, err := record.getServerFromPool()
	if err != nil {
		return "", err
	}

	return serverID, nil
}

// Add a server to the service's server pool.
func (table *MemoryRoutingTable) AddServerToServicePool(serviceID router.ServiceID, serverID router.ServerID) error {
	record, err := table.getOrCreateServiceRecord(serviceID)
	if err != nil {
		return err
	}

	err = record.addServerToPool(serverID)
	if err != nil {
		return err
	}

	return nil
}

// Remove a server from the service's pool of servers.
func (table *MemoryRoutingTable) RemoveServerFromServicePool(serviceID router.ServiceID, serverID router.ServerID) error {
	record, err := table.getOrCreateServiceRecord(serviceID)
	if err != nil {
		return err
	}

	err = record.removeServerFromPool(serverID)
	if err != nil {
		return err
	}

	return nil
}

// Insert new client record in routing table
func (table *MemoryRoutingTable) getOrCreateClientRecord(clientID router.ClientID) (*clientRecord, error) {
	record, ok := table.clientTable[clientID]

	if !ok {
		record = newClientRecord()
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

// Insert new service record in routing table
func (table *MemoryRoutingTable) getOrCreateServiceRecord(serviceID router.ServiceID) (*serviceRecord, error) {
	record, ok := table.serviceTable[serviceID]

	if !ok {
		record = newServiceRecord()
		table.serviceTable[serviceID] = record
	}

	return record, nil
}

// Lookup service information in routing table
func (table *MemoryRoutingTable) getServiceRecord(serviceID router.ServiceID) (*serviceRecord, error) {
	record, ok := table.serviceTable[serviceID]

	if !ok {
		return nil, router.NewRoutingTableError(router.UnknownService, "No service routing info found.")
	}

	return record, nil
}

// Routing information tracked per client
type clientRecord struct {
	messageServer router.ServerID
	serviceMap    serviceMap
}

type serviceMap map[router.ServiceID]router.ServerID

func newClientRecord() *clientRecord {
	record := new(clientRecord)
	record.messageServer = ""
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
	registrar  router.ServerID
	serverPool serverList
}

type serverList []router.ServerID

func (l *serverList) find(serverID router.ServerID) int {
	for i, id := range *l {
		if id == serverID {
			return i
		}
	}
	return -1
}

func (l *serverList) add(serverID router.ServerID) {
	if l.find(serverID) >= 0 {
		return
	}

	*l = append(*l, serverID)
}

func (l *serverList) remove(serverID router.ServerID) {
	pos := l.find(serverID)
	if pos == -1 {
		return
	}

	*l = append((*l)[:pos], (*l)[pos+1:]...)
}

func newServiceRecord() *serviceRecord {
	record := new(serviceRecord)
	record.server = ""
	record.registrar = ""
	record.serverPool = make(serverList, 0, 10)
	return record
}

type serviceTable map[router.ServiceID]*serviceRecord

func (r *serviceRecord) getServer() (router.ServerID, error) {
	if r.server == "" {
		return "", router.NewRoutingTableError(router.ServerNotFoundError, "No catch-all server defined for service.")
	}

	return r.server, nil
}

func (r *serviceRecord) setServer(serverID router.ServerID) error {
	r.server = serverID

	return nil
}

func (r *serviceRecord) getRegistrar() (router.ServerID, error) {
	if r.registrar == "" {
		return "", router.NewRoutingTableError(router.ServerNotFoundError, "No registrar defined for service.")
	}

	return r.registrar, nil
}

func (r *serviceRecord) setRegistrar(registrar router.ServerID) error {
	r.registrar = registrar

	return nil
}

func (r *serviceRecord) getServerFromPool() (router.ServerID, error) {
	pool_size := len(r.serverPool)
	if len(r.serverPool) == 0 {
		return "", router.NewRoutingTableError(router.ServerPoolEmptyError, "No servers in pool.")
	}

	return r.serverPool[rand.Intn(pool_size)], nil
}

func (r *serviceRecord) addServerToPool(serverID router.ServerID) error {
	r.serverPool.add(serverID)

	return nil
}

func (r *serviceRecord) removeServerFromPool(serverID router.ServerID) error {
	r.serverPool.remove(serverID)

	return nil
}
