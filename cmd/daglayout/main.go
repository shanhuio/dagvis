// Copyright (C) 2021  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"shanhu.io/dags"
	"shanhu.io/misc/errcode"
)

type options struct {
	reverse bool
}

func readInput(in string) ([]byte, error) {
	if in == "" {
		return io.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(in)
}

func writeOutput(out string, bs []byte) error {
	if out == "" {
		_, err := os.Stdout.Write(bs)
		return err
	}
	return ioutil.WriteFile(out, bs, 0644)
}

func layoutGraph(g *dags.Graph, opts *options) (
	*dags.Map, *dags.MapView, error,
) {
	if opts.reverse {
		return dags.RevLayout(g)
	}
	return dags.Layout(g)
}

func layout(in, out string, opts *options) error {
	bs, err := readInput(in)
	if err != nil {
		return errcode.Annotate(err, "read input")
	}

	g := new(dags.Graph)
	if err := json.Unmarshal(bs, g); err != nil {
		return errcode.Annotate(err, "parse graph")
	}

	_, v, err := layoutGraph(g, opts)
	if err != nil {
		return errcode.Annotate(err, "layout graph")
	}

	m := dags.Output(v)
	outBytes, err := json.Marshal(m)
	if err != nil {
		return errcode.Annotate(err, "encode output")
	}

	if err := writeOutput(out, outBytes); err != nil {
		return errcode.Annotate(err, "write output")
	}

	return nil
}

func main() {
	reverse := flag.Bool("reverse", false, "if use reverse layout")
	in := flag.String("in", "", "input file")
	out := flag.String("out", "", "output file")
	flag.Parse()

	opts := &options{
		reverse: *reverse,
	}

	if err := layout(*in, *out, opts); err != nil {
		log.Fatal(err)
	}
}
