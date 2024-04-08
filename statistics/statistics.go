package statistics

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type SampleType int

const (
	Bigkeys SampleType = 0
	Memkeys SampleType = 1
	Hotkeys SampleType = 2
)

// GetSampleResult 执行系统命令并返回结果
func GetSampleResult(rh string, rp string, rpd string, st SampleType) (string, error) {
	sample_type := ""
	switch st {
	case Bigkeys:
		sample_type = "--bigkeys"
	case Hotkeys:
		sample_type = "--hotkeys"
	case Memkeys:
		sample_type = "--memkeys"
	default:
		sample_type = ""
	}
	os.Args = []string{"-h", rh, "-p", rp, "-a", rpd, sample_type}

	cmd := exec.Command("redis-cli", os.Args[0:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(fmt.Sprint(err) + ": " + stderr.String())
		//fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return "", err
	}
	return "Result: " + out.String(), nil
	/*fmt.Println("Result: " + out.String())*/

}

// ToResultFile 将结果写入指定文件
func ToResultFile(res string) error {
	// 获取当前工作目录的绝对路径
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}
	// 创建文件，如果文件已存在则会被截断为空
	file, err := os.OpenFile(dir+"/sampleResult", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	// 写入字符串到文件
	_, err = file.WriteString(res)
	if err != nil {
		return err
	}

	log.Println("写入文件成功")
	return nil

}
