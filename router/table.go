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
