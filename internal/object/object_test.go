package object

import (
	"encoding/hex"
	"errors"
	"reflect"
	"testing"

	"github.com/JunNishimura/Goit/internal/sha"
)

func TestNewObject(t *testing.T) {
	type args struct {
		objType Type
		data    []byte
	}
	type test struct {
		name    string
		args    args
		want    *Object
		wantErr error
	}
	tests := []*test{
		func() *test {
			hash, _ := hex.DecodeString("1856e9be02756984c385482a07e42f42efd5d2f3")

			return &test{
				name: "success",
				args: args{
					objType: BlobObject,
					data:    []byte("Hello, World"),
				},
				want: &Object{
					Type: BlobObject,
					Hash: sha.SHA1(hash),
					Size: len("Hello, World"),
					Data: []byte("Hello, World"),
				},
				wantErr: nil,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewObject(tt.args.objType, tt.args.data)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}
