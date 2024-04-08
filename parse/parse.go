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
		fmt.Println("空文件已创建:", filePath)
	} else if err == nil {
		fmt.Println("文件已存在:", filePath)
	} else {
		return err
	}
	return nil
}

// ReadData 读取统计结果的文件
func ReadData() (string, error) {
	// 获取当前工作目录的绝对路径
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return "", err
	}
	filePath := dir + "/sampleResult"
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
func (b *BigKeyResult) ParseBigKeyResult(pType int8, data string) ([]*BigKeyResult, error) {
	// pType 为0时查找so far指标 为1时查找已经统计出的指标
	var bList []*BigKeyResult
	pattern := ``
	switch pType {
	case 0:
		pattern = `\w*\s*(\w*)\s*found so far\s*['](.*?)[']\s*with\s*(\d*)\s*(\w*)`
		break
	case 1:
		pattern = `\w*\s*(\w*)\s*found\s*['](.*?)[']\s*has\s*(\d*)\s*(\w*)`
		break
	default:
		pattern = ``

	}

	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(data, 100)

	for i := 0; i < len(matches); i++ {
		nB := new(BigKeyResult)
		nB.StructureType = matches[i][1]
		nB.KeyName = matches[i][2]
		nB.KeySize, _ = strconv.ParseFloat(matches[i][3], 64)
		nB.KeyUnit = matches[i][4]
		bList = append(bList, nB)

	}
	return bList, nil

}
