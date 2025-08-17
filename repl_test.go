package main

import (
	"testing"
	"time"

	"github.com/heavysider/pokedexcli/internal/pokecache"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		// add more cases here
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("slice sizes don't match :(")
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("words don't match :(")
			}
		}
	}
}

func TestCache(t *testing.T) {
	cache := pokecache.NewCache(1 * time.Second)
	cache.Add("test", []byte{})
	_, ok := cache.Get("test")
	if !ok {
		t.Errorf("no items in cache :(")
	}
	time.Sleep(2 * time.Second)
	_, ok = cache.Get("test")
	if ok {
		t.Errorf("cache didn't clear in time :(")
	}
}
