package binary

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestReadNullTerminatedString(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "normal",
			args: args{
				r: bytes.NewBufferString("read until here\x00not read here"),
			},
			want:    "read until here",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadNullTerminatedString(tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got = %s, want = %s", got, tt.want)
			}
		})
	}
}
