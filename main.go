package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gst/internal/app"
	"gst/internal/utils"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const MV = "gsutil mv %v%v gs://%v"

//var m = make(map[string]string)
var path, bucket, suffix string
var bytes int64

func main() {
	app := app.New()
	app.Action = func(c *cli.Context) (err error) {
		if c.NArg() != 2 {
			fmt.Fprintf(os.Stderr, "Error: %s takes arguments error.\n\n", app.Name)
			cli.ShowAppHelp(c)
			os.Exit(1)
		}

		if !strings.HasSuffix(c.Args().First(), string(os.PathSeparator)) {
			fmt.Printf("Wrong path.\n\n")
			cli.ShowAppHelp(c)
			os.Exit(1)
		}

		path, _ = filepath.Split(c.Args().First())
		//path, _ = filepath.Split("/root/qq/")
		suffix = strings.ToUpper(c.String("ext"))
		bytes = utils.ParseSize(c.String("size"))
		bucket = c.Args().Get(1)
		mtime := c.Uint("time")
		var i uint = mtime
		fmt.Printf("path: %v bucket: %v suffix: %v size: %v  time: %v\n\n", path, bucket, suffix, bytes, mtime)
		//fmt.Print(getExecRet("ps -ef"))
		for {
			if i < mtime {
				fmt.Printf("Perform the next scan after %d minutes.\n", mtime-i)
				i++
				time.Sleep(time.Minute)
				continue
			}

			if files, err := ioutil.ReadDir(path); err == nil {
				fmt.Printf("%v Scan folder %v\n", time.Now().Format("2006-01-02 15:04:05"), path)

				for _, file := range files {
					if !file.IsDir() {
						if file.Size() >= bytes && strings.HasSuffix(strings.ToUpper(file.Name()), suffix) {
							if isProcessing(file.Name()) {
								continue
							}
							//fmt.Printf("%v Start to transfer files %v\n", time.Now().Format("2006-01-02 15:04:05"), path+file.Name())
							go work(file.Name())
						}
					}
				}
				fmt.Printf("Perform the next scan after %d minutes.\n", mtime)
				i = 1
				time.Sleep(time.Minute)
			}
		}

		return
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getExecRet(cmdStr string) (result string) {
	out, err := exec.Command("sh", "-c", cmdStr).Output()
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s", out)
}

func isProcessing(filename string) bool {
	ok := strings.Contains(getExecRet("ps -ef|grep 'gsutil'|grep -v 'grep'"), filename)
	return ok
}

func work(filename string) {
	x := fmt.Sprintf(MV, path, filename, bucket)

	//println("=>", filename)
	fmt.Printf("%v Start to transfer files %v\n", time.Now().Format("2006-01-02 15:04:05"), path+filename)
	cmd := exec.Command("sh", "-c", x)
	if err := cmd.Run(); err != nil {
		println(err.Error())
	}
}
