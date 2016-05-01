package common

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewConfiguration(t *testing.T) {

	assert.NotEmpty(t, NewConfiguration().PublisherURL)
}
