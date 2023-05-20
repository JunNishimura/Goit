package object

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/JunNishimura/Goit/sha"
)

var (
	emailRegexpString     = "([a-zA-Z0-9_.+-]+@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\\.)+[a-zA-Z]{2,})"
	timestampRegexpString = "([1-9][0-9]* [+-][0-9]{4})"
	signRegexp            = regexp.MustCompile("^[^<]* <" + emailRegexpString + "> " + timestampRegexpString + "$")
)

type Sign struct {
	Name      string
	Email     string
	Timestamp time.Time
}

func (s Sign) String() string {
	unixTime := s.Timestamp.Unix()
	_, offsetSec := s.Timestamp.Zone()
	offsetHour := offsetSec / 3600
	offsetMinute := (offsetSec / 60) % 60
	var posNegSign string
	if offsetSec >= 0 {
		posNegSign = "+"
	} else {
		posNegSign = "-"
	}
	offset := fmt.Sprintf("%s%02d%02d", posNegSign, offsetHour, offsetMinute)
	return fmt.Sprintf("%s <%s> %s %s", s.Name, s.Email, fmt.Sprint(unixTime), offset)
}

func NewSign(name, email string) *Sign {
	return &Sign{
		Name:      name,
		Email:     email,
		Timestamp: time.Now(),
	}
}

type Commit struct {
	*Object
	Tree      sha.SHA1
	Parents   []sha.SHA1
	Author    Sign
	Committer Sign
	Message   string
}

func NewCommit(o *Object) (*Commit, error) {
	if o.Type != CommitObject {
		return nil, ErrNotCommitObject
	}

	commit := &Commit{
		Object: o,
	}

	buf := bytes.NewReader(o.Data)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		text := scanner.Text()
		splitText := strings.SplitN(text, " ", 2)
		if len(splitText) != 2 {
			break
		}

		lineType := splitText[0]
		body := splitText[1]

		switch lineType {
		case "tree":
			hash, err := sha.ReadHash(body)
			if err != nil {
				return nil, err
			}
			commit.Tree = hash
		case "parent":
			hash, err := sha.ReadHash(body)
			if err != nil {
				return nil, err
			}
			commit.Parents = append(commit.Parents, hash)
		case "author":
			sign, err := readSign(body)
			if err != nil {
				return nil, err
			}
			commit.Author = sign
		case "committer":
			sign, err := readSign(body)
			if err != nil {
				return nil, err
			}
			commit.Committer = sign
		}
	}

	message := make([]string, 0)
	for scanner.Scan() {
		message = append(message, scanner.Text())
	}
	commit.Message = strings.Join(message, "\n")

	return commit, nil
}

func (c *Commit) String() string {
	var commitString string

	authorSignString := c.Author.String()
	sign1 := strings.SplitN(authorSignString, " <", 2)
	author := sign1[0]
	sign2 := strings.SplitN(sign1[1], "> ", 2)
	email := sign2[0]

	commitString += fmt.Sprintf("commit %s\n", c.Hash)
	commitString += fmt.Sprintf("Author: %s <%s>\n", author, email)
	commitString += fmt.Sprintf("Date: %s\n", c.Author.Timestamp)
	commitString += fmt.Sprintf("\n\t%s\n", c.Message)

	return commitString
}

func (c *Commit) UpdateBranch(branchPath string) error {
	f, err := os.Create(branchPath)
	if err != nil {
		return fmt.Errorf("fail to make %s: %v", branchPath, err)
	}
	defer f.Close()

	if _, err := f.WriteString(c.Hash.String()); err != nil {
		return fmt.Errorf("fail to write hash to %s: %v", branchPath, err)
	}

	return nil
}

func readSign(signString string) (Sign, error) {
	if ok := signRegexp.MatchString(signString); !ok {
		return Sign{}, ErrInvalidCommitObject
	}
	sign1 := strings.SplitN(signString, " <", 2)
	name := sign1[0]
	sign2 := strings.SplitN(sign1[1], "> ", 2)
	email := sign2[0]
	sign3 := strings.SplitN(sign2[1], " ", 2)
	unixTimeString := sign3[0]
	offsetString := sign3[1]

	unixTime, err := strconv.ParseInt(unixTimeString, 10, 64)
	if err != nil {
		return Sign{}, fmt.Errorf("%w : %s", ErrInvalidCommitObject, err)
	}
	var offsetHour, offsetMinute int
	switch offsetString[:1] {
	case "+":
		if _, err := fmt.Sscanf(offsetString, "+%02d%02d", &offsetHour, &offsetMinute); err != nil {
			return Sign{}, fmt.Errorf("%w : %s", ErrInvalidCommitObject, err)
		}
	case "-":
		if _, err := fmt.Sscanf(offsetString, "-%02d%02d", &offsetHour, &offsetMinute); err != nil {
			return Sign{}, fmt.Errorf("%w : %s", ErrInvalidCommitObject, err)
		}
	}
	location := time.FixedZone(" ", 3600*offsetHour+60*offsetMinute)
	timestamp := time.Unix(unixTime, 0).In(location)
	sign := Sign{
		Name:      name,
		Email:     email,
		Timestamp: timestamp,
	}

	return sign, nil
}
