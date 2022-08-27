package repl

import "flag"

func init() {
	store = NewSSHStore()
	center = NewSSHCenter()

	exist = true
	bus = true

	flag.StringVar(&path, "path", "", "")
	flag.IntVar(&cid, "cid", -1, "")
	flag.IntVar(&sid, "sid", -1, "")
	flag.BoolVar(&session, "s", false, "")
}
