package parse

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func checkAndCreateFile(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// 文件不存在，创建空文件
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(file)
		fmt.Println("Empty file created:", filePath)
	} else if err == nil {
		fmt.Println("File already exists:", filePath)
	} else {
		return err
	}
	return nil
}

// ReadData 读取统计结果的文件
func ReadData(fileName string) (string, error) {
	// 获取当前工作目录的绝对路径
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return "", err
	}
	filePath := dir + fmt.Sprintf("/%s_sampleResult", fileName)
	err = checkAndCreateFile(filePath)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil

}

// ParseBigKeyResult 解析统计结果，结果写入BigKeyResult结构体指针列表并返回
func (b *BigKeyResult) ParseBigKeyResult(data string, stFlags string) ([]*BigKeyResult, error) {
	var bList []*BigKeyResult
	pattern := ``
	switch stFlags {
	case "big":
		pattern = `\w*\s*(\w*)\s*found\s*['](.*?)[']\s*has\s*(\d*)\s*(\w*)`
		break
	case "mem":
		pattern = `\w*\s*(\w*)\s*found\s*['](.*?)[']\s*has\s*(\d*)\s*(\w*)`
		break
	case "hot":
		pattern = `\w*\s*\w*\s*found\s*with\s*\w*:\s*(\d*)\s*keyname:\s*["](.*?)["]`
		break
	}

	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(data, 100)

	for i := 0; i < len(matches); i++ {
		nB := new(BigKeyResult)
		if stFlags == "big" || stFlags == "mem" {
			nB.StructureType = matches[i][1]
			nB.KeyName = matches[i][2]
			nB.KeySize, _ = strconv.ParseFloat(matches[i][3], 64)
			nB.KeyUnit = matches[i][4]
		} else if stFlags == "hot" {
			nB.StructureType = "counter"
			nB.KeyName = matches[i][2]
			nB.KeySize, _ = strconv.ParseFloat(matches[i][1], 64)
			nB.KeyUnit = "times"
		}
		nB.SampleType = fmt.Sprintf("%skeys", stFlags)
		bList = append(bList, nB)
	}
	return bList, nil

}
