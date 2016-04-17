package dockerclient

import (
	"reflect"
	"testing"
	"github.com/gaia-adm/mr-burns/dockerclient"
)

func TestMockInterface(t *testing.T) {

	client := reflect.TypeOf((*dockerclient.Client)(nil)).Elem()
	mock := NewMockClient()
	if !reflect.TypeOf(mock).Implements(client) {
		t.Fatalf("Mock does not implement the Client interface")
	}
}
