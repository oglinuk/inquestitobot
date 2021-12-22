package main

import (
	"testing"
	"log"
	"os"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestNewProcessor(t *testing.T) {
	p := NewProcessor()
	assert.NotNil(t, p)
	assert.Nil(t, p.Doc)
	assert.NotNil(t, p.Ttype)
	assert.NotNil(t, p.T)

	t.Cleanup(func() {
		os.Remove(dbname)
	})
}

func TestIterate(t *testing.T) {
	p := NewProcessor()
	f, err := os.Open("./testdata/fourohfournotfound.html")
	assert.Nil(t, err)
	defer f.Close()

	p.Doc = html.NewTokenizer(f)
	p.Ttype = p.Doc.Next()
	token := p.Doc.Token()
	log.Println(token)
	p.iterate()
	log.Println(p.Doc.Token())
	assert.NotEqual(t, token, p.Doc.Token())
}

func TestGetMetadata(t *testing.T) {
	p := NewProcessor()
	assert.NotNil(t, p)

	testHTML := []string{
		"./testdata/fourohfournotfound.html",
		"./testdata/cc.html",
		"./testdata/chaos-theory.html",
	}

	expected := [][]string{
		{"Fourohfournotfound", "fourohfournotfound.com is my personal knowledge node!"},
		{"When we share, everyone wins - Creative Commons", "Help us build a vibrant, collaborative global commons"},
		{"Chaos theory - Wikipedia", "Chaos theory is a branch of mathematics focusing on the study of chaos — dynamical systems whose apparently random states of disorder and irregularities are actually governed by underlying patterns and deterministic laws that are highly sensitive to initial conditions.[1][2] Chaos theory is an interdisciplinary theory stating that, within the apparent randomness of chaotic complex systems, there are underlying patterns, interconnectedness, constant feedback loops, repetition, self-similarity, fractals, and self-organization.[3] The butterfly effect, an underlying principle of chaos, describes how a small change in one state of a deterministic nonlinear system can result in large differences in a later state (meaning that there is sensitive dependence on initial conditions).[4] A metaphor for this behavior is that a butterfly flapping its wings in Texas can cause a hurricane in China.[5]"},
		{"", ""},
	}

	for i, page := range testHTML {
		f, err := os.Open(page)
		assert.Nil(t, err)
		defer f.Close()

		e := expected[i]
		aTitle, aDesc := p.getMetadata(f)
		assert.Equal(t, e[0], aTitle)
		assert.Equal(t, e[1], aDesc)
	}

	t.Cleanup(func() {
		os.Remove(dbname)
	})
}

func TestFilter(t *testing.T) {
	p := NewProcessor()
	assert.NotNil(t, p)

	expected := []string{
		"https://fourohfournotfound.com",
		"https://en.wikipedia.org/wiki/Search_engine",
		"https://creativecommons.org",
	}

	actual := p.filter()
	assert.Equal(t, expected, actual)

	t.Cleanup(func() {
		os.Remove(dbname)
	})
}

/*
func TestProcess(t *testing.T) {
	p := NewProcessor()
	assert.NotNil(t, p)

	originalSeedsLen := len(p.Seeds)

	p.process()

	// TODO: Create an actual expectedProcessor struct as this is a
	// temporary solution to ensure its not the original seeds array
	assert.NotEqual(t, originalSeedsLen, len(p.Seeds))

	t.Cleanup(func() {
		os.Remove(dbname)
	})
}
*/
