package log

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
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

func TestWriteHEAD(t *testing.T) {
	type args struct {
		rec *record
	}
	type test struct {
		name    string
		args    args
		want    string
		wantErr bool
	}
	tests := []*test{
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			now := time.Now()
			rec := NewRecord(CommitRecord, hash, hash, "Test Taro", "test@example.com", now, "test")

			return &test{
				name: "success: commit record",
				args: args{
					rec: rec,
				},
				want:    fmt.Sprintf("%s %s %s <%s> %s %s\t%s: %s\n", rec.from, rec.to, rec.name, rec.email, rec.unixtime, rec.timeDiff, rec.recType, rec.message),
				wantErr: false,
			}
		}(),
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			now := time.Now()
			rec := NewRecord(BranchRecord, hash, hash, "Test Taro", "test@example.com", now, "test")

			return &test{
				name: "success: branch record",
				args: args{
					rec: rec,
				},
				want:    fmt.Sprintf("%s %s %s <%s> %s %s\t%s: %s\n", rec.from, rec.to, rec.name, rec.email, rec.unixtime, rec.timeDiff, rec.recType, rec.message),
				wantErr: false,
			}
		}(),
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			now := time.Now()
			rec := NewRecord(CheckoutRecord, hash, hash, "Test Taro", "test@example.com", now, "test")

			return &test{
				name: "success: checkout record",
				args: args{
					rec: rec,
				},
				want:    fmt.Sprintf("%s %s %s <%s> %s %s\t%s: %s\n", rec.from, rec.to, rec.name, rec.email, rec.unixtime, rec.timeDiff, rec.recType, rec.message),
				wantErr: false,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			gLogger := NewGoitLogger(tmpDir)
			if err := gLogger.WriteHEAD(tt.args.rec); (err != nil) != tt.wantErr {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}

			headPath := filepath.Join(tmpDir, "logs", "HEAD")
			b, err := os.ReadFile(headPath)
			if err != nil {
				t.Log(err)
			}
			got := string(b)
			if got != tt.want {
				t.Errorf("got = %s, want = %s", got, tt.want)
			}
		})
	}
}

func TestWriteBranch(t *testing.T) {
	type args struct {
		rec        *record
		branchName string
	}
	type test struct {
		name    string
		args    args
		want    string
		wantErr bool
	}
	tests := []*test{
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			now := time.Now()
			rec := NewRecord(CommitRecord, hash, hash, "Test Taro", "test@example.com", now, "test")

			return &test{
				name: "success: commit record",
				args: args{
					rec:        rec,
					branchName: "test",
				},
				want:    fmt.Sprintf("%s %s %s <%s> %s %s\t%s: %s\n", rec.from, rec.to, rec.name, rec.email, rec.unixtime, rec.timeDiff, rec.recType, rec.message),
				wantErr: false,
			}
		}(),
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			now := time.Now()
			rec := NewRecord(BranchRecord, hash, hash, "Test Taro", "test@example.com", now, "test")

			return &test{
				name: "success: branch record",
				args: args{
					rec:        rec,
					branchName: "test",
				},
				want:    fmt.Sprintf("%s %s %s <%s> %s %s\t%s: %s\n", rec.from, rec.to, rec.name, rec.email, rec.unixtime, rec.timeDiff, rec.recType, rec.message),
				wantErr: false,
			}
		}(),
		func() *test {
			hash, _ := hex.DecodeString("87f3c49bccf2597484ece08746d3ee5defaba335")
			now := time.Now()
			rec := NewRecord(CheckoutRecord, hash, hash, "Test Taro", "test@example.com", now, "test")

			return &test{
				name: "success: checkout record",
				args: args{
					rec:        rec,
					branchName: "test",
				},
				want:    fmt.Sprintf("%s %s %s <%s> %s %s\t%s: %s\n", rec.from, rec.to, rec.name, rec.email, rec.unixtime, rec.timeDiff, rec.recType, rec.message),
				wantErr: false,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			gLogger := NewGoitLogger(tmpDir)
			if err := gLogger.WriteBranch(tt.args.rec, tt.args.branchName); (err != nil) != tt.wantErr {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}

			branchPath := filepath.Join(tmpDir, "logs", "refs", "heads", tt.args.branchName)
			b, err := os.ReadFile(branchPath)
			if err != nil {
				t.Log(err)
			}
			got := string(b)
			if got != tt.want {
				t.Errorf("got = %s, want = %s", got, tt.want)
			}
		})
	}
}
