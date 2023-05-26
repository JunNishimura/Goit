package file

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestFindGoitRoot(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "success",
			args: args{
				path: ".",
			},
			want:    ".goit",
			wantErr: nil,
		},
		{
			name: "fail",
			args: args{
				path: "../",
			},
			want:    "",
			wantErr: ErrGoitRootNotFound,
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

			rootPath := filepath.Join(tmpDir, tt.args.path)
			goitRootPath, err := FindGoitRoot(rootPath)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}

			var wantPath string
			if tt.want != "" {
				wantPath = filepath.Join(tmpDir, tt.want)
			}
			if goitRootPath != wantPath {
				t.Errorf("got = %v, want = %v", goitRootPath, wantPath)
			}
		})
	}
}

func TestGetFilePathsUnderDirectory(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr error
	}{
		{
			name:    "success",
			want:    []string{"test.txt", "test2.txt", "dir/test1.txt", "dir/test2.txt"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			f, err := os.Create(filepath.Join(tmpDir, "test.txt"))
			if err != nil {
				t.Log(err)
			}
			f.Close()
			f, err = os.Create(filepath.Join(tmpDir, "test2.txt"))
			if err != nil {
				t.Log(err)
			}
			f.Close()
			if err := os.Mkdir(filepath.Join(tmpDir, "dir"), os.ModePerm); err != nil {
				t.Log(err)
			}
			f, err = os.Create(filepath.Join(tmpDir, "dir/test1.txt"))
			if err != nil {
				t.Log(err)
			}
			f.Close()
			f, err = os.Create(filepath.Join(tmpDir, "dir/test2.txt"))
			if err != nil {
				t.Log(err)
			}
			f.Close()

			got, err := GetFilePathsUnderDirectory(tmpDir)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			var filePaths []string
			for _, path := range tt.want {
				filePaths = append(filePaths, filepath.Join(tmpDir, path))
			}
			sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
			sort.Slice(filePaths, func(i, j int) bool { return filePaths[i] < filePaths[j] })
			if !reflect.DeepEqual(got, filePaths) {
				t.Errorf("got = %v, want = %v", got, filePaths)
			}
		})
	}
}
