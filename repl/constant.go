package repl

const (
	NULL = iota
	ILLEGAL

	LACK

	NEW
	FAST
	NONE

	PROMPT    = ":> "
	HEADER    = "ac "
	NILSTRING = ""
)

var (
	center  *SSHCenter
	store   *SSHStore
	path    string
	cid     int
	sid     int
	session bool
	bus     bool
	exist   bool
)

type ACTYPE int
