package object

import "testing"

func TestNewType(t *testing.T) {
	type args struct {
		typeString string
	}
	tests := []struct {
		name string
		args args
		want Type
	}{
		{
			name: "blob",
			args: args{
				typeString: "blob",
			},
			want: BlobObject,
		},
		{
			name: "tree",
			args: args{
				typeString: "tree",
			},
			want: TreeObject,
		},
		{
			name: "commit",
			args: args{
				typeString: "commit",
			},
			want: CommitObject,
		},
		{
			name: "tag",
			args: args{
				typeString: "tag",
			},
			want: TagObject,
		},
		{
			name: "undefined",
			args: args{
				typeString: "xxx",
			},
			want: UndefinedObject,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := NewType(tt.args.typeString); got != tt.want {
				t.Errorf("got = %s, want = %s", got, tt.want)
			}
		})
	}
}
