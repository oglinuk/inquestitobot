package main

import (
	//"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	dbi = NewDBInstance()
)

func TestNewDBInstance(t *testing.T) {
	assert.NotNil(t, dbi.DB)
}

func TestInsert(t *testing.T) {
	id, err := dbi.Insert(NewDocument("Test", "testing.com", "this is a test"))
	assert.Nil(t, err)
	assert.NotEqual(t, "", id)

	t.Cleanup(func() {
		//os.Remove(dbname)
	})
}
