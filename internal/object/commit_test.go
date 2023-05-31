package object

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestNewSign(t *testing.T) {
	type args struct {
		name  string
		email string
	}
	tests := []struct {
		name string
		args args
		want *Sign
	}{
		{
			name: "success",
			args: args{
				name:  "test taro",
				email: "test@example.com",
			},
			want: &Sign{
				Name:      "test taro",
				Email:     "test@example.com",
				Timestamp: time.Now(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSign(tt.args.name, tt.args.email)
			if got.Name != tt.want.Name {
				t.Errorf("got = %s, want = %s", got.Name, tt.want.Name)
			}
			if got.Email != tt.want.Email {
				t.Errorf("got = %s, want = %s", got.Email, tt.want.Email)
			}
		})
	}
}

func TestNewCommit(t *testing.T) {
	type args struct {
		obj *Object
	}
	type test struct {
		name    string
		args    args
		want    *Commit
		wantErr error
	}
	tests := []*test{
		func() *test {
			obj, _ := NewObject(BlobObject, []byte("blob 12\x00Hello, World"))

			return &test{
				name: "fail: blob object",
				args: args{
					obj: obj,
				},
				want:    nil,
				wantErr: ErrNotCommitObject,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCommit(tt.args.obj)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}
