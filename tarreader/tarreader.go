package tarreader

import (
	"fmt"
)

func tarreader() {
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
