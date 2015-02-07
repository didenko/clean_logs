package flow

import (
	"bufio"
	"io/ioutil"
	"os"
)

func Stash(in, done, out chan string, err chan error) {

	defer close(out)

	fd, e := ioutil.TempFile("", "stash_")

	defer func(fname string) {
		os.Remove(fname)
	}(fd.Name())

	defer fd.Close()

	if e != nil {
		err <- e
		return
	}

	for line := range in {
		_, e = fd.WriteString(line + "\n")
		if e != nil {
			err <- e
			return
		}
	}

	<-done

	_, e = fd.Seek(0, 0)
	if e != nil {
		err <- e
		return
	}

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		out <- scanner.Text()
	}

	if e := scanner.Err(); e != nil {
		err <- e
	}
}

func Fork(in, outA, outB chan string) {
	defer close(outA)
	defer close(outB)

	for line := range in {
		outA <- line
		outB <- line
	}
}
