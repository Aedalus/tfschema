package v1schema

import (
	"fmt"
	"io/ioutil"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestParseV1Schema(t *testing.T) {
	file, err := ioutil.ReadFile("./file-schema.json")

	assert.Nil(t, err)
	assert.NotNil(t, file)

	schema, err := ParseV1Schema(file)

	assert.Nil(t, err)
	assert.NotNil(t, schema)

	fmt.Printf("%+v", schema)
}
