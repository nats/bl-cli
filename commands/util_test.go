package commands

import (
	"testing"
)

func Test_looksLikeUUID(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{
			name: "UUIDv1",
			arg:  "a235f190-01a4-11ea-af17-4989e8155574",
			want: true,
		},
		{
			name: "UUIDv4",
			arg:  "21464f77-42fc-4a32-9aa4-a101843e94c0",
			want: true,
		},
		{
			name: "my-cluster",
			arg:  "my-cluster",
			want: false,
		},
		{
			name: "f8291060b73f4fa7b60586fe51a1d862",
			arg:  "f8291060b73f4fa7b60586fe51a1d862",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := looksLikeUUID(tt.arg); got != tt.want {
				t.Errorf("looksLikeUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
