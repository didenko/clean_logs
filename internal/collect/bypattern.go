package collect

import (
	"os"
	"path"
	"regexp"
)

func ByPattern(dir string, re_dir, re_file *regexp.Regexp, out chan string, err chan error) {
	defer close(out)
	scan_dir(dir, re_dir, re_file, out, err, false)
}

func scan_dir(current string, re_dir, re_file *regexp.Regexp, out chan string, err chan error, all bool) {

	f, e := os.Open(current)
	if e != nil {
		err <- e
		return
	}

	entries, e := f.Readdir(0)
	f.Close()

	if e != nil {
		err <- e
		return
	}

	for _, candidate := range entries {

		basename := candidate.Name()
		child := path.Join(current, basename)

		if candidate.IsDir() {

			if all || re_dir.MatchString(basename) {
				scan_dir(child, re_dir, re_file, out, err, true)
			} else {
				scan_dir(child, re_dir, re_file, out, err, false)
			}

		} else {
			if all || re_file.MatchString(basename) {
				out <- child
			}
		}
	}
}
