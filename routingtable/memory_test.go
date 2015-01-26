package routingtable

import (
	"github.com/robertkluin/message-flow/router"
	"testing"
)

func TestMemoryGetClientMessageServer(t *testing.T) {
	table := NewMemoryRoutingTable()

	// Client with no mapped services.
	table.SetClientMessageServer("client.2", "server.1")

	// Client with a mapped service, but no message server.
	table.SetClientServiceServer("client.3", "service.1", "server.1")

	type TestCase struct {
		ClientID router.ClientID
		Result   router.ServerID
		Err      *router.RoutingTableError
	}

	tests := []TestCase{
		// client.1 does not exist, there are no mappings for it.
		TestCase{"client.1", "", router.NewRoutingTableError(router.UnknownClient, "")},

		// client.2 messages are mapped to server.1.
		TestCase{"client.2", "server.1", nil},

		// client.3 has service mappings, but no message server defined.
		TestCase{"client.3", "", router.NewRoutingTableError(router.MappingNotFoundError, "")},
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

	// Client with no mapped services.
	table.SetClientMessageServer("client.2", "server.1")
	//table.clientTable["client.2"] = newClientRecord("server.1")

	// Client with service.2 mapped.
	table.SetClientServiceServer("client.3", "service.2", "server.1")
	table.SetClientMessageServer("client.3", "server.2")

	type TestCase struct {
		ClientID  router.ClientID
		ServiceID router.ServiceID
		Result    router.ServerID
		Err       *router.RoutingTableError
	}

	tests := []TestCase{
		// client.1 does not exist, there are no mappings for it.
		TestCase{"client.1", "service.1", "", router.NewRoutingTableError(router.UnknownClient, "")},

		// client.2 exists, but there are no service mappings for it.
		TestCase{"client.2", "service.1", "", router.NewRoutingTableError(router.MappingNotFoundError, "")},

		// client.3 exists, but there is no service mappings for service.1.
		TestCase{"client.3", "service.1", "", router.NewRoutingTableError(router.MappingNotFoundError, "")},

		// client.3 exists and there is a mapping to service.2 to server.1.
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
