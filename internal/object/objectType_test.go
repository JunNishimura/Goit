package object

import (
	"errors"
	"testing"
)

func TestNewType(t *testing.T) {
	type args struct {
		typeString string
	}
	tests := []struct {
		name    string
		args    args
		want    Type
		wantErr error
	}{
		{
			name: "blob",
			args: args{
				typeString: "blob",
			},
			want:    BlobObject,
			wantErr: nil,
		},
		{
			name: "tree",
			args: args{
				typeString: "tree",
			},
			want:    TreeObject,
			wantErr: nil,
		},
		{
			name: "commit",
			args: args{
				typeString: "commit",
			},
			want:    CommitObject,
			wantErr: nil,
		},
		{
			name: "tag",
			args: args{
				typeString: "tag",
			},
			want:    TagObject,
			wantErr: nil,
		},
		{
			name: "undefined",
			args: args{
				typeString: "xxx",
			},
			want:    UndefinedObject,
			wantErr: ErrInvalidObject,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewType(tt.args.typeString)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got = %s, want = %s", got, tt.want)
			}
		})
	}
}
