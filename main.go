package main

import (
	"bytes"
	"fmt"
	"github.com/urfave/cli/v2"
	"gst/internal/app"
	"gst/internal/utils"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const MV = "/root/google-cloud-sdk/bin/gsutil mv %v%v gs://%v"

//var m = make(map[string]string)
var path, bucket, suffix string
var bytesN int64

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
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
		suffix = strings.ToLower(c.String("ext"))
		bytesN = utils.ParseSize(c.String("size"))
		bucket = c.Args().Get(1)
		mtime := c.Uint("time")
		var i uint = mtime
		fmt.Printf("path: %v bucket: %v suffix: %v size: %v  time: %v\n\n", path, bucket, suffix, bytesN, mtime)
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
						if file.Size() >= bytesN && strings.HasSuffix(strings.ToLower(file.Name()), suffix) {
							if isProcessing(file.Name()) {
								fmt.Printf("The same file name exists: %s\n", file.Name())
								continue
							}
							//fmt.Printf("%v Start to transfer files %v\n", time.Now().Format("2006-01-02 15:04:05"), path+file.Name())
							go workTrans(file.Name())
						}
					}
				}
				time.Sleep(time.Second * 3)
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
	//runtime.Goexit()
}

func getExecRet(cmdStr string) (result string) {
	cmd := exec.Command("sh", "-c", cmdStr)
	//out, err := exec.Command("sh", "-c", cmdStr).Output()
	//output, err := cmd.CombinedOutput()
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil && stderr.Len() > 0 {
		//fmt.Printf("err: %v\n", err)
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		//fmt.Printf("err len: %d\n", stderr.Len())
	}
	return fmt.Sprintf("%s", out.String())
}

func isProcessing(filename string) bool {
	ok := strings.Contains(getExecRet("ps -ef|grep 'gsutil'|grep -v 'grep'"), filename)
	return ok
}

func workTrans(filename string) {
	//x := fmt.Sprintf(MV, path, filename, bucket)
	//fmt.Println(x)
	//println("=>", filename)
	fmt.Printf("%v Start to transfer files %v\n", time.Now().Format("2006-01-02 15:04:05"), path+filename)
	cmd := exec.Command("/root/google-cloud-sdk/bin/gsutil", "mv", path+filename, "gs://"+bucket)
	//cmd := exec.Command("/usr/bin/python3", "/root/google-cloud-sdk/bin/bootstrapping/gsutil.py", "mv", path+filename, "gs://"+bucket)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	fmt.Printf("%s transfer Done,%s\n", filename, out.String())

	if err != nil && stderr.Len() > 0 {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}
	//fmt.Println(out.Len())
	//fmt.Printf("%s\n", out.String())
}
