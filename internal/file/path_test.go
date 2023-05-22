package file

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFindGoitRoot(t *testing.T) {
	type args struct {
		path string
	}
	type test struct {
		name    string
		args    args
		want    string
		wantErr error
	}
	tests := []*test{
		func() *test {
			rootPath, _ := filepath.Abs("../../")
			testdataPath := filepath.Join(rootPath, "testdata")
			goitPath := filepath.Join(testdataPath, ".goit")

			return &test{
				name: "success",
				args: args{
					path: testdataPath,
				},
				want:    goitPath,
				wantErr: nil,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goitPath := filepath.Join(tt.args.path, ".goit")
			os.Mkdir(goitPath, os.ModePerm)
			defer func() {
				os.RemoveAll(goitPath)
			}()

			goitRootPath, err := FindGoitRoot(tt.args.path)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if goitRootPath != tt.want {
				t.Errorf("got = %v, want = %v", goitRootPath, tt.want)
			}
		})
	}
}
