package store

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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

func TestNewRefs(t *testing.T) {
	type fields struct {
		name string
		hash sha.SHA1
	}
	tests := []struct {
		name    string
		fields  []*fields
		want    *Refs
		wantErr error
	}{
		{
			name:   "success: no heads",
			fields: nil,
			want: &Refs{
				Heads: make([]*branch, 0),
			},
			wantErr: nil,
		},
		{
			name: "success: some heads",
			fields: []*fields{
				{
					name: "main",
					hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
				},
				{
					name: "test",
					hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
				},
			},
			want: &Refs{
				Heads: []*branch{
					{
						Name: "main",
						hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
					},
					{
						Name: "test",
						hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
					},
				},
			},
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
			for _, field := range tt.fields {
				// make main branch
				branchPath := filepath.Join(headsDir, field.name)
				f, err := os.Create(branchPath)
				if err != nil {
					t.Logf("%v: %s", err, branchPath)
				}
				if _, err := f.WriteString(field.hash.String()); err != nil {
					t.Log(err)
				}
				f.Close()
			}

			got, err := NewRefs(goitDir)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestAddBranch(t *testing.T) {
	type args struct {
		newBranchName string
		newBranchHash sha.SHA1
	}
	type fields struct {
		name string
		hash sha.SHA1
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		want    *Refs
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				newBranchName: "main",
				newBranchHash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
			},
			fields: fields{
				name: "test",
				hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
			},
			want: &Refs{
				Heads: []*branch{
					{
						Name: "main",
						hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
					},
					{
						Name: "test",
						hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "failure",
			args: args{
				newBranchName: "main",
				newBranchHash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
			},
			fields: fields{
				name: "main",
				hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
			},
			want: &Refs{
				Heads: []*branch{
					{
						Name: "main",
						hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
					},
				},
			},
			wantErr: true,
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
			branchPath := filepath.Join(headsDir, tt.fields.name)
			f, err := os.Create(branchPath)
			if err != nil {
				t.Logf("%v: %s", err, branchPath)
			}
			if _, err := f.WriteString(tt.fields.hash.String()); err != nil {
				t.Log(err)
			}
			f.Close()

			r, err := NewRefs(goitDir)
			if err != nil {
				t.Log(err)
			}

			if err := r.AddBranch(goitDir, tt.args.newBranchName, tt.args.newBranchHash); (err != nil) != tt.wantErr {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(r, tt.want) {
				t.Errorf("got = %v, want = %v", r, tt.want)
			}
		})
	}
}

func TestGetBranchPos(t *testing.T) {
	type args struct {
		branchName string
	}
	type fields struct {
		name string
		hash sha.SHA1
	}
	tests := []struct {
		name       string
		args       args
		fieldsList []*fields
		want       int
	}{
		{
			name: "found existing branch",
			args: args{
				branchName: "main",
			},
			fieldsList: []*fields{
				{
					name: "main",
					hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
				},
				{
					name: "test",
					hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
				},
			},
			want: 0,
		},
		{
			name: "not found",
			args: args{
				branchName: "xxxx",
			},
			fieldsList: []*fields{
				{
					name: "main",
					hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
				},
				{
					name: "test",
					hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
				},
			},
			want: NewBranchFlag,
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
			for _, field := range tt.fieldsList {
				// make main branch
				branchPath := filepath.Join(headsDir, field.name)
				f, err := os.Create(branchPath)
				if err != nil {
					t.Logf("%v: %s", err, branchPath)
				}
				if _, err := f.WriteString(field.hash.String()); err != nil {
					t.Log(err)
				}
				f.Close()
			}

			r, err := NewRefs(goitDir)
			if err != nil {
				t.Log(err)
			}
			got := r.getBranchPos(tt.args.branchName)
			if got != tt.want {
				t.Errorf("got = %d, want = %d", got, tt.want)
			}
		})
	}
}

func TestDeleteBranch(t *testing.T) {
	type args struct {
		headBranchName   string
		deleteBranchName string
	}
	type fields struct {
		name string
		hash sha.SHA1
	}
	tests := []struct {
		name       string
		args       args
		fieldsList []*fields
		want       *Refs
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				headBranchName:   "main",
				deleteBranchName: "test",
			},
			fieldsList: []*fields{
				{
					name: "main",
					hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
				},
				{
					name: "test",
					hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
				},
			},
			want: &Refs{
				Heads: []*branch{
					{
						Name: "main",
						hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "failure",
			args: args{
				headBranchName:   "main",
				deleteBranchName: "main",
			},
			fieldsList: []*fields{
				{
					name: "main",
					hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
				},
				{
					name: "test",
					hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
				},
			},
			want: &Refs{
				Heads: []*branch{
					{
						Name: "main",
						hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
					},
					{
						Name: "test",
						hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
					},
				},
			},
			wantErr: true,
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
			for _, field := range tt.fieldsList {
				// make main branch
				branchPath := filepath.Join(headsDir, field.name)
				f, err := os.Create(branchPath)
				if err != nil {
					t.Logf("%v: %s", err, branchPath)
				}
				if _, err := f.WriteString(field.hash.String()); err != nil {
					t.Log(err)
				}
				f.Close()
			}

			r, err := NewRefs(goitDir)
			if err != nil {
				t.Log(err)
			}
			if err := r.DeleteBranch(goitDir, tt.args.headBranchName, tt.args.deleteBranchName); (err != nil) != tt.wantErr {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(r, tt.want) {
				t.Errorf("got = %v, want = %v", r, tt.want)
			}
		})
	}
}

func TestUpdateBranchHash(t *testing.T) {
	type args struct {
		branchName string
		newHash    sha.SHA1
	}
	type fields struct {
		name string
		hash sha.SHA1
	}
	tests := []struct {
		name       string
		args       args
		fieldsList []*fields
		want       *Refs
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				branchName: "main",
				newHash:    sha.SHA1([]byte(strings.Repeat("4", 40))),
			},
			fieldsList: []*fields{
				{
					name: "main",
					hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
				},
				{
					name: "test",
					hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
				},
			},
			want: &Refs{
				Heads: []*branch{
					{
						Name: "main",
						hash: sha.SHA1([]byte(strings.Repeat("4", 40))),
					},
					{
						Name: "test",
						hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "failure",
			args: args{
				branchName: "xxxx",
				newHash:    sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
			},
			fieldsList: []*fields{
				{
					name: "main",
					hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
				},
				{
					name: "test",
					hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
				},
			},
			want: &Refs{
				Heads: []*branch{
					{
						Name: "main",
						hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
					},
					{
						Name: "test",
						hash: sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335")),
					},
				},
			},
			wantErr: true,
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
			for _, field := range tt.fieldsList {
				// make main branch
				branchPath := filepath.Join(headsDir, field.name)
				f, err := os.Create(branchPath)
				if err != nil {
					t.Logf("%v: %s", err, branchPath)
				}
				if _, err := f.WriteString(field.hash.String()); err != nil {
					t.Log(err)
				}
				f.Close()
			}

			r, err := NewRefs(goitDir)
			if err != nil {
				t.Log(err)
			}
			if err := r.UpdateBranchHash(goitDir, tt.args.branchName, tt.args.newHash); (err != nil) != tt.wantErr {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(r, tt.want) {
				t.Errorf("got = %v, want = %v", r, tt.want)
			}
		})
	}
}
