package object

import (
	"encoding/hex"
	"errors"
	"reflect"
	"testing"
)

func TestNewGetPaths(t *testing.T) {
	type fields struct {
		node Node
	}
	type test struct {
		name   string
		fields fields
		want   []string
	}
	tests := []*test{
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")

			return &test{
				name: "success: empty",
				fields: fields{
					node: Node{
						Hash:     hash,
						Name:     "test.txt",
						Children: nil,
					},
				},
				want: []string{"test.txt"},
			}
		}(),
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")

			return &test{
				name: "success: multiple nodes",
				fields: fields{
					node: Node{
						Hash: hash,
						Name: "dir",
						Children: []*Node{
							{
								Hash: hash,
								Name: "dir2",
								Children: []*Node{
									{
										Hash:     hash,
										Name:     "test.txt",
										Children: nil,
									},
									{
										Hash:     hash,
										Name:     "test2.txt",
										Children: nil,
									},
								},
							},
							{
								Hash:     hash,
								Name:     "test.txt",
								Children: nil,
							},
							{
								Hash:     hash,
								Name:     "test2.txt",
								Children: nil,
							},
						},
					},
				},
				want: []string{
					"dir/dir2/test.txt",
					"dir/dir2/test2.txt",
					"dir/test.txt",
					"dir/test2.txt",
				},
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.node.GetPaths()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}

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
		func() *test {
			object, _ := NewObject(BlobObject, []byte("blob 12\x00Hello, World"))

			return &test{
				name: "fail: invalid tree object",
				args: args{
					object: object,
				},
				want:    nil,
				wantErr: ErrInvalidTreeObject,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTree("", tt.args.object)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if tt.want == nil {
				return
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
				if !gotChild.Hash.Compare(wantChild.Hash) {
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

func TestString(t *testing.T) {
	type args struct {
		object *Object
	}
	type test struct {
		name    string
		args    args
		want    string
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

			return &test{
				name: "success",
				args: args{
					object: object,
				},
				want:    "100644 blob 1856e9be02756984c385482a07e42f42efd5d2f3	test.txt",
				wantErr: nil,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, _ := NewTree("", tt.args.object)
			if tree.String() != tt.want {
				t.Errorf("got = %s, want = %s", tree, tt.want)
			}
		})
	}
}
