package archive

import (
	"archive/zip"
	"io"
	"os"
	"path"
)

func Add(arch, basedir string, fnames, done chan string, err chan error) {

	defer close(done)

	zip_file, e := os.Create(arch)
	defer zip_file.Close()
	if e != nil {
		err <- e
		return
	}

	zip_writer := zip.NewWriter(zip_file)

	dest := zip_dest{zip_writer, basedir, err}

	for file := range fnames {
		dest.push(file)
	}

	e = zip_writer.Close()
	if e != nil {
		err <- e
	}

}

type zip_dest struct {
	writer  *zip.Writer
	basedir string
	err     chan error
}

func (zd *zip_dest) push(file string) {

	file_reader, e := os.Open(path.Join(zd.basedir, file))
	if e != nil {
		zd.err <- e
		return
	}
	defer file_reader.Close()

	zip_member, e := zd.writer.Create(file)
	if e != nil {
		zd.err <- e
		return
	}

	io.Copy(zip_member, file_reader)
}
