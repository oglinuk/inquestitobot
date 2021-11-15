package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	expected := "098f6bcd4621d373cade4e832627b4f6"
	actual := hash("test")
	assert.Equal(t, expected, actual)
}

func TestIHash(t *testing.T) {
	expected := "8df1c3f47e4cf4bc4603ddd0eebcb4ac"
	actual := ihash(&Document{})
	assert.Equal(t, expected, actual)
}

func TestNewDocument(t *testing.T) {
	expected := &Document{
		"f75c2ec8e1253fb6f411875cebefcab9",
		"",
		"https://library.fourohfournotfound.com",
		`This is my library of (digital) books that are freely available
		under their respective permissive licenses. If you have books that
		you know of which aren’t here please get in touch!`,
		"834d3bdc08a880ee13290bff3ba411a0",
	}

	actual := NewDocument(
		"",
		"https://library.fourohfournotfound.com",
		`This is my library of (digital) books that are freely available
		under their respective permissive licenses. If you have books that
		you know of which aren’t here please get in touch!`)

	assert.Equal(t, expected, actual)
}
