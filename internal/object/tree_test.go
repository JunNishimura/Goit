package object

import (
	"encoding/hex"
	"errors"
	"reflect"
	"testing"
)

func TestNewTree(t *testing.T) {
	type args struct {
		object *Object
	}
	type test struct {
		name    string
		args    args
		want    *Tree
		wantErr error
	}
	tests := []*test{
		func() *test {
			hash, _ := hex.DecodeString("1856e9be02756984c385482a07e42f42efd5d2f3")

			var data []byte
			data = append(data, []byte("100644 test.txt")...)
			data = append(data, 0x00)
			data = append(data, hash...)

			object, _ := NewObject(TreeObject, data)

			node := &Node{
				Hash:     hash,
				Name:     "test.txt",
				Children: []*Node{},
			}

			return &test{
				name: "success",
				args: args{
					object: object,
				},
				want: &Tree{
					object:   object,
					Children: []*Node{node},
				},
				wantErr: nil,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTree("", tt.args.object)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got.object, tt.want.object) {
				t.Errorf("got = %v, want = %v", got.object, tt.want.object)
			}
			if len(got.Children) != len(tt.want.Children) {
				t.Errorf("got = %v, want = %v", len(got.Children), len(tt.want.Children))
			}
			for i := 0; i < len(got.Children); i++ {
				gotChild := got.Children[i]
				wantChild := tt.want.Children[i]
				if gotChild.Hash.String() != wantChild.Hash.String() {
					t.Errorf("got = %v, want = %v", gotChild.Hash.String(), wantChild.Hash.String())
				}
				if gotChild.Name != wantChild.Name {
					t.Errorf("got = %v, want = %v", gotChild.Name, wantChild.Name)
				}
				if len(gotChild.Children) != len(wantChild.Children) {
					t.Errorf("got = %v, want = %v", len(gotChild.Children), len(wantChild.Children))
				}
			}
		})
	}
}