package store

import (
	"reflect"
	"testing"

	"github.com/JunNishimura/Goit/internal/sha"
)

func TestNewEntry(t *testing.T) {
	type args struct {
		hash       sha.SHA1
		nameLength uint16
		path       []byte
	}
	type test struct {
		name string
		args args
		want *Entry
	}
	tests := []*test{
		func() *test {
			hash := sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335"))
			path := []byte("cmd/main.go")
			nameLength := uint16(len(path))

			return &test{
				name: "success",
				args: args{
					hash:       hash,
					nameLength: nameLength,
					path:       path,
				},
				want: &Entry{
					Hash:       hash,
					NameLength: nameLength,
					Path:       path,
				},
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEntry(tt.args.hash, tt.args.path)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}
