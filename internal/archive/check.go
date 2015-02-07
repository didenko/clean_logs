package archive

import (
	"archive/zip"
	"fmt"
	"os"
)

type Diff struct {
	asize int64
	dsize int64
}

func Check(arch string, files, done chan string, err chan error) {
	defer close(done)

	compare2disk := func(zf *zip.File) *Diff {

		asize := zf.FileInfo().Size()

		finfo, e := os.Stat(zf.Name)
		if e != nil {
			err <- e
			return &Diff{asize, 0}
		}

		dsize := finfo.Size()

		if dsize != asize {
			return &Diff{asize, dsize}
		}

		return nil
	}

	var (
		Missed    map[string]bool
		Extra     map[string]bool
		Different map[string]*Diff

		fexp  string
		fok   bool
		farch *zip.File
		aok   bool
	)

	fexp, fok = <-files

	areader, e := zip.OpenReader(arch)
	defer areader.Close()
	if e != nil {
		err <- e
	}
	afiles := slice2chan(areader.File)

	for {
		if !fok {
			break
		}

		farch, aok = <-afiles
		if !aok {
			break
		}

		if fexp == farch.Name {

			if d := compare2disk(farch); d != nil {
				Different[farch.Name] = d
			}

			fexp, fok = <-files
			continue
		}

		if Missed[farch.Name] {
			delete(Missed, farch.Name)

			if d := compare2disk(farch); d != nil {
				Different[farch.Name] = d
			}

		} else {
			Extra[farch.Name] = true
		}

		if Extra[fexp] {
			delete(Extra, fexp)

			if d := compare2disk(farch); d != nil {
				Different[farch.Name] = d
			}

		} else {
			Missed[fexp] = true
		}

		fexp, fok = <-files
	}

	for fexp = range files {
		Missed[fexp] = true
	}

	for farch = range afiles {
		Extra[farch.Name] = true
	}

	if len(Missed) > 0 {
		fmt.Println()
		fmt.Println("Missing from the archive:")
	}
	for mf, in := range Missed {
		if in {
			fmt.Println("\t", mf)
		}
	}

	if len(Extra) > 0 {
		fmt.Println()
		fmt.Println("Extra in the archive:")
	}
	for ef, in := range Extra {
		if in {
			fmt.Println("\t", ef)
		}
	}

	if len(Different) > 0 {
		fmt.Println()
		fmt.Println("Different from the archive {archive, disk size}:")
	}
	for df, sz := range Different {
		fmt.Println("\t", df, " ", sz)
	}
}

func slice2chan(s []*zip.File) chan *zip.File {
	ch := make(chan *zip.File)
	go func() {
		defer close(ch)
		for _, f := range s {
			ch <- f
		}
	}()
	return ch
}
