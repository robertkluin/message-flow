package routingtable

import (
	"github.com/robertkluin/message-flow/router"
	"testing"
)

func TestMemoryGetClientMessageServer(t *testing.T) {
	table := NewMemoryRoutingTable()
	router.TestGetClientMessageServer(t, table)
}

func TestMemoryGetClientServiceServer(t *testing.T) {
	table := NewMemoryRoutingTable()
	router.TestGetClientServiceServer(t, table)
}

func TestMemoryGetServiceServer(t *testing.T) {
	table := NewMemoryRoutingTable()
	router.TestGetServiceServer(t, table)
}

func TestMemoryGetServiceRegistrar(t *testing.T) {
	table := NewMemoryRoutingTable()
	router.TestGetServiceRegistrar(t, table)
}

func TestMemoryGetServiceRandomServer(t *testing.T) {
	table := NewMemoryRoutingTable()
	router.TestGetServiceRandomServer(t, table)
}
