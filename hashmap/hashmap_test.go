package hashmap

import (
	"math/rand"
	"testing"
)

var (
	redistributionTuples = []tuple{}
)

func TestMain(m *testing.M) {
	rand.Seed(42)

	ma := make(map[string]struct{})

	var i int
	for {
		key := wordList[rand.Intn(len(wordList))]
		if _, ok := ma[key]; ok {
			continue
		}
		ma[key] = struct{}{}
		tpl := tuple{key, wordList[rand.Intn(len(wordList))]}
		redistributionTuples = append(redistributionTuples, tpl)

		i++
		if i > 50 {
			break
		}
	}

	m.Run()
}

func TestMap(t *testing.T) {
	tests := map[string]struct {
		debug   bool
		adds    []tuple
		lookups []tuple
	}{
		"SingleValue": {
			false,
			[]tuple{{"hello", "world"}},
			[]tuple{{"hello", "world"}},
		},
		"OverwriteValue": {
			false,
			[]tuple{{"hello", "world"}, {"hello", "foo"}},
			[]tuple{{"hello", "foo"}},
		},
		"MultipleValues": {
			false,
			[]tuple{{"hello", "world"}, {"foo", "bar"}, {"baz", "bubbles"}, {"hello", "foo"}},
			[]tuple{{"foo", "bar"}, {"baz", "bubbles"}, {"hello", "foo"}},
		},
		"Redistribute": {
			false,
			redistributionTuples,
			redistributionTuples,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			debug = test.debug

			m := NewHashmap()
			for _, tpl := range test.adds {
				m.Add(tpl.key, tpl.value)
			}

			for _, tpl := range test.lookups {
				v := m.Lookup(tpl.key)
				if v != tpl.value {
					t.Fatalf("%s: expected '%v', got '%v'", tpl.key, tpl.value, v)
				}
			}
		})
	}
}
