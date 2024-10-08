package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Request struct {
	path string
	file string
}

var osGetwd = os.Getwd
var osStat = os.Stat

func main() {
	target := getTarget()
	dryRun := isDryRun()

	request := makeRequest(target)
	if request == nil {
		log.Println("can not continue to process")
		return
	}

	fp, err := os.Open(request.path + request.file)
	if err != nil {
		log.Println("Can not open file", request.path+request.file, err)
	}
	defer fp.Close()

	log.Println("directory: ", request.path, "file: ", request.file)

	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		delete(request.path, scanner.Text(), dryRun)
	}
}

var t string
var target string
var d bool
var dryRun bool

func init() {
	flag.StringVar(&t, "t", "", "dryRun")
	flag.StringVar(&target, "target", "", "dryRun")
	flag.BoolVar(&d, "d", false, "dryRun")
	flag.BoolVar(&dryRun, "dryRun", false, "dryRun")
}

func getTarget() string {
	if t != "" {
		return t
	}
	return target
}

func isDryRun() bool {
	return d || dryRun
}

func makeRequest(arg string) *Request {
	if arg == "" {
		dir, e := osGetwd()
		if e != nil {
			log.Println("can not get working directory: ", e)
			return nil
		}
		return &Request{path: dir + "/", file: ".gitignore"}
	} else {
		stat, e := osStat(arg)
		if e != nil {
			log.Println("can not discriminate arg (directory or file): ", e)
			return nil
		}
		if stat.IsDir() {
			if strings.HasSuffix(arg, "/") {
				return &Request{path: arg, file: ".gitignore"}
			} else {
				return &Request{path: arg + "/", file: ".gitignore"}
			}
		} else {
			l := strings.Split(arg, "/")
			if len(l) == 1 {
				dir, e := osGetwd()
				if e != nil {
					log.Println("can not get working directory: ", e)
					return nil
				}
				return &Request{path: dir + "/", file: l[0]}
			} else {
				fileName := l[len(l)-1]
				path := arg[0 : len(arg)-len(fileName)]
				return &Request{path: path, file: fileName}
			}
		}
	}
}

func delete(path, pattern string, dryRun bool) {
	files, err := filepath.Glob(path + pattern)
	if err != nil {
		log.Println("can not get: ", err)
		return
	}
	for _, v := range files {
		if dryRun {
			log.Println("remove(without dryRun): ", v)
		} else {
			e := os.Remove(v)
			if e != nil {
				log.Println("Can not remove", e)
			} else {
				log.Println("removed: ", v)
			}
		}
	}
}
