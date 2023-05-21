package binary

import (
	"bytes"
	"io"
	"testing"
)

func TestReadNullTerminatedString(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{
				r: bytes.NewBufferString("read until here\x00not read here"),
			},
			want: "read until here",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := ReadNullTerminatedString(tt.args.r); got != tt.want {
				t.Errorf("got = %s, want = %s", got, tt.want)
			}
		})
	}
}
