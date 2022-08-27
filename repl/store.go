package repl

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

type Store interface {
	Write(string)
	Read() string
}

type SSHStore struct {
	Cid      int
	Sid      int
	Instance map[int]any
}

func (s *SSHStore) Write(data string) {
	if s.Cid > -1 && s.Sid > -1 {
		tmpl := s.Instance[s.Cid].(map[int][]string)
		tmpl[s.Sid] = append(tmpl[s.Sid], data)
		s.Instance[s.Cid] = tmpl
	}
}

func (s *SSHStore) Read() string {
	res := ""
	tmpl := s.Instance[s.Cid].(map[int][]string)
	for _, v := range tmpl[s.Sid] {
		res += v
	}
	return res
}

func (s *SSHStore) Save() {
	fp, err := OpenFile("./record.txt")
	if err != nil {
		logrus.Fatal(err.Error())
	}
	defer fp.Close()
	_, err = io.WriteString(fp, s.Read())
	if err != nil {
		logrus.Fatal(err.Error())
	}
	fmt.Print("Save successfully\n")
}

func (s *SSHStore) Set(cid, sid int) {
	s.Cid = cid
	s.Sid = sid
}

func GetStore() *SSHStore {
	return store
}

func NewSSHStore() *SSHStore {
	return &SSHStore{
		Cid:      -1,
		Sid:      -1,
		Instance: make(map[int]any),
	}
}
