package utils

import (
	"testing"
)

func TestFilesize_String(t *testing.T) {
	tests := []struct {
		name string
		size Filesize
		want string
	}{
		{
			name: "0 B",
			size: 0,
			want: "0 B",
		},
		{
			name: "100 B",
			size: 100,
			want: "100 B",
		},
		{
			name: "100 KB",
			size: 100 << 10,
			want: "100.0 KB",
		},
		{
			name: "100 MB",
			size: 100 << 20,
			want: "100.0 MB",
		},
		{
			name: "100 GB",
			size: 100 << 30,
			want: "100.0 GB",
		},
		{
			name: "1.0 TB",
			size: 1001 << 30,
			want: "1.0 TB",
		},
		{
			name: "1001.0 TB",
			size: 1001 << 40,
			want: "1001.0 TB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.size.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
