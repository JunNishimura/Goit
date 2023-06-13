package log

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/JunNishimura/Goit/internal/sha"
)

func TestNewRecord(t *testing.T) {
	type args struct {
		recType recordType
		from    sha.SHA1
		to      sha.SHA1
		name    string
		email   string
		t       time.Time
		message string
	}
	type test struct {
		name string
		args args
		want *record
	}
	tests := []*test{
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			now := time.Now()
			unixtime := fmt.Sprint(now.Unix())
			_, offset := now.Zone()
			offsetMinutes := offset / 60
			timeDiff := fmt.Sprintf("%+03d%02d", offsetMinutes/60, offsetMinutes%60)

			return &test{
				name: "success: commit record",
				args: args{
					recType: CommitRecord,
					from:    sha.SHA1(hash),
					to:      sha.SHA1(hash),
					name:    "Test Taro",
					email:   "test@example.com",
					t:       now,
					message: "test",
				},
				want: &record{
					recType:  CommitRecord,
					from:     sha.SHA1(hash),
					to:       sha.SHA1(hash),
					name:     "Test Taro",
					email:    "test@example.com",
					unixtime: unixtime,
					timeDiff: timeDiff,
					message:  "test",
				},
			}
		}(),
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			now := time.Now()
			unixtime := fmt.Sprint(now.Unix())
			_, offset := now.Zone()
			offsetMinutes := offset / 60
			timeDiff := fmt.Sprintf("%+03d%02d", offsetMinutes/60, offsetMinutes%60)

			return &test{
				name: "success: branch record",
				args: args{
					recType: BranchRecord,
					from:    sha.SHA1(hash),
					to:      sha.SHA1(hash),
					name:    "Test Taro",
					email:   "test@example.com",
					t:       now,
					message: "test",
				},
				want: &record{
					recType:  BranchRecord,
					from:     sha.SHA1(hash),
					to:       sha.SHA1(hash),
					name:     "Test Taro",
					email:    "test@example.com",
					unixtime: unixtime,
					timeDiff: timeDiff,
					message:  "test",
				},
			}
		}(),
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			now := time.Now()
			unixtime := fmt.Sprint(now.Unix())
			_, offset := now.Zone()
			offsetMinutes := offset / 60
			timeDiff := fmt.Sprintf("%+03d%02d", offsetMinutes/60, offsetMinutes%60)

			return &test{
				name: "success: checkout record",
				args: args{
					recType: CheckoutRecord,
					from:    sha.SHA1(hash),
					to:      sha.SHA1(hash),
					name:    "Test Taro",
					email:   "test@example.com",
					t:       now,
					message: "test",
				},
				want: &record{
					recType:  CheckoutRecord,
					from:     sha.SHA1(hash),
					to:       sha.SHA1(hash),
					name:     "Test Taro",
					email:    "test@example.com",
					unixtime: unixtime,
					timeDiff: timeDiff,
					message:  "test",
				},
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRecord(tt.args.recType, tt.args.from, tt.args.to, tt.args.name, tt.args.email, tt.args.t, tt.args.message)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}
