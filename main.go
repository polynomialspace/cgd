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

var cmd = flag.String("c", "", "CGI `prog`ram to run")
var pwd = flag.String("w", "", "Working `dir` for CGI")
var serveFcgi = flag.Bool("f", false, "Serve FastCGI instead of HTTP")
var address = flag.String("a", ":3333", "Listen on `addr`ess")
var envVars = flag.String("e", "", "Comma-separated list of `env`ironment variables passed through to prog")

func main() {
	flag.Usage = func () {
		os.Stderr.WriteString("usage: cgd [-f] -c prog [-w wdir] [-a addr] [-e VAR1,VAR2]\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *cmd == "" {
		flag.Usage()
		os.Exit(2)
	}

	// This is a hack to make p9p's rc happier for some unknown reason.
	c := *cmd
	if c[0] != '/' {
		c = "./" + c
	}

	os.Setenv("PATH", os.Getenv("PATH")+":.")

	envList := []string{"PATH", "PLAN9"}
	for _, envVar := range strings.Split(*envVars, ",") {
		envList = append(envList, envVar)
	}

	h := &cgi.Handler{
		Path:       c,
		Root:       "/",
		Dir:        *pwd,
		InheritEnv: envList,
	}

	var err error
	if *serveFcgi {
		if l, err := net.Listen("tcp", *address); err == nil {
			log.Println("Starting FastCGI daemon listening on", *address)
			err = fcgi.Serve(l, h)
		}

	} else {
		log.Println("Starting HTTP server listening on", *address)
		err = http.ListenAndServe(*address, h)
	}

	if err != nil {
		log.Fatal(err)
	}
}

