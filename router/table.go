package router

// A RoutingTable provides all core interfaces.
type RoutingTable interface {
	ClientTable
}

// Routing tables meeting the ClientTable spec answer basic questions about
// a client such as which server messages should be routed to for a given
// service and which message-flow front-end is handling communication with the
// client.
type ClientTable interface {
	// Get client information
	// Which message server handles communication for client.
	GetClientMessageServer(ClientID) (string, error)

	// Which server for service should messages from client be routed to.
	GetClientServiceServer(ClientID, ServiceID) (string, error)
}
