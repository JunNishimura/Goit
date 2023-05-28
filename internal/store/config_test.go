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
		name          string
		localContent  string
		globalContent string
		want          *Config
		wantErr       error
	}
	tests := []*test{
		func() *test {
			return &test{
				name: "success: no load",
				want: &Config{
					local:  make(map[string]kv),
					global: make(map[string]kv),
				},
				wantErr: nil,
			}
		}(),
		func() *test {
			config := newConfig()
			config.Add("user", "name", "test taro", false)
			config.Add("user", "email", "test@example.com", false)

			return &test{
				name:         "success: local load",
				localContent: "[user]\n\tname = test taro\n\temail = test@example.com\n",
				want:         config,
				wantErr:      nil,
			}
		}(),
		func() *test {
			config := newConfig()
			config.Add("user", "name", "test taro", true)
			config.Add("user", "email", "test@example.com", true)

			return &test{
				name:          "success: global load",
				globalContent: "[user]\n\tname = test taro\n\temail = test@example.com\n",
				want:          config,
				wantErr:       nil,
			}
		}(),
		func() *test {
			config := newConfig()
			config.Add("hoge", "piyo", "poyo", false)
			config.Add("hoge", "foo", "bar", false)
			config.Add("user", "name", "test taro", true)
			config.Add("user", "email", "test@example.com", true)

			return &test{
				name:          "success: global load",
				localContent:  "[hoge]\n\tpiyo = poyo\n\tfoo = bar\n",
				globalContent: "[user]\n\tname = test taro\n\temail = test@example.com\n",
				want:          config,
				wantErr:       nil,
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
		localConfigFile := filepath.Join(goitDir, "config")
		f, err := os.Create(localConfigFile)
		if err != nil {
			t.Logf("%v: %s", err, localConfigFile)
		}
		if tt.localContent != "" {
			_, err := f.WriteString(tt.localContent)
			if err != nil {
				t.Logf("%v: %s", err, localConfigFile)
			}
		}
		f.Close()
		// make global/config file
		userHomePath, err := os.UserHomeDir()
		if err != nil {
			t.Logf("%v: %s", err, userHomePath)
		}
		globalConfigPath := filepath.Join(userHomePath, ".goitconfig")
		f, err = os.Create(globalConfigPath)
		if err != nil {
			t.Logf("%v: %s", err, globalConfigPath)
		}
		if tt.globalContent != "" {
			_, err := f.WriteString(tt.globalContent)
			if err != nil {
				t.Logf("%v: %s", err, globalConfigPath)
			}
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

		if !reflect.DeepEqual(cfg, tt.want) {
			t.Errorf("got = %v, want = %v", cfg, tt.want)
		}
	}
}

func TestIsUserSet(t *testing.T) {
	tests := []struct {
		name          string
		localContent  string
		globalContent string
		want          bool
	}{
		{
			name:         "success: local set",
			localContent: "[user]\n\tname = Test Taro\n\temail = test@example.com\n",
			want:         true,
		},
		{
			name:          "success: global set",
			globalContent: "[user]\n\tname = Test Taro\n\temail = test@example.com\n",
			want:          true,
		},
		{
			name: "fail",
			want: false,
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
		// make .goit/config file
		localConfigFile := filepath.Join(goitDir, "config")
		f, err := os.Create(localConfigFile)
		if err != nil {
			t.Logf("%v: %s", err, localConfigFile)
		}
		if tt.localContent != "" {
			_, err := f.WriteString(tt.localContent)
			if err != nil {
				t.Logf("%v: %s", err, localConfigFile)
			}
		}
		f.Close()
		// make global/config file
		userHomePath, err := os.UserHomeDir()
		if err != nil {
			t.Logf("%v: %s", err, userHomePath)
		}
		globalConfigPath := filepath.Join(userHomePath, ".goitconfig")
		f, err = os.Create(globalConfigPath)
		if err != nil {
			t.Logf("%v: %s", err, globalConfigPath)
		}
		if tt.globalContent != "" {
			_, err := f.WriteString(tt.globalContent)
			if err != nil {
				t.Logf("%v: %s", err, globalConfigPath)
			}
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

func TestGetUserName(t *testing.T) {
	tests := []struct {
		name          string
		localContent  string
		globalContent string
		want          string
	}{
		{
			name:         "success: get from local",
			localContent: "[user]\n\tname = Test Taro\n\temail = test@example.com\n",
			want:         "Test Taro",
		},
		{
			name:          "success: get from global",
			globalContent: "[user]\n\tname = Test Taro\n\temail = test@example.com\n",
			want:          "Test Taro",
		},
		{
			name:          "success: get from local when both local and global set",
			localContent:  "[user]\n\tname = Test Hanako\n\temail = test@example.com\n",
			globalContent: "[user]\n\tname = Test Taro\n\temail = test@example.com\n",
			want:          "Test Hanako",
		},
		{
			name: "fail",
			want: "",
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
		// make .goit/config file
		localConfigFile := filepath.Join(goitDir, "config")
		f, err := os.Create(localConfigFile)
		if err != nil {
			t.Logf("%v: %s", err, localConfigFile)
		}
		if tt.localContent != "" {
			_, err := f.WriteString(tt.localContent)
			if err != nil {
				t.Logf("%v: %s", err, localConfigFile)
			}
		}
		f.Close()
		// make global/config file
		userHomePath, err := os.UserHomeDir()
		if err != nil {
			t.Logf("%v: %s", err, userHomePath)
		}
		globalConfigPath := filepath.Join(userHomePath, ".goitconfig")
		f, err = os.Create(globalConfigPath)
		if err != nil {
			t.Logf("%v: %s", err, globalConfigPath)
		}
		if tt.globalContent != "" {
			_, err := f.WriteString(tt.globalContent)
			if err != nil {
				t.Logf("%v: %s", err, globalConfigPath)
			}
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
		userName := cfg.GetUserName()
		if userName != tt.want {
			t.Errorf("got = %s, want = %s", userName, tt.want)
		}
	}
}

// func TestAdd(t *testing.T) {
// 	type test struct {
// 		name    string
// 		content map[string]KV
// 		want    *Config
// 	}
// 	tests := []*test{
// 		func() *test {
// 			config := newConfig()
// 			config.Map = make(map[string]KV)
// 			config.Map["user"] = make(KV)
// 			config.Map["user"]["name"] = "Test Taro"
// 			config.Map["user"]["email"] = "test@example.com"

// 			m := make(map[string]KV)
// 			m["user"] = make(KV)
// 			m["user"]["name"] = "Test Taro"
// 			m["user"]["email"] = "test@example.com"

// 			return &test{
// 				name:    "success",
// 				content: m,
// 				want:    config,
// 			}
// 		}(),
// 	}
// 	for _, tt := range tests {
// 		tmpDir := t.TempDir()
// 		// .goit initialization
// 		goitDir := filepath.Join(tmpDir, ".goit")
// 		if err := os.Mkdir(goitDir, os.ModePerm); err != nil {
// 			t.Logf("%v: %s", err, goitDir)
// 		}
// 		// make .goit/config file
// 		configFile := filepath.Join(goitDir, "config")
// 		f, err := os.Create(configFile)
// 		if err != nil {
// 			t.Logf("%v: %s", err, configFile)
// 		}
// 		f.Close()
// 		// make .goit/HEAD file and write main branch
// 		headFile := filepath.Join(goitDir, "HEAD")
// 		f, err = os.Create(headFile)
// 		if err != nil {
// 			t.Logf("%v: %s", err, headFile)
// 		}
// 		// set 'main' as default branch
// 		if _, err := f.WriteString("ref: refs/heads/main"); err != nil {
// 			t.Logf("%v: %s", err, headFile)
// 		}
// 		f.Close()
// 		// make .goit/objects directory
// 		objectsDir := filepath.Join(goitDir, "objects")
// 		if err := os.Mkdir(objectsDir, os.ModePerm); err != nil {
// 			t.Logf("%v: %s", err, objectsDir)
// 		}
// 		// make .goit/refs directory
// 		refsDir := filepath.Join(goitDir, "refs")
// 		if err := os.Mkdir(refsDir, os.ModePerm); err != nil {
// 			t.Logf("%v: %s", err, refsDir)
// 		}
// 		// make .goit/refs/heads directory
// 		headsDir := filepath.Join(refsDir, "heads")
// 		if err := os.Mkdir(headsDir, os.ModePerm); err != nil {
// 			t.Logf("%v: %s", err, headsDir)
// 		}
// 		// make .goit/refs/tags directory
// 		tagsDir := filepath.Join(refsDir, "tags")
// 		if err := os.Mkdir(tagsDir, os.ModePerm); err != nil {
// 			t.Logf("%v: %s", err, tagsDir)
// 		}

// 		cfg, err := NewConfig(goitDir)
// 		if err != nil {
// 			t.Log(err)
// 		}
// 		for ident, kv := range tt.content {
// 			for k, v := range kv {
// 				cfg.Add(ident, k, v)
// 			}
// 		}
// 		if !reflect.DeepEqual(cfg, tt.want) {
// 			t.Errorf("got = %v, want = %v", cfg, tt.want)
// 		}
// 	}
// }

// func TestWrite(t *testing.T) {
// 	type test struct {
// 		name    string
// 		content map[string]KV
// 		want1   string
// 		want2   string
// 		wantErr error
// 	}
// 	tests := []*test{
// 		func() *test {
// 			m := make(map[string]KV)
// 			m["user"] = make(KV)
// 			m["user"]["name"] = "Test Taro"
// 			m["user"]["email"] = "test@example.com"

// 			return &test{
// 				name:    "success",
// 				content: m,
// 				want1:   "[user]\n\tname = Test Taro\n\temail = test@example.com\n",
// 				want2:   "[user]\n\temail = test@example.com\n\tname = Test Taro\n",
// 				wantErr: nil,
// 			}
// 		}(),
// 	}
// 	for _, tt := range tests {
// 		tmpDir := t.TempDir()
// 		// .goit initialization
// 		goitDir := filepath.Join(tmpDir, ".goit")
// 		if err := os.Mkdir(goitDir, os.ModePerm); err != nil {
// 			t.Logf("%v: %s", err, goitDir)
// 		}
// 		// make .goit/config file
// 		configFile := filepath.Join(goitDir, "config")
// 		f, err := os.Create(configFile)
// 		if err != nil {
// 			t.Logf("%v: %s", err, configFile)
// 		}
// 		f.Close()
// 		// make .goit/HEAD file and write main branch
// 		headFile := filepath.Join(goitDir, "HEAD")
// 		f, err = os.Create(headFile)
// 		if err != nil {
// 			t.Logf("%v: %s", err, headFile)
// 		}
// 		// set 'main' as default branch
// 		if _, err := f.WriteString("ref: refs/heads/main"); err != nil {
// 			t.Logf("%v: %s", err, headFile)
// 		}
// 		f.Close()
// 		// make .goit/objects directory
// 		objectsDir := filepath.Join(goitDir, "objects")
// 		if err := os.Mkdir(objectsDir, os.ModePerm); err != nil {
// 			t.Logf("%v: %s", err, objectsDir)
// 		}
// 		// make .goit/refs directory
// 		refsDir := filepath.Join(goitDir, "refs")
// 		if err := os.Mkdir(refsDir, os.ModePerm); err != nil {
// 			t.Logf("%v: %s", err, refsDir)
// 		}
// 		// make .goit/refs/heads directory
// 		headsDir := filepath.Join(refsDir, "heads")
// 		if err := os.Mkdir(headsDir, os.ModePerm); err != nil {
// 			t.Logf("%v: %s", err, headsDir)
// 		}
// 		// make .goit/refs/tags directory
// 		tagsDir := filepath.Join(refsDir, "tags")
// 		if err := os.Mkdir(tagsDir, os.ModePerm); err != nil {
// 			t.Logf("%v: %s", err, tagsDir)
// 		}

// 		cfg, err := NewConfig(goitDir)
// 		if err != nil {
// 			t.Log(err)
// 		}
// 		for ident, kv := range tt.content {
// 			for k, v := range kv {
// 				cfg.Add(ident, k, v)
// 			}
// 		}
// 		err = cfg.Write(goitDir)
// 		if !errors.Is(err, tt.wantErr) {
// 			t.Log(err)
// 		}
// 		contentBytes, err := os.ReadFile(configFile)
// 		if err != nil {
// 			t.Log(err)
// 		}
// 		content := string(contentBytes)
// 		if content != tt.want1 && content != tt.want2 {
// 			t.Errorf("got = %s, want1 = %s, want2 = %s", content, tt.want1, tt.want2)
// 		}
// 	}
// }
