package store

import (
	"reflect"
	"testing"

	"github.com/JunNishimura/Goit/internal/sha"
)

func TestNewBanch(t *testing.T) {
	type args struct {
		name string
		hash sha.SHA1
	}
	tests := []struct {
		name string
		args args
		want *branch
	}{
		{
			name: "success",
			args: args{
				name: "main",
				hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
			},
			want: &branch{
				Name: "main",
				hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newBranch(tt.args.name, tt.args.hash)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}
