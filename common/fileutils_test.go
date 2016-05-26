package common

import (
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {

	const VAR1 string = "COREOS_PRIVATE_IPV4=10.60.4.229"
	const VAR2 string = "a=b"
	d1 := []byte(VAR1 + "\n" + VAR2)
	err := ioutil.WriteFile("test.env", d1, 0644)
	if err != nil {
		panic(err)
	}
	var envVars []string
	count := 0
	ReadFile("test.env", func(line string) {
		count++
		envVars = append(envVars, line)
	})
	assert.Equal(t, count, 2)
	assert.Equal(t, envVars[0], VAR1)
	assert.Equal(t, envVars[1], VAR2)
}
