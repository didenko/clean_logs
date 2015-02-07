package collect

import (
	"regexp"
	"testing"
)

func TestByPattern(t *testing.T) {
	var (
		out      = make(chan string, 10)
		err      = make(chan error, 3)
		outslice = make([]string, 0)
		errslice = make([]error, 0)
		re_file  = regexp.MustCompile("02$")
		re_dir   = regexp.MustCompile("02")

		expected = []string{
			"mock/010101/010102",
			"mock/010102",
			"mock/010104/010102",
			"mock/010104/010202",
			"mock/010203/010102",
			"mock/010203/010203a",
			"mock/010203/010203b",
		}
	)

	ByPattern("mock", re_dir, re_file, out, err)
	close(err)

	for f := range out {
		outslice = append(outslice, f)
	}

	for e := range err {
		errslice = append(errslice, e)
	}

	if differ(expected, outslice) {
		t.Error("Wrong collector results: ", outslice)
	}

	if len(errslice) > 0 {
		t.Error("Collector errored: ", errslice)
	}
}

func differ(a, b []string) bool {
	if len(a) != len(b) {
		return true
	}

	for i, _a := range a {
		if _a != b[i] {
			return true
		}
	}

	return false
}
