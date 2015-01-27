package routingtable

import (
	"github.com/robertkluin/message-flow/router"
	"testing"
)

type TestCase struct {
	Args   []interface{}
	Result router.ServerID
	Err    *router.RoutingTableError
}

func evalTests(t *testing.T, tests []TestCase, runner func(test TestCase) (router.ServerID, error)) {
	for _, test := range tests {
		//result, err := table.GetClientMessageServer(test.ClientID)
		result, err := runner(test)
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

func TestMemoryGetClientMessageServer(t *testing.T) {
	table := NewMemoryRoutingTable()

	// Client with no mapped services.
	table.SetClientMessageServer("client.2", "server.1")

	// Client with a mapped service, but no message server.
	table.SetClientServiceServer("client.3", "service.1", "server.1")

	mkArgs := func(clientID router.ClientID) []interface{} {
		return []interface{}{clientID}
	}

	tests := []TestCase{
		// client.1 does not exist, there are no mappings for it.
		TestCase{mkArgs("client.1"), "", router.NewRoutingTableError(router.UnknownClient, "")},

		// client.2 messages are mapped to server.1.
		TestCase{mkArgs("client.2"), "server.1", nil},

		// client.3 has service mappings, but no message server defined.
		TestCase{mkArgs("client.3"), "", router.NewRoutingTableError(router.MappingNotFoundError, "")},
	}

	evalTests(t, tests, func(test TestCase) (router.ServerID, error) {
		arg, _ := test.Args[0].(router.ClientID)
		return table.GetClientMessageServer(arg)
	})
}

func TestMemoryGetClientServiceServer(t *testing.T) {
	table := NewMemoryRoutingTable()

	// Client with no mapped services.
	table.SetClientMessageServer("client.2", "server.1")
	//table.clientTable["client.2"] = newClientRecord("server.1")

	// Client with service.2 mapped.
	table.SetClientServiceServer("client.3", "service.2", "server.1")
	table.SetClientMessageServer("client.3", "server.2")

	mkArgs := func(clientID router.ClientID, serviceID router.ServiceID) []interface{} {
		return []interface{}{clientID, serviceID}
	}

	tests := []TestCase{
		// client.1 does not exist, there are no mappings for it.
		TestCase{mkArgs("client.1", "service.1"), "", router.NewRoutingTableError(router.UnknownClient, "")},

		// client.2 exists, but there are no service mappings for it.
		TestCase{mkArgs("client.2", "service.1"), "", router.NewRoutingTableError(router.MappingNotFoundError, "")},

		// client.3 exists, but there is no service mappings for service.1.
		TestCase{mkArgs("client.3", "service.1"), "", router.NewRoutingTableError(router.MappingNotFoundError, "")},

		// client.3 exists and there is a mapping to service.2 to server.1.
		TestCase{mkArgs("client.3", "service.2"), "server.1", nil},
	}

	evalTests(t, tests, func(test TestCase) (router.ServerID, error) {
		clientID, _ := test.Args[0].(router.ClientID)
		serviceID, _ := test.Args[1].(router.ServiceID)
		return table.GetClientServiceServer(clientID, serviceID)
	})
}

func TestMemoryGetServiceServer(t *testing.T) {
	table := NewMemoryRoutingTable()

	// Service with a catch-all server.
	table.SetServiceServer("service.2", "server.1")

	// Service with an empty catch-all server.
	table.SetServiceServer("service.3", "")

	mkArgs := func(serviceID router.ServiceID) []interface{} {
		return []interface{}{serviceID}
	}

	tests := []TestCase{
		// service.1 does not exist, there is no mapping.
		TestCase{mkArgs("service.1"), "", router.NewRoutingTableError(router.UnknownService, "")},

		// service.2 is mapped to server.1.
		TestCase{mkArgs("service.2"), "server.1", nil},

		// service.3 has a mappings, but no catch-all server defined.
		TestCase{mkArgs("service.3"), "", router.NewRoutingTableError(router.ServerNotFoundError, "")},
	}

	evalTests(t, tests, func(test TestCase) (router.ServerID, error) {
		serviceID, _ := test.Args[0].(router.ServiceID)
		return table.GetServiceServer(serviceID)
	})
}

func TestMemoryGetServiceRegistrar(t *testing.T) {
	table := NewMemoryRoutingTable()

	// Service with a catch-all server, but no registrar.
	table.SetServiceServer("service.2", "server.1")

	// Service with a registrar, and no server.
	table.SetServiceRegistrar("service.3", "registrar.1")

	// Service with a registrar and server.
	table.SetServiceServer("service.4", "server.2")
	table.SetServiceRegistrar("service.4", "registrar.2")

	// Service with an empty registrar.
	table.SetServiceRegistrar("service.5", "")

	mkArgs := func(serviceID router.ServiceID) []interface{} {
		return []interface{}{serviceID}
	}

	tests := []TestCase{
		// service.1 does not exist, there is no mapping.
		TestCase{mkArgs("service.1"), "", router.NewRoutingTableError(router.UnknownService, "")},

		// service.2 has a server, but no registrar.
		TestCase{mkArgs("service.2"), "", router.NewRoutingTableError(router.ServerNotFoundError, "")},

		// service.3 has a registrar.
		TestCase{mkArgs("service.3"), "registrar.1", nil},

		// service.4 has a server and a registrar.
		TestCase{mkArgs("service.4"), "registrar.2", nil},

		// service.5 has an empty registrar.
		TestCase{mkArgs("service.5"), "", router.NewRoutingTableError(router.ServerNotFoundError, "")},
	}

	evalTests(t, tests, func(test TestCase) (router.ServerID, error) {
		serviceID, _ := test.Args[0].(router.ServiceID)
		return table.GetServiceRegistrar(serviceID)
	})
}

func TestMemoryGetServiceRandomServer(t *testing.T) {
	table := NewMemoryRoutingTable()

	// Service with a catch-all server, but no server pool.
	table.SetServiceServer("service.2", "server.1")

	// Service with a registrar, and no server pool.
	table.SetServiceRegistrar("service.3", "registrar.1")

	// Service with a registrar and server, but no server pool.
	table.SetServiceServer("service.4", "server.2")
	table.SetServiceRegistrar("service.4", "registrar.2")

	// Service with a registrar and server and one server in pool.
	table.SetServiceServer("service.5", "server.2")
	table.SetServiceRegistrar("service.5", "registrar.2")
	table.AddServerToServicePool("service.5", "pool.1")

	// Service with only one server in pool.
	table.AddServerToServicePool("service.6", "pool.1")

	// Service with server removed from pool
	table.AddServerToServicePool("service.7", "pool.1")
	table.AddServerToServicePool("service.7", "pool.2")
	table.AddServerToServicePool("service.7", "pool.3")

	table.RemoveServerFromServicePool("service.7", "pool.1")
	table.RemoveServerFromServicePool("service.7", "pool.2")
	table.RemoveServerFromServicePool("service.7", "pool.3")

	// Service with all but one servers removed from pool
	table.AddServerToServicePool("service.8", "pool.1")
	table.AddServerToServicePool("service.8", "pool.2")
	table.RemoveServerFromServicePool("service.8", "pool.1")

	// Tests of serverList add/remove
	table.RemoveServerFromServicePool("service.9", "pool.1")
	table.AddServerToServicePool("service.9", "pool.1")
	table.AddServerToServicePool("service.9", "pool.1")
	table.RemoveServerFromServicePool("service.9", "pool.1")

	mkArgs := func(serviceID router.ServiceID) []interface{} {
		return []interface{}{serviceID}
	}

	tests := []TestCase{
		// service.1 does not exist, there is no mapping.
		TestCase{mkArgs("service.1"), "", router.NewRoutingTableError(router.UnknownService, "")},

		// service.2 has a server, but no server pool.
		TestCase{mkArgs("service.2"), "", router.NewRoutingTableError(router.ServerPoolEmptyError, "")},

		// service.3 has a registrar, but no server pool.
		TestCase{mkArgs("service.3"), "", router.NewRoutingTableError(router.ServerPoolEmptyError, "")},

		// service.4 has a server and registrar, but no server pool.
		TestCase{mkArgs("service.4"), "", router.NewRoutingTableError(router.ServerPoolEmptyError, "")},

		// service.5 has a server, registrar, and single server in the pool.
		TestCase{mkArgs("service.5"), "pool.1", nil},

		// service.6 has only a single server in the pool.
		TestCase{mkArgs("service.6"), "pool.1", nil},

		// service.7 has an emptied server pool.
		TestCase{mkArgs("service.7"), "", router.NewRoutingTableError(router.ServerPoolEmptyError, "")},

		// service.8 has servers added and removed, but one left in the pool.
		TestCase{mkArgs("service.8"), "pool.2", nil},

		// service.9 should have an emptied server pool.
		TestCase{mkArgs("service.9"), "", router.NewRoutingTableError(router.ServerPoolEmptyError, "")},
	}

	evalTests(t, tests, func(test TestCase) (router.ServerID, error) {
		serviceID, _ := test.Args[0].(router.ServiceID)
		return table.GetServiceRandomServer(serviceID)
	})
}
