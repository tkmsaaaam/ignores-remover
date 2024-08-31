package main

import (
	"bufio"
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
		return
	}

	fp, err := os.Open(request.path + request.file)
	if err != nil {
		log.Println("ignore file is not present", request.path+request.file)
	}
	defer fp.Close()

	log.Println("directory: ", request.path, "file: ", request.file)

	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		delete(request.path, scanner.Text(), dryRun)
	}
}

func getTarget() string {
	if len(os.Args) < 2 {
		return ""
	}
	return os.Args[1]
}

func isDryRun() bool {
	if len(os.Args) < 3 {
		return false
	}
	return strings.EqualFold(os.Args[2], "--dryRun=true")
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
			os.Remove(v)
			log.Println("removed: ", v)
		}
	}
}
