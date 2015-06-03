package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
	"read_comp/tarreader"
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
		tarreader.Readfile(i)
	}
	app.Run(os.Args)
}
