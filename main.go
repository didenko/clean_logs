package main

import (
	"flag"
	"github.com/didenko/clean_logs/internal/archive"
	"github.com/didenko/clean_logs/internal/collect"
	"github.com/didenko/clean_logs/internal/flow"
	"log"
	"os"
	"os/user"
	"path"
	"regexp"
	"time"
)

var re_file, re_dir *regexp.Regexp
var basedir, month, store string
var problems bool = false

var (
	files         = make(chan string, 1000)
	filesToArch   = make(chan string, 1000)
	filesToDisk   = make(chan string, 1000)
	filesFromDisk = make(chan string, 1000)
	done_archv    = make(chan string)
	done_check    = make(chan string)
	err           = make(chan error)
)

func init() {
	usr, e := user.Current()
	if e != nil {
		log.Fatal(e)
	}

	flag.StringVar(&month, "month", time.Now().AddDate(0, -1, 0).Format("200601"), "in the YYYYMM format")
	flag.StringVar(&store, "store", usr.HomeDir, "Default storage location")
	flag.StringVar(&basedir, "base", ".", "Path to the directory to scan")

	flag.Parse()

	re_file = regexp.MustCompile("_" + month + "\\d\\d.*\\.log$")
	re_dir = regexp.MustCompile("^" + month + "\\d\\d$")
	store = path.Join(store, "log_"+month+".zip")
}

func main() {
	defer close(err)

	go func() {
		for e := range err {
			log.Println(e)
			problems = true
		}
	}()

	e := os.Chdir(basedir)
	if e != nil {
		err <- e
	} else {
		go archive.Check(store, filesFromDisk, done_check, err)
		go flow.Stash(filesToDisk, done_archv, filesFromDisk, err)
		go archive.Add(store, basedir, filesToArch, done_archv, err)
		go flow.Fork(files, filesToDisk, filesToArch)
		go collect.ByPattern(".", re_dir, re_file, files, err)
	}

	<-done_check
	if problems {
		os.Exit(1)
	}
}
