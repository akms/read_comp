package tarreader

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var (
	reader chan *tar.Reader = make(chan *tar.Reader)
	header chan *tar.Header = make(chan *tar.Header)
	body        []byte
	fileinfo    os.FileInfo
)

func unCompress() {
	tr := <-reader
	hdr := <-header
	fileinfo = hdr.FileInfo()
	func() {
		mkdir_name, _ := filepath.Split(hdr.Name)

		if err := os.MkdirAll(mkdir_name, fileinfo.Mode()); err != nil {
			log.Fatal(err)
		}
		if hdr.Typeflag == '0' {
			wfile, werr := os.Create(hdr.Name)
			if werr != nil {
				log.Fatal(werr)
			}
			defer wfile.Close()

			body = make([]byte, 8192)
			for {
				c, rerr := tr.Read(body)
				if c == 0 {
					break
				}
				if rerr != nil {
					log.Fatal(rerr)
				}
				wfile.Write(body[:c])
			}
		}
	}()
}

func Readfile(args []string) {
	var (
		buf         bytes.Buffer
		fileReader  io.ReadCloser
		file        *os.File
		file_name   string
		target_name string
		err         error
		hdr         *tar.Header
		dir_name    string
		tr          *tar.Reader
	)
	go unCompress()
	dd, _ := os.Getwd()
	dir_name, file_name = filepath.Split(args[0])
	changeDir(dir_name)
	target_name = args[1]
	default_Regexp := regexp.MustCompile(target_name)
	file, err = os.Open(file_name)
	if err != nil {
		log.Fatal("Can'ft open file \n")
	}
	defer file.Close()

	_, err = io.Copy(&buf, file)
	if fileReader, err = gzip.NewReader(&buf); err != nil {
		log.Fatal(err)
	}
	defer fileReader.Close()

	tr = tar.NewReader(fileReader)
	for {
		hdr, err = tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Can't read hdr \n")
			break
		}
		if default_Regexp.MatchString(hdr.Name) {
			changeDir(dd)
			fmt.Println(hdr.Name)
			reader <- tr
			header <- hdr
			changeDir(dir_name)
		}
	}
}

func changeDir(dir_name string) {
	var err error
	if dir_name != "" {
		err = os.Chdir(dir_name)
		if err != nil {
			log.Fatal(err)
		}
	}
}
