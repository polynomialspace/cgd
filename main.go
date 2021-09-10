package main

import (
	"flag"
	"log"
	"os"
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

	h := &cgi.Handler{
		Path: *cmd,
		Root: "/",
		Dir:  *pwd,
		InheritEnv: append(
			[]string{"PATH", "PLAN9"},
			strings.Split(*envVars, ",")...,
		),
	}

	// This is a hack to make p9p's rc happier for some unknown reason.
	if h.Path[0] != '/' && strings.Split(h.Path, "/")[0][0] != '.' {
		h.Path = "./" + h.Path
	}

	err := os.Setenv("PATH", os.Getenv("PATH")+":.")
	if err != nil {
		log.Fatal(err)
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

