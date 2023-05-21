package sha

import (
	"testing"
)

func TestReadHash(t *testing.T) {
	type args struct {
		hashString string
	}
	tests := []struct {
		name string
		args args
		want SHA1
	}{
		{
			name: "success",
			args: args{
				hashString: "1856e9be02756984c385482a07e42f42efd5d2f3",
			},
			want: SHA1([]byte{24, 86, 233, 190, 2, 117, 105, 132, 195, 133, 72, 42, 7, 228, 47, 66, 239, 213, 210, 243}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := ReadHash(tt.args.hashString); string(got) != string(tt.want) {
				t.Errorf("got = %s, want = %s", got, tt.want)
			}
		})
	}
}
