// Package dagvis provides clear visualization of a DAG.
package dagvis

import (
	"flag"
	"log"

	"shanhu.io/aries"
	"shanhu.io/misc/errcode"
	"shanhu.io/misc/osutil"
)

type server struct {
}

func newServer(h *osutil.Home) *server {
	return new(server)
}

func (s *server) serveIndex(c *aries.C) error {
	return nil
}

func makeService(home string) (aries.Service, error) {
	h, err := osutil.NewHome(home)
	if err != nil {
		return nil, errcode.Annotate(err, "make new home")
	}

	s := newServer(h)

	r := aries.NewRouter()
	r.Index(s.serveIndex)
	return r, nil
}

// Main is main.
func Main() {
	addr := aries.DeclareAddrFlag("")
	home := flag.String("home", ".", "home dir")
	flag.Parse()

	s, err := makeService(*home)
	if err != nil {
		log.Fatal(err)
	}
	if err := aries.ListenAndServe(*addr, s); err != nil {
		log.Fatal(err)
	}
}
