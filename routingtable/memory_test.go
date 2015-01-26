package routingtable

import (
	"github.com/robertkluin/message-flow/router"
	"testing"
)

func TestMemoryGetClientMessageServer(t *testing.T) {
	table := NewMemoryRoutingTable()

	// Client with no mapped services.
	table.SetClientMessageServer("client.2", "server.1")

	type TestCase struct {
		ClientID router.ClientID
		Result   router.ServerID
		Err      *router.RoutingTableError
	}

	tests := []TestCase{
		TestCase{"client.1", "", router.NewRoutingTableError(router.UnknownClient, "")},
		TestCase{"client.2", "server.1", nil},
	}

	// Load routing data

	for _, test := range tests {
		result, err := table.GetClientMessageServer(test.ClientID)
		if result != test.Result {
			t.Errorf("FAIL: Results didn't match.\n\tTest Case: %+v\n\tActual: {result: \"%v\", err: %+v}",
				test, result, err)
		} else if err != nil && test.Err == nil {
			t.Errorf("FAIL: Got an unexpected error.\n\tTest Case: %+v\n\tActual: {result: \"%v\", err: %+v}",
				test, result, err)
		} else if err == nil && test.Err != nil {
			t.Errorf("FAIL: Didn't get an expected error.\n\tTest Case: %+v\n\tActual: {result: \"%v\", err: %+v}",
				test, result, err)
		} else if err != nil && test.Err != nil && err.(*router.RoutingTableError).Code != test.Err.Code {
			t.Errorf("FAIL: Got the wrong error.\n\tTest Case: %+v\n\tActual: {result: \"%v\", err: %+v}",
				test, result, err)
		}
	}
}

func TestMemoryGetClientServiceServer(t *testing.T) {
	table := NewMemoryRoutingTable()

	table.SetClientMessageServer("client.2", "server.1")

	table.SetClientMessageServer("client.3", "server.2")
	clientRecord.setServiceServer("service.2", "server.1")

	type TestCase struct {
		ClientID  router.ClientID
		ServiceID router.ServiceID
		Result    router.ServerID
		Err       *router.RoutingTableError
	}

	tests := []TestCase{
		TestCase{"client.1", "service.1", "", router.NewRoutingTableError(router.UnknownClient, "")},
		TestCase{"client.2", "service.1", "", router.NewRoutingTableError(router.MappingNotFoundError, "")},
		TestCase{"client.3", "service.1", "", router.NewRoutingTableError(router.MappingNotFoundError, "")},
		TestCase{"client.3", "service.2", "server.1", nil},
	}

	// Load routing data

	for _, test := range tests {
		result, err := table.GetClientServiceServer(test.ClientID, test.ServiceID)
		if result != test.Result {
			t.Errorf("FAIL: Results didn't match.\n\tTest Case: %+v\n\tActual: {result: \"%v\", err: %+v}",
				test, result, err)
		} else if err != nil && test.Err == nil {
			t.Errorf("FAIL: Got an unexpected error.\n\tTest Case: %+v\n\tActual: {result: \"%v\", err: %+v}",
				test, result, err)
		} else if err == nil && test.Err != nil {
			t.Errorf("FAIL: Didn't get an expected error.\n\tTest Case: %+v\n\tActual: {result: \"%v\", err: %+v}",
				test, result, err)
		} else if err != nil && test.Err != nil && err.(*router.RoutingTableError).Code != test.Err.Code {
			t.Errorf("FAIL: Got the wrong error.\n\tTest Case: %+v\n\tActual: {result: \"%v\", err: %+v}",
				test, result, err)
		}
	}
}
