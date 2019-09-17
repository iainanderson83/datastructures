package stack

import "testing"

func TestArrStack(t *testing.T) {
	tests := map[string]struct {
		in       string
		balanced bool
	}{
		"SingleType": {
			"(()())",
			true,
		},
		"Unbalanced": {
			"(((())(",
			false,
		},
		"Unbalanced2": {
			"((())))",
			false,
		},
		"MultipleTypes": {
			"({[]})",
			true,
		},
		"MultipleTypesUnbalanced": {
			"((({}]))",
			false,
		},
		"Empty": {
			"",
			true,
		},
		"UnknownChars": {
			"this (is a) sentence {kinda}",
			true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var s ArrStack
			for _, c := range test.in {
				switch c {
				case '(', '{', '[':
					s.Push(c)
				case ')', '}', ']':
					v := s.Pop()
					if v == nil {
						if !test.balanced {
							return
						}
						t.Fatal("expected balanced")
					}

					switch v {
					case ')':
						if c != '(' && test.balanced {
							t.Fatal("expected balanced")
						}
					case '}':
						if c != '{' && test.balanced {
							t.Fatal("expected balanced")
						}
					case ']':
						if c != '[' && test.balanced {
							t.Fatal("expected balanced")
						}
					default:
					}
				default:
				}
			}

			if s.Peek() != nil && test.balanced {
				t.Fatal("expected balanced")
			}
		})
	}
}
