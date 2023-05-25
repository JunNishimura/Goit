package store

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	type test struct {
		name    string
		isLoad  bool
		content map[string]KV
		want    *Config
		wantErr error
	}
	tests := []*test{
		func() *test {
			config := newConfig()

			return &test{
				name:    "not load",
				isLoad:  false,
				content: nil,
				want:    config,
				wantErr: nil,
			}
		}(),
		func() *test {
			config := newConfig()
			config.Add("user", "name", "test taro")
			config.Add("user", "email", "test@example.com")

			m := make(map[string]KV)
			m["user"] = make(KV)
			m["user"]["name"] = "test taro"
			m["user"]["email"] = "test@example.com"

			return &test{
				name:    "load",
				isLoad:  true,
				content: m,
				want:    config,
				wantErr: nil,
			}
		}(),
	}
	for _, tt := range tests {
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

		cfg, err := NewConfig(goitDir)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("got = %v, want = %v", err, tt.wantErr)
		}

		if tt.isLoad {
			for ident, kv := range tt.content {
				for k, v := range kv {
					cfg.Add(ident, k, v)
				}
			}
		}

		if !reflect.DeepEqual(cfg, tt.want) {
			t.Errorf("got = %v, want = %v", cfg, tt.want)
		}
	}
}

func TestLoad(t *testing.T) {
	type test struct {
		name    string
		content string
		want    *Config
		wantErr error
	}
	tests := []*test{
		func() *test {
			config := newConfig()
			config.Add("user", "name", "Test Taro")
			config.Add("user", "email", "test@example.com")

			return &test{
				name:    "success",
				content: "[user]\n\tname = Test Taro\n\temail = test@example.com\n",
				want:    config,
				wantErr: nil,
			}
		}(),
		{
			name:    "fail",
			content: "[]\n\tname = Test Taro\n\temail = test@example.com\n",
			want:    nil,
			wantErr: ErrInvalidIdentifier,
		},
	}
	for _, tt := range tests {
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
		if _, err := f.WriteString(tt.content); err != nil {
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

		cfg, err := load(configFile)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("got = %v, want = %v", err, tt.wantErr)
		}
		if !reflect.DeepEqual(cfg, tt.want) {
			t.Errorf("got = %v, want = %v", cfg, tt.want)
		}
	}
}

func TestIsUserSet(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "success",
			content: "[user]\n\tname = Test Taro\n\temail = test@example.com\n",
			want:    true,
		},
		{
			name:    "fail",
			content: "",
			want:    false,
		},
	}
	for _, tt := range tests {
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
		if _, err := f.WriteString(tt.content); err != nil {
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

		cfg, err := NewConfig(goitDir)
		if err != nil {
			t.Log(err)
		}
		isUserSet := cfg.IsUserSet()
		if isUserSet != tt.want {
			t.Errorf("got = %v, want = %v", isUserSet, tt.want)
		}
	}
}

func TestAdd(t *testing.T) {
	type test struct {
		name    string
		content map[string]KV
		want    *Config
	}
	tests := []*test{
		func() *test {
			config := newConfig()
			config.Map = make(map[string]KV)
			config.Map["user"] = make(KV)
			config.Map["user"]["name"] = "Test Taro"
			config.Map["user"]["email"] = "test@example.com"

			m := make(map[string]KV)
			m["user"] = make(KV)
			m["user"]["name"] = "Test Taro"
			m["user"]["email"] = "test@example.com"

			return &test{
				name:    "success",
				content: m,
				want:    config,
			}
		}(),
	}
	for _, tt := range tests {
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

		cfg, err := NewConfig(goitDir)
		if err != nil {
			t.Log(err)
		}
		for ident, kv := range tt.content {
			for k, v := range kv {
				cfg.Add(ident, k, v)
			}
		}
		if !reflect.DeepEqual(cfg, tt.want) {
			t.Errorf("got = %v, want = %v", cfg, tt.want)
		}
	}
}

func TestWrite(t *testing.T) {
	type test struct {
		name    string
		content map[string]KV
		want    string
		wantErr error
	}
	tests := []*test{
		func() *test {
			m := make(map[string]KV)
			m["user"] = make(KV)
			m["user"]["name"] = "Test Taro"
			m["user"]["email"] = "test@example.com"

			return &test{
				name:    "success",
				content: m,
				want:    "[user]\n\tname = Test Taro\n\temail = test@example.com\n",
				wantErr: nil,
			}
		}(),
	}
	for _, tt := range tests {
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

		cfg, err := NewConfig(goitDir)
		if err != nil {
			t.Log(err)
		}
		for ident, kv := range tt.content {
			for k, v := range kv {
				cfg.Add(ident, k, v)
			}
		}
		err = cfg.Write(goitDir)
		if !errors.Is(err, tt.wantErr) {
			t.Log(err)
		}
		contentBytes, err := os.ReadFile(configFile)
		if err != nil {
			t.Log(err)
		}
		content := string(contentBytes)
		if content != tt.want {
			t.Errorf("got = %s, want = %s", content, tt.want)
		}
	}
}
