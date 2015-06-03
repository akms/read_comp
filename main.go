package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	app := cli.NewApp()
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} [options] [absolute path] [relative path]
   ex. read_comp ../../hoge.tar.gz
       read_comp /hoem/hoge/hoge.tar.gz

VERSION:
   {{.Version}}{{if or .Author .Email}}

AUTHOR:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
`
	app.Name = "read_comp"
	app.Usage = "read .tar.gz files and uncomppress files."
	app.Version = "2.0"
	app.Action = func(c *cli.Context) {
		if len(c.Args()) > 2 || len(c.Args()) < 2 {
			log.Fatal("too many or too little args.")
		}

		i := c.Args()
		readFile(i)
	}
	app.Run(os.Args)
}

func readFile(args []string) {
	var (
		buf         bytes.Buffer
		fileReader  io.ReadCloser
		file        *os.File
		file_name   string
		target_name string
		err         error
		hdr         *tar.Header
		dir_name    string
		fileinfo    os.FileInfo
		body        []byte
	)
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
	
	tr := tar.NewReader(fileReader)
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
			fmt.Println(hdr.Name)
			fileinfo = hdr.FileInfo()
			func() {
				mkdir_name, _ := filepath.Split(hdr.Name)
				changeDir(dd)

				if err = os.MkdirAll(mkdir_name, fileinfo.Mode()); err != nil {
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
