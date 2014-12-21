package routingtable

import (
	"github.com/robertkluin/message-flow/router"
	"testing"
)

type TableTestCase struct {
	ClientID router.ClientID
	Result   router.ServerID
	Err      *router.RoutingTableError
}

func TestMemoryGetClientMessageServer(t *testing.T) {
	table := NewMemoryRoutingTable()

	table.clientTable["CLIENT"] = newClientRecord("server1")

	tests := []TableTestCase{
		TableTestCase{"NEWCLIENT", "", router.NewRoutingTableError(router.UnknownClient, "")},
		TableTestCase{"CLIENT", "server1", nil},
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
