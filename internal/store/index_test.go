package store

import (
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/JunNishimura/Goit/internal/object"
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
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			hash = sha.SHA1(hash)
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
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			hash = sha.SHA1(hash)
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
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			hash = sha.SHA1(hash)
			path := []byte("cmd/main.go")

			return &test{
				name: "not update",
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

func TestGetEntry(t *testing.T) {
	type fields struct {
		hash sha.SHA1
		path []byte
	}
	type args struct {
		path []byte
	}
	type test struct {
		name      string
		fields    fields
		args      args
		wantPos   int
		wantEntry *Entry
		wantFound bool
	}
	tests := []*test{
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			hash = sha.SHA1(hash)
			path := []byte("cmd/main.go")

			entry := NewEntry(hash, path)

			return &test{
				name: "found",
				fields: fields{
					hash: hash,
					path: path,
				},
				args: args{
					path: path,
				},
				wantPos:   0,
				wantEntry: entry,
				wantFound: true,
			}
		}(),
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			hash = sha.SHA1(hash)
			path := []byte("cmd/main.go")

			return &test{
				name: "not found",
				fields: fields{
					hash: hash,
					path: path,
				},
				args: args{
					path: []byte("notfound.txt"),
				},
				wantPos:   -1,
				wantEntry: nil,
				wantFound: false,
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

			_, err = index.Update(goitDir, tt.fields.hash, tt.fields.path)
			if err != nil {
				t.Log(err)
			}

			pos, entry, isFound := index.GetEntry(tt.args.path)
			if pos != tt.wantPos {
				t.Errorf("got = %d, want = %d", pos, tt.wantPos)
			}
			if !reflect.DeepEqual(entry, tt.wantEntry) {
				t.Errorf("got = %v, want = %v", entry, tt.wantEntry)
			}
			if isFound != tt.wantFound {
				t.Errorf("got = %v, want = %v", isFound, tt.wantFound)
			}
		})
	}
}

func TestGetEntriesByDirectory(t *testing.T) {
	type args struct {
		dirName string
	}
	type fields struct {
		entries []*Entry
	}
	type test struct {
		name   string
		args   args
		fields fields
		want   []*Entry
	}
	tests := []*test{
		func() *test {
			return &test{
				name: "success: empty",
				args: args{
					dirName: "sample",
				},
				fields: fields{
					entries: []*Entry{},
				},
				want: nil,
			}
		}(),
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			hash = sha.SHA1(hash)
			return &test{
				name: "success: directory",
				args: args{
					dirName: "dir",
				},
				fields: fields{
					entries: []*Entry{
						NewEntry(hash, []byte("dir/dir2/test.txt")),
						NewEntry(hash, []byte("dir/test.txt")),
						NewEntry(hash, []byte("test.txt")),
					},
				},
				want: []*Entry{
					NewEntry(hash, []byte("dir/dir2/test.txt")),
					NewEntry(hash, []byte("dir/test.txt")),
				},
			}
		}(),
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			hash = sha.SHA1(hash)
			return &test{
				name: "true: sub directory",
				args: args{
					dirName: "dir/dir2",
				},
				fields: fields{
					entries: []*Entry{
						NewEntry(hash, []byte("dir/dir2/test.txt")),
						NewEntry(hash, []byte("dir/test.txt")),
						NewEntry(hash, []byte("test.txt")),
					},
				},
				want: []*Entry{
					NewEntry(hash, []byte("dir/dir2/test.txt")),
				},
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

			index, err := NewIndex(goitDir)
			if err != nil {
				t.Log(err)
			}

			for _, entry := range tt.fields.entries {
				_, err = index.Update(goitDir, entry.Hash, entry.Path)
				if err != nil {
					t.Log(err)
				}
			}

			got := index.GetEntriesByDirectory(tt.args.dirName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestIsRegisteredAsDirectory(t *testing.T) {
	type args struct {
		dirName string
	}
	type fields struct {
		entries []*Entry
	}
	type test struct {
		name   string
		args   args
		fields fields
		want   bool
	}
	tests := []*test{
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			hash = sha.SHA1(hash)
			return &test{
				name: "true",
				args: args{
					dirName: "dir",
				},
				fields: fields{
					entries: []*Entry{
						NewEntry(hash, []byte("dir/dir2/test.txt")),
						NewEntry(hash, []byte("dir/test.txt")),
						NewEntry(hash, []byte("test.txt")),
					},
				},
				want: true,
			}
		}(),
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			hash = sha.SHA1(hash)
			return &test{
				name: "true: sub directory",
				args: args{
					dirName: "dir/dir2",
				},
				fields: fields{
					entries: []*Entry{
						NewEntry(hash, []byte("dir/dir2/test.txt")),
						NewEntry(hash, []byte("dir/test.txt")),
						NewEntry(hash, []byte("test.txt")),
					},
				},
				want: true,
			}
		}(),
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			hash = sha.SHA1(hash)
			return &test{
				name: "false",
				args: args{
					dirName: "sample",
				},
				fields: fields{
					entries: []*Entry{
						NewEntry(hash, []byte("dir/dir2/test.txt")),
						NewEntry(hash, []byte("dir/test.txt")),
						NewEntry(hash, []byte("test.txt")),
					},
				},
				want: false,
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

			index, err := NewIndex(goitDir)
			if err != nil {
				t.Log(err)
			}

			for _, entry := range tt.fields.entries {
				_, err = index.Update(goitDir, entry.Hash, entry.Path)
				if err != nil {
					t.Log(err)
				}
			}

			got := index.IsRegisteredAsDirectory(tt.args.dirName)
			if got != tt.want {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestDeleteEntry(t *testing.T) {
	type fields struct {
		hash sha.SHA1
		path []byte
	}
	type args struct {
		path []byte
	}
	type test struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	tests := []*test{
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			hash = sha.SHA1(hash)
			path := []byte("cmd/main.go")

			return &test{
				name: "success",
				fields: fields{
					hash: hash,
					path: path,
				},
				args: args{
					path: path,
				},
				wantErr: false,
			}
		}(),
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			hash = sha.SHA1(hash)
			path := []byte("cmd/main.go")

			return &test{
				name: "not found",
				fields: fields{
					hash: hash,
					path: path,
				},
				args: args{
					path: []byte("not_exist.txt"),
				},
				wantErr: true,
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

			_, err = index.Update(goitDir, tt.fields.hash, tt.fields.path)
			if err != nil {
				t.Log(err)
			}

			if err := index.DeleteEntry(goitDir, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetEntriesFromTree(t *testing.T) {
	type args struct {
		rootName string
		nodes    []*object.Node
	}
	type test struct {
		name    string
		args    args
		want    []*Entry
		wantErr bool
	}
	tests := []*test{
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")

			return &test{
				name: "success",
				args: args{
					rootName: "",
					nodes: []*object.Node{
						{
							Hash:     hash,
							Name:     "a.txt",
							Children: []*object.Node{},
						},
						{
							Hash: hash,
							Name: "b",
							Children: []*object.Node{
								{
									Hash:     hash,
									Name:     "a.txt",
									Children: []*object.Node{},
								},
								{
									Hash:     hash,
									Name:     "b.txt",
									Children: []*object.Node{},
								},
							},
						},
						{
							Hash:     hash,
							Name:     "c.txt",
							Children: []*object.Node{},
						},
					},
				},
				want: []*Entry{
					{
						Hash:       hash,
						NameLength: uint16(len("a.txt")),
						Path:       []byte("a.txt"),
					},
					{
						Hash:       hash,
						NameLength: uint16(len("b/a.txt")),
						Path:       []byte("b/a.txt"),
					},
					{
						Hash:       hash,
						NameLength: uint16(len("b/b.txt")),
						Path:       []byte("b/b.txt"),
					},
					{
						Hash:       hash,
						NameLength: uint16(len("c.txt")),
						Path:       []byte("c.txt"),
					},
				},
				wantErr: false,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEntriesFromTree(tt.args.rootName, tt.args.nodes)
			if (err != nil) != tt.wantErr {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}
