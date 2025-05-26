// Package gobag offers a collection of small, generic utility functions for slices, maps, string parsing, and more.
// It provides practical helpers like deduplication, key extraction, and balanced string splitting.
package gobag

import (
	"errors"
	"strings"
)

// Fields splits a string by the given separator rune, respecting
// balanced parentheses, quoted substrings, and escape sequences.
// Returns an error if quotes or parentheses are unbalanced.
func Fields(s string, sep rune) ([]string, error) {
	var sb strings.Builder
	fields := make([]string, 0)
	var balance int
	var inSingle, inDouble, isEscaped bool

	for _, r := range s {
		if isEscaped {
			sb.WriteRune(r)
			isEscaped = false
			continue
		}

		switch r {
		case '\\':
			isEscaped = true
			continue
		case sep:
			if balance == 0 && !inSingle && !inDouble {
				fields = append(fields, sb.String())
				sb.Reset()
				continue
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '(':
			if !inSingle && !inDouble {
				balance++
			}
		case ')':
			if !inSingle && !inDouble {
				balance--
			}
		}
		sb.WriteRune(r)
	}

	if isEscaped {
		return nil, errors.New("dangling escape character at end of string")
	}
	if balance < 0 {
		return nil, errors.New("too many closing parentheses")
	}
	if balance != 0 {
		return nil, errors.New("unbalanced parentheses in string")
	}
	if inSingle {
		return nil, errors.New("unbalanced single quote in string")
	}
	if inDouble {
		return nil, errors.New("unbalanced double quote in string")
	}

	if sb.Len() > 0 {
		fields = append(fields, sb.String())
	}

	return fields, nil
}

// UnquoteStrings unquote double quote strings in a string slice.
func UnquoteStrings(elems []string) ([]string, error) {
	unquoted := make([]string, 0, len(elems))
	for _, elem := range elems {
		s, err := UnquoteString(elem)
		if err != nil {
			return nil, err
		}
		unquoted = append(unquoted, s)
	}

	return unquoted, nil
}

// UnquoteString unquotes double quotes in a string.
func UnquoteString(s string) (string, error) {
	var sb strings.Builder

	inQuote := false
	escape := false
	for _, r := range s {
		switch {
		case escape:
			switch r {
			case '"', '\\':
				sb.WriteRune(r)
			default:
				sb.WriteRune('\\')
				sb.WriteRune(r)
			}
			escape = false
		case r == '\\':
			if !inQuote {
				return "", errors.New("escape character found outside a quote")
			}
			escape = true
		case r == '"':
			inQuote = !inQuote
		default:
			sb.WriteRune(r)
		}
	}

	if escape {
		return "", errors.New("dangling escape character at end of string")
	}
	if inQuote {
		return "", errors.New("unterminated double quote")
	}
	return sb.String(), nil
}

// Ternary returns v1 if the condition is true, otherwise it returns v2.
func Ternary[T any](cond bool, v1, v2 T) T {
	if cond {
		return v1
	}
	return v2
}

// Deduplicate returns a new slice with duplicates removed, preserving
// the order of first occurrence.
func Deduplicate[T comparable](s []T) []T {
	if len(s) == 0 {
		return []T{}
	}

	seen := make(map[T]struct{}, len(s))
	result := make([]T, 0, len(s))
	for _, e := range s {
		if _, ok := seen[e]; !ok {
			seen[e] = struct{}{}
			result = append(result, e)
		}
	}

	return result
}

// In reports whether the given element is present in the provided slice, using equality comparison.
func In[T comparable](s []T, e T) bool {
	for _, element := range s {
		if element == e {
			return true
		}
	}
	return false
}

// Keys returns a slice of keys from the given map.
// The order of keys is not guaranteed.
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
