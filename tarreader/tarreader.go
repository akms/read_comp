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
	"runtime"
	"time"
)

var (
	reader chan *tar.Reader = make(chan *tar.Reader)
	header chan *tar.Header = make(chan *tar.Header)
)

func Uncompress() {
	var (
		chdr     *tar.Header
		ctr      *tar.Reader
		fileinfo os.FileInfo
		body     []byte
	)
	for {
		select {
		case chdr = <-header:
			fileinfo = chdr.FileInfo()
			mkdir_name, _ := filepath.Split(chdr.Name)
			if err := os.MkdirAll(mkdir_name, fileinfo.Mode()); err != nil {
				log.Fatal(err)
			}
		case ctr = <-reader:
			func() {
				if chdr.Typeflag == '0' {
					wfile, werr := os.Create(chdr.Name)
					if werr != nil {
						log.Fatal(werr)
					}
					defer wfile.Close()
					body = make([]byte, 512)
					for {
						c, rerr := ctr.Read(body)
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
	}
}

func readFile(tr *tar.Reader, dd string, dir_name string, default_Regexp *regexp.Regexp) {
	var (
		hdr *tar.Header
		err error
	)
	for {
		hdr, err = tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Can't read hdr \n ", err)
			break
		}
		if default_Regexp.MatchString(hdr.Name) {
			changeDir(dd)
			fmt.Println(hdr.Name)
			header <- hdr
			reader <- tr
			changeDir(dir_name)
			time.Sleep(time.Second / 8)
		}
	}
	close(header)
	close(reader)
}

func Readarchive(args []string) {
	var (
		rbuf        bytes.Buffer
		rbody       []byte
		fileReader  io.ReadCloser
		file        *os.File
		file_name   string
		target_name string
		err         error
		dir_name    string
		tr          *tar.Reader
	)
	go Uncompress()
	cpus := runtime.NumCPU()
	fmt.Println(cpus)
	runtime.GOMAXPROCS(cpus)
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
	rbody = make([]byte, 512)
	for {
		cr, err := file.Read(rbody)
		if cr == 0 {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		rbuf.Write(rbody[:cr])
	}
	if fileReader, err = gzip.NewReader(&rbuf); err != nil {
		log.Fatal(err)
	}
	defer fileReader.Close()

	tr = tar.NewReader(fileReader)
	readFile(tr, dd, dir_name, default_Regexp)
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
