package store

import (
	"errors"
	"os"
	"path/filepath"
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

func TestNewIndex(t *testing.T) {
	type test struct {
		name    string
		want    *Index
		wantErr error
	}
	tests := []*test{
		func() *test {
			return &test{
				name: "not exist index",
				want: &Index{
					Header: Header{
						Signature: [4]byte{'D', 'I', 'R', 'C'},
						Version:   uint32(1),
						EntryNum:  uint32(0),
					},
				},
				wantErr: nil,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			// .goit initialization
			goitDir := filepath.Join(tmpDir, ".goit")
			if err := os.Mkdir(goitDir, os.ModePerm); err != nil {
				t.Logf("%v: %s", err, goitDir)
			}
			// make .goit/config file
			configFile := filepath.Join(goitDir, "config")
			f, err := os.Create(configFile)
			if err != nil {
				t.Logf("%v: %s", err, configFile)
			}
			f.Close()
			// make .goit/HEAD file and write main branch
			headFile := filepath.Join(goitDir, "HEAD")
			f, err = os.Create(headFile)
			if err != nil {
				t.Logf("%v: %s", err, headFile)
			}
			// set 'main' as default branch
			if _, err := f.WriteString("ref: refs/heads/main"); err != nil {
				t.Logf("%v: %s", err, headFile)
			}
			f.Close()
			// make .goit/objects directory
			objectsDir := filepath.Join(goitDir, "objects")
			if err := os.Mkdir(objectsDir, os.ModePerm); err != nil {
				t.Logf("%v: %s", err, objectsDir)
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
			// make .goit/refs/tags directory
			tagsDir := filepath.Join(refsDir, "tags")
			if err := os.Mkdir(tagsDir, os.ModePerm); err != nil {
				t.Logf("%v: %s", err, tagsDir)
			}

			index, err := NewIndex(goitDir)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(index, tt.want) {
				t.Errorf("got = %v, want = %v", index, tt.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		hash sha.SHA1
		path []byte
	}
	type test struct {
		name     string
		args     args
		isRepeat bool
		want     bool
		wantErr  error
	}
	tests := []*test{
		func() *test {
			hash := sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335"))
			path := []byte("cmd/main.go")

			return &test{
				name: "success",
				args: args{
					hash: hash,
					path: path,
				},
				isRepeat: false,
				want:     true,
				wantErr:  nil,
			}
		}(),
		func() *test {
			hash := sha.SHA1([]byte("87f3c49bccf2597484ece08746d3ee5defaba335"))
			path := []byte("cmd/main.go")

			return &test{
				name: "success",
				args: args{
					hash: hash,
					path: path,
				},
				isRepeat: true,
				want:     false,
				wantErr:  nil,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			// .goit initialization
			goitDir := filepath.Join(tmpDir, ".goit")
			if err := os.Mkdir(goitDir, os.ModePerm); err != nil {
				t.Logf("%v: %s", err, goitDir)
			}
			// make .goit/config file
			configFile := filepath.Join(goitDir, "config")
			f, err := os.Create(configFile)
			if err != nil {
				t.Logf("%v: %s", err, configFile)
			}
			f.Close()
			// make .goit/HEAD file and write main branch
			headFile := filepath.Join(goitDir, "HEAD")
			f, err = os.Create(headFile)
			if err != nil {
				t.Logf("%v: %s", err, headFile)
			}
			// set 'main' as default branch
			if _, err := f.WriteString("ref: refs/heads/main"); err != nil {
				t.Logf("%v: %s", err, headFile)
			}
			f.Close()
			// make .goit/objects directory
			objectsDir := filepath.Join(goitDir, "objects")
			if err := os.Mkdir(objectsDir, os.ModePerm); err != nil {
				t.Logf("%v: %s", err, objectsDir)
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
			// make .goit/refs/tags directory
			tagsDir := filepath.Join(refsDir, "tags")
			if err := os.Mkdir(tagsDir, os.ModePerm); err != nil {
				t.Logf("%v: %s", err, tagsDir)
			}

			index, err := NewIndex(goitDir)
			if err != nil {
				t.Log(err)
			}

			isUpdated, err := index.Update(goitDir, tt.args.hash, tt.args.path)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if tt.isRepeat {
				isUpdated, err = index.Update(goitDir, tt.args.hash, tt.args.path)
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("got = %v, want = %v", err, tt.wantErr)
				}
			}
			if isUpdated != tt.want {
				t.Errorf("got = %v, want = %v", isUpdated, tt.want)
			}
		})
	}
}
