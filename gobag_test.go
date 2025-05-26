package gobag

import (
	"errors"
	"reflect"
	"testing"
)

func TestFields(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		err      error
	}{
		// Normal cases
		{
			name:     "simple comma-separated",
			input:    "a,b,c",
			expected: []string{"a", "b", "c"},
			err:      nil,
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
			err:      nil,
		},
		{
			name:     "single field",
			input:    "hello",
			expected: []string{"hello"},
			err:      nil,
		},
		{
			name:     "spaces around commas",
			input:    " a , b , c ",
			expected: []string{" a ", " b ", " c "},
			err:      nil,
		},

		// Quoted strings
		{
			name:     "double quotes with comma",
			input:    `a,"b,c",d`,
			expected: []string{"a", `"b,c"`, "d"},
			err:      nil,
		},
		{
			name:     "single quotes with comma",
			input:    `a,'b,c',d`,
			expected: []string{"a", `'b,c'`, "d"},
			err:      nil,
		},
		{
			name:     "nested quotes",
			input:    `a,"b,'c,d',e",f`,
			expected: []string{"a", `"b,'c,d',e"`, "f"},
			err:      nil,
		},

		// Parentheses
		{
			name:     "balanced parentheses",
			input:    "a,(b,c),d",
			expected: []string{"a", "(b,c)", "d"},
			err:      nil,
		},
		{
			name:     "nested parentheses",
			input:    "a,(b,(c,d),e),f",
			expected: []string{"a", "(b,(c,d),e)", "f"},
			err:      nil,
		},

		// Escape sequences
		{
			name:     "escaped comma",
			input:    `a\,b,c`,
			expected: []string{"a,b", "c"},
			err:      nil,
		},
		{
			name:     "escaped quotes",
			input:    `a,\"b\",c`,
			expected: []string{"a", `"b"`, "c"},
			err:      nil,
		},
		{
			name:     "escaped backslash",
			input:    `a\\b,c`,
			expected: []string{`a\b`, "c"},
			err:      nil,
		},

		// Unicode
		{
			name:     "unicode characters",
			input:    "α,β,γ",
			expected: []string{"α", "β", "γ"},
			err:      nil,
		},

		// Error cases
		{
			name:     "dangling escape",
			input:    `a,b\`,
			expected: nil,
			err:      errors.New("dangling escape character at end of string"),
		},
		{
			name:     "unbalanced single quote",
			input:    `a,'b,c`,
			expected: nil,
			err:      errors.New("unbalanced single quote in string"),
		},
		{
			name:     "unbalanced double quote",
			input:    `a,"b,c`,
			expected: nil,
			err:      errors.New("unbalanced double quote in string"),
		},
		{
			name:     "unbalanced parentheses",
			input:    `a,(b,c`,
			expected: nil,
			err:      errors.New("unbalanced parentheses in string"),
		},
		{
			name:     "too many closing parentheses",
			input:    `a,b)c`,
			expected: nil,
			err:      errors.New("too many closing parentheses"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Fields(tt.input, ',')
			if tt.err != nil {
				if err == nil || err.Error() != tt.err.Error() {
					t.Errorf("commaFields(%q) error = %v, want %v", tt.input, err, tt.err)
				}
			} else if err != nil {
				t.Errorf("commaFields(%q) error = %v, want nil", tt.input, err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("commaFields(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestUnquoteString(t *testing.T) {
	positiveTests := []struct {
		str     string
		expects string
	}{
		{str: `"foo ","bar"`, expects: "foo ,bar"},
		{str: `"foo\"","bar"`, expects: `foo",bar`},
		{str: `"foo\"\n","bar"`, expects: `foo"\n,bar`},
		{str: `"foo\"\n","bar\\\\"`, expects: `foo"\n,bar\\`},
	}

	for i, test := range positiveTests {
		s, err := UnquoteString(test.str)
		if err != nil {
			t.Fatalf("positive test %d [%s]: %v", i, test.str, err)
		}
		if s != test.expects {
			t.Fatalf("positive test %d [%s]: expected [%s] got [%s]", i+1, test.str, test.expects, s)
		}
	}

	negativeTests := []string{
		`\""foo"\"`,
		`"foo""`,
		`"foo\"`,
		`"foo","bar`,
		`"foo","bar\"`,
	}

	for i, test := range negativeTests {
		s, err := UnquoteString(test)
		if err == nil {
			t.Fatalf("negative test %d [%s]: ought to fail, but didn't, result [%s]", i+1, test, s)
		}
	}

}

func TestDeduplicate(t *testing.T) {
	tests := []struct {
		input    []int
		expected []int
	}{
		{[]int{1, 2, 2, 3, 1}, []int{1, 2, 3}},
		{[]int{}, []int{}},
		{[]int{1}, []int{1}},
		{nil, []int{}},
	}
	for _, tt := range tests {
		result := Deduplicate(tt.input)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("deduplicate(%v) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}
