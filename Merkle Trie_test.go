package main

import (
	"testing"
)

func TestHash_String(t *testing.T) {
	tests := []struct {
		name string
		hash Hash
		want string
	}{
		{
			name: "Test1",
			hash: NewHash([]byte{
				0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
				0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
				0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c,
				0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00,
			}),
			want: "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hash.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
