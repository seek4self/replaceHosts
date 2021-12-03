package main

import (
	"reflect"
	"testing"
)

func Test_commentDomain(t *testing.T) {
	tests := []struct {
		name string
		line string
		want []byte
	}{
		// TODO: Add test cases.
		{"1", "140.82.113.21                central.github.com", []byte("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := commentDomain(tt.line); !reflect.DeepEqual(got, tt.want) {
				t.Log("out", string(got))
				t.Errorf("commentDomain() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
