package store

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReflogLoad(t *testing.T) {
	type args struct {
		head *Head
		refs *Refs
	}
	type test struct {
		name    string
		args    args
		fields  string
		want    *Reflog
		wantErr bool
	}
	tests := []*test{
		func() *test {
			return &test{
				name: "success",
				args: args{
					head: &Head{},
					refs: &Refs{},
				},
				fields: "",
				want: &Reflog{
					records: make([]*LogRecord, 0),
				},
				wantErr: false,
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
			// make .goit/logs directory
			logsDir := filepath.Join(goitDir, "logs")
			if err := os.Mkdir(logsDir, os.ModePerm); err != nil {
				t.Logf("%v: %s", err, logsDir)
			}
			// make HEAD file
			headPath := filepath.Join(logsDir, "HEAD")
			f, err := os.Create(headPath)
			if err != nil {
				t.Logf("%v: %s", err, headPath)
			}
			if _, err := f.WriteString(tt.fields); err != nil {
				t.Log(err)
			}
			f.Close()

			reflog := newReflog()
			if err := reflog.load(goitDir, tt.args.head, tt.args.refs); (err != nil) != tt.wantErr {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(reflog, tt.want) {
				t.Errorf("got = %v, want = %v", reflog, tt.want)
			}
		})
	}
}
