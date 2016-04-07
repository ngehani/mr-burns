package mockclient

import (
	"reflect"
	"testing"
	"github.com/gaia-adm/mr-burns/container"
)

func TestMockInterface(t *testing.T) {
	iface := reflect.TypeOf((*container.Client)(nil)).Elem()
	mock := NewMockClient()

	if !reflect.TypeOf(mock).Implements(iface) {
		t.Fatalf("Mock does not implement the Client interface")
	}
}
