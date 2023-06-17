package store

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewIgnore(t *testing.T) {
	type fields struct {
		content string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Ignore
		wantErr bool
	}{
		{
			name: "success: empty",
			fields: fields{
				content: "",
			},
			want: &Ignore{
				paths: []string{`\.goit/.*`},
			},
			wantErr: false,
		},
		{
			name: "success: some ignore list",
			fields: fields{
				content: "*.exe\ndir/",
			},
			want: &Ignore{
				paths: []string{`\.goit/.*`, `.*\.exe`, `dir/.*`},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			goitignorePath := filepath.Join(tmpDir, ".goitignore")
			f, err := os.Create(goitignorePath)
			if err != nil {
				t.Log(err)
			}
			_, err = f.WriteString(tt.fields.content)
			if err != nil {
				t.Log(err)
			}
			f.Close()

			goitPath := filepath.Join(tmpDir, ".goit")
			i, err := NewIgnore(goitPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(i, tt.want) {
				t.Errorf("got = %v, want = %v", i, tt.want)
			}
		})
	}
}
