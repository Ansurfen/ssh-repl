package repl

import (
	"bufio"
	"os"
	"strings"
)

func Parse(cmd string) string {
	var actype ACTYPE
	cmd, actype = ParseHelper(cmd)
	switch actype {
	case NULL:
		if !bus {
			center.SetCommand(cmd)
			return NILSTRING
		}
		println()
		return NILSTRING
	case ILLEGAL:
		println("Fail to find the header of ac or args is less than 1.")
		return NILSTRING
	case NEW:
		println("Preparing to create a client connet...")
		center.NewClient(path)
		return NILSTRING
	case FAST:
		println("Preparing to create a client connet and session...")
		center.NewClient(path)
		center.SetClient(0)
		center.CurClient().NewSession()
		center.SetSession(0)
		bus = false
		return NILSTRING
	case NONE:
		return cmd
	}
	resetFlags()
	return cmd
}

func ParseHelper(cmd string) (string, ACTYPE) {
	cmd = Strip(strings.TrimSpace(cmd))
	if len(cmd) == 0 {
		return NILSTRING, NULL
	}
	if cmd == ":quit" {
		center.SetCommand("clear")
		center.Sid = -1
		bus = true
		return NILSTRING, NULL
	}
	if bus && cmd == "quit" {
		exist = false
		return NILSTRING, NULL
	}
	if !bus {
		return cmd, NULL
	}
	if !strings.HasPrefix(cmd, HEADER) {
		return NILSTRING, ILLEGAL
	}
	cmd = cmd[len(HEADER):]
	args := strings.Split(cmd, " ")
	if len(args) <= 0 {
		return NILSTRING, LACK
	}
	var err error
	switch strings.ToLower(args[0]) {
	case "ls":
		if len(center.Clients) == 0 {
			println("no available client")
			return NILSTRING, NULL
		}
		print("available client: ")
		for c := range center.Clients {
			print(c, " ")
		}
		println()
	case "new":
		if args, err = nextArg(args); err != nil {
			return NILSTRING, LACK
		}
		parseFlags(args)
		if center.ClientIsValid() && session {
			center.CurClient().NewSession()
		}
		if len(path) != 0 {
			return NILSTRING, NEW
		}
		println("Fail to parse path, try ac new -path [confName]")
		return NILSTRING, NULL
	case "set":
		if args, err = nextArg(args); err != nil {
			return NILSTRING, LACK
		}
		parseFlags(args)
		if sid == -1 && cid == -1 {
			return NILSTRING, NULL
		}
		if cid >= 0 {
			center.SetClient(cid)
		}
		if sid != -1 && cid == -1 && !center.ClientIsValid() {
			println("please set cid")
		} else if sid != -1 && center.ClientIsValid() {
			center.SetSession(sid)
		}
		return NILSTRING, NULL
	case "cur":
		println("Cid: ", center.Cid, " Sid: ", center.Sid)
		return NILSTRING, NULL
	case "fast":
		if args, err = nextArg(args); err != nil {
			return NILSTRING, LACK
		}
		parseFlags(args)
		if len(path) != 0 {
			return NILSTRING, FAST
		}
		return NILSTRING, NULL
	}
	return cmd, NONE
}

func Launch() {
	go center.ReadWithStore()
	go center.Write()
	buf := bufio.NewReader(os.Stdin)
	for exist {
		if bus {
			print(":> ")
		}
		cmd, _ := buf.ReadString('\n')
		if cmd = Parse(cmd); len(cmd) == 0 {
			continue
		}
	}
	store.Save()
}
