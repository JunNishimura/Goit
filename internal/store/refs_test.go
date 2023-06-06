package store

import (
	"errors"
	"os"
	"path/filepath"
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

func TestLoadHash(t *testing.T) {
	type fields struct {
		name string
		hash sha.SHA1
	}
	tests := []struct {
		name    string
		fields  fields
		want    sha.SHA1
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				name: "main",
				hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
			},
			want:    sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			// .goit initialization
			goitDir := filepath.Join(tmpDir, ".goit")
			if err := os.Mkdir(goitDir, os.ModePerm); err != nil {
				t.Logf("%v: %s", err, goitDir)
			}
			// make .goit/refs directory
			refsDir := filepath.Join(goitDir, "refs")
			if err := os.Mkdir(refsDir, os.ModePerm); err != nil {
				t.Logf("%v: %s", err, refsDir)
			}
			// make .goit/refs/heads directory
			headsDir := filepath.Join(refsDir, "heads")
			if err := os.Mkdir(headsDir, os.ModePerm); err != nil {
				t.Logf("%v: %s", err, headsDir)
			}
			// make main branch
			mainBranchPath := filepath.Join(headsDir, tt.fields.name)
			f, err := os.Create(mainBranchPath)
			if err != nil {
				t.Logf("%v: %s", err, mainBranchPath)
			}
			if _, err := f.WriteString(tt.fields.hash.String()); err != nil {
				t.Log(err)
			}
			f.Close()

			b := newBranch(tt.fields.name, nil)
			if err := b.loadHash(goitDir); !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if b.hash.String() != tt.want.String() {
				t.Errorf("got = %s, want = %s", b.hash, tt.want)
			}
		})
	}
}
