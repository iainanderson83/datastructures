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
			m := NewFNV1aHashmap()
			m2 := NewRuntimeHashmap()
			for _, tpl := range test.adds {
				m.Add(tpl.key, tpl.value)
				m2.Add(tpl.key, tpl.value)
			}

			for _, tpl := range test.lookups {
				v := m.Lookup(tpl.key)
				if v != tpl.value {
					t.Fatalf("%s: expected '%v', got '%v'", tpl.key, tpl.value, v)
				}

				v2 := m2.Lookup(tpl.key)
				if v2 != tpl.value {
					t.Fatalf("%s: expected '%v', got '%v'", tpl.key, tpl.value, v2)
				}
			}
		})
	}
}

func BenchmarkRuntimeHashmap(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		m := NewRuntimeHashmap()
		for _, tpl := range redistributionTuples {
			m.Add(tpl.key, tpl.value)
		}

		for _, tpl := range redistributionTuples {
			v := m.Lookup(tpl.key)
			_ = v
		}
	}
}

func BenchmarkFNV1aHashmap(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		m := NewFNV1aHashmap()
		for _, tpl := range redistributionTuples {
			m.Add(tpl.key, tpl.value)
		}

		for _, tpl := range redistributionTuples {
			v := m.Lookup(tpl.key)
			_ = v
		}
	}
}

func BenchmarkGoHashmap(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		m := make(map[string]interface{})
		for _, tpl := range redistributionTuples {
			m[tpl.key] = tpl.value
		}

		for _, tpl := range redistributionTuples {
			v := m[tpl.key]
			_ = v
		}
	}
}
