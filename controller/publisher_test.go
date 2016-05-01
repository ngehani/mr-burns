package controller

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
)

func TestPublisher_Publish(t *testing.T) {

	var actual string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if body, err := ioutil.ReadAll(r.Body); err != nil {
			t.Fail()
		} else {
			actual = string(body)
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	const EXPECTED = "feed me seymour"
	NewPublisher(server.URL).Publish(EXPECTED)
	assert.Equal(t, EXPECTED, actual);
}
