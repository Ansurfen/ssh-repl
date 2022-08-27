package repl

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

type SSHReader struct {
	channel chan string
}

func NewSSHReader() *SSHReader {
	return &SSHReader{channel: make(chan string, 2)}
}

type SSHWriter struct {
	channel chan string
}

func NewSSHWriter() *SSHWriter {
	return &SSHWriter{channel: make(chan string, 2)}
}

func (r *SSHReader) Read(p []byte) (n int, err error) {
	cmd := <-r.channel
	tmpl := []byte(cmd + "\n")
	copy(p, tmpl)
	return len(tmpl), err
}

func (w *SSHWriter) Write(p []byte) (n int, err error) {
	w.channel <- string(p)
	return len(p), err
}

type SSHCenter struct {
	Cid     int
	Sid     int
	Clients []*SSHClient
	Command string
}

type SSHClient struct {
	Cid      int
	Instance *ssh.Client
	Opts     *SSHOpts
	Sessions []*SSHSession
}

type SSHOpts struct {
	Conf    *viper.Viper
	Network string
	Addr    string
	Port    string
	User    string
	Key     string
	Kpath   string
}

type SSHSession struct {
	Sid      int
	Instance *ssh.Session
	Writer   *SSHWriter
	Reader   *SSHReader
	Lock     sync.Mutex
}

func (opts *SSHOpts) publicKeyAuthFunc() ssh.AuthMethod {
	if len(opts.Key) == 0 {
		key, err := ioutil.ReadFile(opts.Kpath)
		opts.Key = string(key)
		PanicWithLogger(err, "Fail to read the key's file of ssh")
	}
	signer, err := ssh.ParsePrivateKey([]byte(opts.Key))
	PanicWithLogger(err, "Fail to sign the key of ssh")
	return ssh.PublicKeys(signer)
}

func (c *SSHClient) Connect() {
	var err error
	c.Instance, err = ssh.Dial(c.Opts.Network, c.Opts.Addr+":"+c.Opts.Port, &ssh.ClientConfig{
		User:            c.Opts.User,
		Auth:            []ssh.AuthMethod{c.Opts.publicKeyAuthFunc()},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	PanicWithLogger(err, "Fail to connect client")
}

func (c *SSHClient) NewSession() *SSHSession {
	s := &SSHSession{}
	session, err := c.Instance.NewSession()
	PanicWithLogger(err, "Fail to create a new session")
	s.Instance = session
	s.Reader = NewSSHReader()
	s.Writer = NewSSHWriter()
	session.Stdout = s.Writer
	session.Stdin = s.Reader
	session.Stderr = os.Stderr
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("xterm", 25, 80, modes)
	PanicWithLogger(err, "Fail to set `request pty`")
	err = session.Shell()
	PanicWithLogger(err, "Fail to start shell")
	s.Sid = len(c.Sessions)
	s.ReadWithStore()
	c.Sessions = append(c.Sessions, s)
	return s
}

func (s *SSHSession) ReadWithStore() {
	res := s.Read()
	var tmpl map[int][]string
	if store.Instance[store.Cid] == nil {
		store.Instance[store.Cid] = make(map[int][]string)
	}
	tmpl = store.Instance[store.Cid].(map[int][]string)
	tmpl[s.Sid] = append(tmpl[s.Sid], res)
	store.Instance[store.Cid] = tmpl
	store.Sid = s.Sid
}

func (s *SSHSession) Read() string {
	s.Lock.Lock()
	res := <-s.Writer.channel
	s.Lock.Unlock()
	fmt.Print(res)
	return res
}

func (s *SSHSession) Write(cmd string) {
	s.Reader.channel <- Strip(strings.TrimSpace(cmd))
}

func (s *SSHSession) ExecCommand(cmd string) {
	s.Write(cmd)
	s.ReadWithStore()
}

func (center *SSHCenter) NewClient(conf string) {
	client := &SSHClient{
		Opts: center.NewConf(conf),
	}
	client.Cid = len(center.Clients)
	center.Clients = append(center.Clients, client)
	client.Connect()
}

func (center *SSHCenter) NewConf(target string) *SSHOpts {
	opts := &SSHOpts{
		Conf: GetConf(target, "yaml", "."),
	}
	opts.Network = opts.Conf.GetString("opts.network")
	opts.User = opts.Conf.GetString("opts.user")
	opts.Addr = opts.Conf.GetString("opts.addr")
	opts.Port = opts.Conf.GetString("opts.port")
	opts.Kpath = opts.Conf.GetString("opts.kpath")
	opts.Key = opts.Conf.GetString("opts.key")
	return opts
}

func (center *SSHCenter) SetClient(cid int) {
	if cid >= len(center.Clients) && cid <= 0 {
		logrus.Warn("cid is illegal")
		return
	}
	store.Cid = cid
	center.Cid = cid
}

func (center *SSHCenter) SetSession(sid int) {
	if sid >= len(center.CurClient().Sessions) && sid <= 0 {
		logrus.Warn("sid is illegal")
		return
	}
	store.Sid = sid
	center.Sid = sid
}

func (center *SSHCenter) CurClient() *SSHClient {
	return center.Clients[center.Cid]
}

func (center *SSHCenter) CurSession() *SSHSession {
	return center.CurClient().Sessions[center.Sid]
}

func (center *SSHCenter) ClientIsValid() bool {
	return center.Cid > -1
}

func (center *SSHCenter) SessionIsValid() bool {
	return center.Sid > -1
}

func (center *SSHCenter) HasCommand() bool {
	return len(center.Command) != 0
}

func (center *SSHCenter) ClearCommand() {
	center.Command = ""
}

func (center *SSHCenter) SetCommand(cmd string) {
	center.Command = cmd
}

func (center *SSHCenter) FastExecCommand(cmd string) string {
	session, err := center.CurClient().Instance.NewSession()
	defer session.Close()
	data, err := session.CombinedOutput(cmd)
	PanicWithLogger(err, "Fail to exec command")
	return string(data)
}

func (center *SSHCenter) ReadWithStore() {
	for {
		if center.SessionIsValid() {
			center.CurSession().ReadWithStore()
		}
	}
}

func (center *SSHCenter) Write() {
	for {
		if center.SessionIsValid() && center.HasCommand() {
			center.CurSession().Write(center.Command)
			center.ClearCommand()
		}
	}
}

func NewSSHCenter() *SSHCenter {
	return &SSHCenter{
		Cid:     -1,
		Sid:     -1,
		Clients: make([]*SSHClient, 0),
	}
}
