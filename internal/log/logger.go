package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/JunNishimura/Goit/internal/sha"
)

type recordType int

const (
	CommitRecord recordType = iota
	CheckoutRecord
	BranchRecord
)

func (t recordType) String() string {
	switch t {
	case CommitRecord:
		return "commit"
	case CheckoutRecord:
		return "checkout"
	case BranchRecord:
		return "branch"
	default:
		return "undefined"
	}
}

type record struct {
	recType  recordType
	from     sha.SHA1
	to       sha.SHA1
	name     string
	email    string
	unixtime string
	timeDiff string
	message  string
}

func NewRecord(recType recordType, from, to sha.SHA1, name, email string, t time.Time, message string) *record {
	unixtime := fmt.Sprint(t.Unix())
	_, offset := t.Zone()
	offsetMinutes := offset / 60
	timeDiff := fmt.Sprintf("%+03d%02d", offsetMinutes/60, offsetMinutes%60)

	return &record{
		recType:  recType,
		from:     from,
		to:       to,
		name:     name,
		email:    email,
		unixtime: unixtime,
		timeDiff: timeDiff,
		message:  message,
	}
}

func (r *record) String() string {
	var fromStr string
	if r.from == nil {
		fromStr = strings.Repeat("0", 40)
	} else {
		fromStr = r.from.String()
	}
	return fmt.Sprintf("%s %s %s <%s> %s %s\t%s: %s\n", fromStr, r.to, r.name, r.email, r.unixtime, r.timeDiff, r.recType, r.message)
}

type GoitLogger struct {
	rootGoitPath string
}

func NewGoitLogger(rootGoitPath string) *GoitLogger {
	return &GoitLogger{
		rootGoitPath: rootGoitPath,
	}
}

func (l *GoitLogger) WriteHEAD(r *record) error {
	// init logs
	logsPath := filepath.Join(l.rootGoitPath, "logs")
	if _, err := os.Stat(logsPath); os.IsNotExist(err) {
		if err := os.Mkdir(logsPath, os.ModePerm); err != nil {
			return fmt.Errorf("fail to make dir %s: %w", logsPath, err)
		}
	}

	// write HEAD log
	headPath := filepath.Join(logsPath, "HEAD")
	f, err := os.OpenFile(headPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("fail to create %s: %w", headPath, err)
	}
	defer f.Close()

	if _, err := f.WriteString(r.String()); err != nil {
		return fmt.Errorf("fail to write %s: %w", headPath, err)
	}

	return nil
}

func (l *GoitLogger) WriteBranch(r *record, branchName string) error {
	// init logs
	logsPath := filepath.Join(l.rootGoitPath, "logs")
	if _, err := os.Stat(logsPath); os.IsNotExist(err) {
		if err := os.Mkdir(logsPath, os.ModePerm); err != nil {
			return fmt.Errorf("fail to make dir %s: %w", logsPath, err)
		}
	}
	logsRefsPath := filepath.Join(logsPath, "refs")
	if _, err := os.Stat(logsRefsPath); os.IsNotExist(err) {
		if err := os.Mkdir(logsRefsPath, os.ModePerm); err != nil {
			return fmt.Errorf("fail to make dir %s: %w", logsRefsPath, err)
		}
		logsHeadsPath := filepath.Join(logsRefsPath, "heads")
		if err := os.Mkdir(logsHeadsPath, os.ModePerm); err != nil {
			return fmt.Errorf("fail to make dir %s: %w", logsHeadsPath, err)
		}
	}

	// write HEAD log
	branchPath := filepath.Join(logsRefsPath, "heads", branchName)
	f, err := os.OpenFile(branchPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("fail to create %s: %w", branchPath, err)
	}
	defer f.Close()

	if _, err := f.WriteString(r.String()); err != nil {
		return fmt.Errorf("fail to write %s: %w", branchPath, err)
	}

	return nil
}
