package store

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestNewHead(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		isCreated bool
		want      Head
		wantErr   error
	}{
		{
			name:      "success",
			content:   "ref: refs/heads/main",
			isCreated: true,
			want:      "main",
			wantErr:   nil,
		}, {
			name:      "invalid HEAD format",
			content:   "ref: ***",
			isCreated: true,
			want:      "",
			wantErr:   ErrInvalidHead,
		}, {
			name:      "no HEAD file",
			content:   "ref: refs/heads/main",
			isCreated: false,
			want:      "",
			wantErr:   nil,
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
			if tt.isCreated {
				headFile := filepath.Join(goitDir, "HEAD")
				f, err := os.Create(headFile)
				if err != nil {
					t.Logf("%v: %s", err, headFile)
				}
				// set 'main' as default branch
				if _, err := f.WriteString(tt.content); err != nil {
					t.Logf("%v: %s", err, headFile)
				}
				f.Close()
			}
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

			head, err := NewHead(goitDir)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if head != tt.want {
				t.Errorf("got = %s, want = %s", head, tt.want)
			}
		})
	}
}
