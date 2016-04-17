package dockerclient

import (
	"reflect"
	"testing"
)

func TestMockInterface(t *testing.T) {

	client := reflect.TypeOf((*Client)(nil)).Elem()
	mock := NewMockClient()
	if !reflect.TypeOf(mock).Implements(client) {
		t.Fatalf("Mock does not implement the Client interface")
	}
}
