package router

import (
	"fmt"
)

type RoutingTableErrorCode int

const (
	_                                  = iota
	ServiceError RoutingTableErrorCode = iota
	LookupError

	UnknownClient
	UnknownService

	ServerPoolEmptyError
	ServerNotFoundError
	MappingNotFoundError
)

type RoutingTableError struct {
	Code    RoutingTableErrorCode
	Message string
}

func (err RoutingTableError) Error() string {
	return fmt.Sprintf("%v (Routing Error Code: %d)", err.Message, err.Code)
}

func NewRoutingTableError(code RoutingTableErrorCode, message string) *RoutingTableError {
	return &RoutingTableError{Code: code, Message: message}
}

// A RoutingTable provides all core interfaces.
type RoutingTable interface {
	ClientTable
	ServiceTable
}

// Routing tables meeting the ClientTable spec answer basic questions about
// a client such as which server messages should be routed to for a given
// service and which message-flow front-end is handling communication with the
// client.
type ClientTable interface {
	// Which message server handles communication for client.
	GetClientMessageServer(ClientID) (ServerID, error)

	// Set the message server handling communication for the client.
	SetClientMessageServer(ClientID, ServerID) error

	// Which server for service should messages from client be routed to.
	GetClientServiceServer(ClientID, ServiceID) (ServerID, error)

	// Set server for service responsible for handling messages from client.
	SetClientServiceServer(ClientID, ServiceID, ServerID) error
}

// Routing tables meeting the ServiceTable spec answer questions about a
// service such as is there a single server defined for a service, is there a
// registrar defined, or what server should messages be directed to.
type ServiceTable interface {
	// Get the catch-all server, if defined, for the service.
	GetServiceServer(ServiceID) (ServerID, error)

	// Set a catch-all server for the service.
	SetServiceServer(ServiceID, ServerID) error

	// Get the registrar, if defined, for the service.
	GetServiceRegistrar(ServiceID) (ServerID, error)

	// Set the registrar for the service.
	SetServiceRegistrar(ServiceID, ServerID) error

	// Get a server from the pool of the service's registered servers
	GetServiceRandomServer(ServiceID) (ServerID, error)

	// Add a server to the service's server pool.
	AddServerToServicePool(ServiceID, ServerID) error

	// Remove a server from the service's pool of servers.
	RemoveServerFromServicePool(ServiceID, ServerID) error
}
