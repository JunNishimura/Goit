package sha

import (
	"encoding/hex"
	"errors"
	"strings"
	"testing"
)

func TestCompare(t *testing.T) {
	type args struct {
		sha SHA1
	}
	type fields struct {
		sha SHA1
	}
	type test struct {
		name   string
		args   args
		fields fields
		want   bool
	}
	tests := []*test{
		func() *test {
			s, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")

			return &test{
				name: "success: true",
				args: args{
					sha: s,
				},
				fields: fields{
					sha: s,
				},
				want: true,
			}
		}(),
		func() *test {
			s, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			s2, _ := hex.DecodeString(strings.Repeat("0", 40))

			return &test{
				name: "success: false",
				args: args{
					sha: s,
				},
				fields: fields{
					sha: s2,
				},
				want: false,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.sha.Compare(tt.args.sha)
			if got != tt.want {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestReadHash(t *testing.T) {
	type args struct {
		hashString string
	}
	tests := []struct {
		name    string
		args    args
		want    SHA1
		wantErr error
	}{
		{
			name: "success",
			args: args{
				hashString: "1856e9be02756984c385482a07e42f42efd5d2f3",
			},
			want:    SHA1([]byte{24, 86, 233, 190, 2, 117, 105, 132, 195, 133, 72, 42, 7, 228, 47, 66, 239, 213, 210, 243}),
			wantErr: nil,
		},
		{
			name: "failure",
			args: args{
				hashString: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			},
			want:    SHA1([]byte{}),
			wantErr: ErrInvalidHash,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadHash(tt.args.hashString)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
			if string(got) != string(tt.want) {
				t.Errorf("got = %s, want = %s", got, tt.want)
			}
		})
	}
}
