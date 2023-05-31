package object

import (
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %s, want = %s", got, tt.want)
			}
		})
	}
}
