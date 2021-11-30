package see

import (
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/junbin-yang/golib/bytesconv"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type acc struct {
	obj     *logrus.Logger
	channel *dataContainer

	FileName string
	Path     string
	Rotate   bool
	KeepDays int64
}

func (this *acc) New() *acc {
	this.obj = logrus.New()
	this.obj.SetLevel(logrus.InfoLevel)

	if this.Path == "" {
		this.Path = "/var/log"
	}

	if this.KeepDays == 0 {
		this.KeepDays = 7
	}

	os.MkdirAll(this.Path, 0777)

	if this.FileName != "" {
		logfile := this.Path + "/" + this.FileName + "-" + carbon.Now().ToDateString() + ".log"
		src, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend|0644)
		if err == nil {
			this.obj.SetOutput(src)

			if this.Rotate {
				go func(Path, FileName string) {
					c := cron.New()
					c.AddFunc("0 0 0 * * ?", func() {
						logfile := Path + "/" + FileName + "-" + carbon.Now().ToDateString() + ".log"
						src, _ := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend|0644)
						this.obj.SetOutput(src)

						var diff_time int64 = 3600 * 24 * this.KeepDays
						now_time := time.Now().Unix()
						delpath := this.Path
						err := filepath.Walk(delpath, func(delpath string, f os.FileInfo, err error) error {
							if f == nil {
								return err
							}
							file_time := f.ModTime().Unix()
							if (now_time - file_time) > diff_time {
								if isFile(delpath) && strings.Contains(delpath, this.FileName) {
									os.RemoveAll(delpath)
								}
							}
							return nil
						})
						if err != nil {
							fmt.Printf("filepath.Walk() returned %s\r\n", err.Error())
						}
					})
					c.Start()
					select {}
				}(this.Path, this.FileName)
			}
		} else {
			fmt.Println("打开文件失败:", err.Error())
			this.writerStdout()
		}
	} else {
		this.writerStdout()
	}
	this.obj.SetFormatter(new(logFormatter))

	this.channel = initDataContainer()

	go func() {
		for {
			itemInterface := this.channel.Pop()
			if itemInterface != nil {
				item := itemInterface.(map[string]interface{})
				this.obj.WithFields(item).Info()
			}
		}
	}()

	return this
}

func (this *acc) writerStdout() {
	writers := []io.Writer{os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	this.obj.SetOutput(fileAndStdoutWriter)
}

func (this *acc) Writer(msg map[string]interface{}) {
	this.channel.Push(msg)
}

func (this *acc) PrintRoute(httpMethod, absolutePath string) {
	msg := map[string]interface{}{
		"Type":         "PrintRoute",
		"httpMethod":   httpMethod,
		"absolutePath": absolutePath,
	}
	this.channel.Push(msg)
}

func (this *acc) Println(o ...interface{}) {
	m := ""
	for _, v := range o {
		m += fmt.Sprint(v) + " "
	}

	msg := map[string]interface{}{
		"Type":    "Println",
		"Message": m,
	}
	this.channel.Push(msg)
}

//自定义格式
type logFormatter struct{}

func (s *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	t := carbon.Now().ToDateTimeString() + "," + fmt.Sprint(time.Now().Nanosecond()/1e6)
	msg := fmt.Sprintf("%s %s %s %s %s %d %d %s\n", t, entry.Data["userAgent"], entry.Data["clientIP"], entry.Data["reqMethod"], entry.Data["reqUri"], entry.Data["statusCode"], entry.Data["latencyTime"], entry.Data["Body"])
	if entry.Data["Type"] == "PrintRoute" {
		msg = fmt.Sprintf("%s %s %s %s\n", t, "Register Route:", entry.Data["httpMethod"], entry.Data["absolutePath"])
	}
	if entry.Data["Type"] == "Println" {
		msg = fmt.Sprintf("%s %s \n", t, entry.Data["Message"])
	}
	return bytesconv.StringToBytes(msg), nil
}

type dataContainer struct {
	queue chan interface{}
}

func initDataContainer() (dc *dataContainer) {
	dc = &dataContainer{}
	dc.queue = make(chan interface{})
	return dc
}

//非阻塞push
func (dc *dataContainer) Push(data interface{}) bool {
	click := time.After(time.Millisecond * 20)
	select {
	case dc.queue <- data:
		return true
	case <-click:
		return false
	}
}

//非阻塞pop
func (dc *dataContainer) Pop() (data interface{}) {
	click := time.After(time.Millisecond * 20)
	select {
	case data = <-dc.queue:
		return data
	case <-click:
		return nil
	}
}

func isFile(path string) bool {
	fi, err := os.Stat(path)
	if nil == err {
		return !fi.IsDir()
	} else {
		return false
	}
}
