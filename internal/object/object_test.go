package object

import (
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
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

func TestGetObject(t *testing.T) {
	type args struct {
		hash sha.SHA1
	}
	type test struct {
		name      string
		args      args
		objType   Type
		data      []byte
		isCreated bool
		want      *Object
		wantErr   error
	}
	tests := []*test{
		func() *test {
			hash, _ := hex.DecodeString("1856e9be02756984c385482a07e42f42efd5d2f3")

			return &test{
				name: "success",
				args: args{
					hash: hash,
				},
				objType:   BlobObject,
				data:      []byte("Hello, World"),
				isCreated: true,
				want: &Object{
					Type: BlobObject,
					Hash: sha.SHA1(hash),
					Size: len("Hello, World"),
					Data: []byte("Hello, World"),
				},
				wantErr: nil,
			}
		}(),
		func() *test {
			hash, _ := hex.DecodeString("1856e9be02756984c385482a07e42f42efd5d2f3")

			return &test{
				name: "fail: no object",
				args: args{
					hash: hash,
				},
				objType:   BlobObject,
				data:      []byte(""),
				isCreated: false,
				want:      nil,
				wantErr:   ErrIOHandling,
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

			// make object
			if tt.isCreated {
				obj, _ := NewObject(tt.objType, tt.data)
				_ = obj.Write(goitDir)
			}

			got, err := GetObject(goitDir, tt.args.hash)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}
