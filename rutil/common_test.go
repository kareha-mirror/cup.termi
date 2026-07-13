package rutil

import (
	"testing"
)

func compareStringSlices(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestFilterByPrefix(t *testing.T) {
	tests := []struct {
		name     string
		list     []string
		prefix   string
		wantList []string
	}{
		{
			"typical",
			[]string{"hello", "test1", "test2"},
			"test",
			[]string{"test1", "test2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotList := FilterByPrefix(tt.list, tt.prefix)
			if !compareStringSlices(gotList, tt.wantList) {
				t.Errorf(
					"FilterByPrefix(%v, \"%s\") = %v; wanted %v",
					tt.list, tt.prefix, gotList, tt.wantList,
				)
			}
		})
	}
}

func TestCommonPrefix(t *testing.T) {
	tests := []struct {
		name       string
		list       []string
		wantPrefix string
	}{
		{
			"typical",
			[]string{"test1", "test2"},
			"test",
		},
		{
			"none",
			[]string{"hello", "world"},
			"",
		},
		{
			"empty",
			[]string{},
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrefix := CommonPrefix(tt.list)
			if gotPrefix != tt.wantPrefix {
				t.Errorf(
					"CommonPrefix(%v) = \"%s\"; wanted %s",
					tt.list, gotPrefix, tt.wantPrefix,
				)
			}
		})
	}
}
