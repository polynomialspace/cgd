package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"net"
	"net/http"
	"net/http/cgi"
	"net/http/fcgi"
)

var cmd = flag.String("c", "", "CGI `prog`ram to run; relative paths are relative to -w dir")
var pwd = flag.String("w", "", "Working `dir` for CGI")
var serveFcgi = flag.Bool("f", false, "Serve FastCGI instead of HTTP")
var address = flag.String("a", ":3333", "Listen on `addr`ess")
var envVars = flag.String("e", "", "Comma-separated list of `env`ironment variables passed through to prog")

func main() {
	flag.Usage = func() {
		os.Stderr.WriteString("usage: cgd [-f] -c prog [-w wdir] [-a addr] [-e VAR1,VAR2]\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *cmd == "" {
		flag.Usage()
		os.Exit(2)
	}

	c, err := filepath.Abs(*cmd)
	if err != nil {
		log.Fatal(err)
	}

	// cgi already sets a default $PATH if empty
	// don't clobber that 'feature' here
	if env := os.Getenv("PATH"); env != "" {
		err = os.Setenv("PATH", env+":.")
		if err != nil {
			log.Fatal(err)
		}
	}

	h := &cgi.Handler{
		Path: c,
		Root: "/",
		Dir:  *pwd,
		InheritEnv: append(
			[]string{"PATH", "PLAN9"},
			strings.Split(*envVars, ",")...,
		),
	}

	l, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatal(err)
	}

	serve, msg := http.Serve, "HTTP server"
	if *serveFcgi {
		serve, msg = fcgi.Serve, "FastCGI daemon"
	}

	log.Println("Starting", msg, "listening on", *address)
	log.Fatal(serve(l, h))
}
