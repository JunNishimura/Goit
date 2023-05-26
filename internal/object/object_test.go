package object

import (
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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

func TestReadHeader(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name     string
		args     args
		wantType Type
		wantSize int
		wantErr  error
	}{
		{
			name: "success",
			args: args{
				r: strings.NewReader("blob 12\x00Hello, World"),
			},
			wantType: BlobObject,
			wantSize: 12,
			wantErr:  nil,
		},
		{
			name: "fail: empty header",
			args: args{
				r: strings.NewReader(""),
			},
			wantType: UndefinedObject,
			wantSize: 0,
			wantErr:  ErrInvalidObject,
		},
		{
			name: "fail: invalid object type",
			args: args{
				r: strings.NewReader("blub 12\x00Hello, World"),
			},
			wantType: UndefinedObject,
			wantSize: 0,
			wantErr:  ErrInvalidObject,
		},
		{
			name: "fail: invalid size",
			args: args{
				r: strings.NewReader("blob xx\x00Hello, World"),
			},
			wantType: UndefinedObject,
			wantSize: 0,
			wantErr:  ErrInvalidObject,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getType, n, err := readHeader(tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if getType != tt.wantType {
				t.Errorf("got = %v, want = %v", getType, tt.wantType)
			}
			if n != tt.wantSize {
				t.Errorf("got = %d, want = %d", n, tt.wantSize)
			}
		})
	}
}

func TestHeader(t *testing.T) {
	tests := []struct {
		name    string
		objType Type
		data    []byte
		want    []byte
	}{
		{
			name:    "success",
			objType: BlobObject,
			data:    []byte("Hello, World"),
			want:    []byte("blob 12\x00"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, _ := NewObject(tt.objType, tt.data)
			b := obj.Header()
			if !reflect.DeepEqual(b, tt.want) {
				t.Errorf("got = %v, want = %v", b, tt.want)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name    string
		objType Type
		data    []byte
		wantErr error
	}{
		{
			name:    "success",
			objType: BlobObject,
			data:    []byte("Hello, World"),
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
			obj, _ := NewObject(tt.objType, tt.data)
			err = obj.Write(goitDir)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			dirPath := filepath.Join(goitDir, "objects", obj.Hash.String()[:2])
			filepath := filepath.Join(dirPath, obj.Hash.String()[2:])
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				t.Error("fail to write goit objects")
			}
			if _, err := os.Stat(filepath); os.IsNotExist(err) {
				t.Error("fail to write goit objects")
			}
		})
	}
}
