package statistics

import (
	"bytes"
	"fmt"
	"github.com/oklog/run"
	"log"
	"os"
	"os/exec"
)

// GetSampleResultConcurrent Execute system commands concurrently and return the result dictionary set
func GetSampleResultConcurrent(rh string, rp string, rpd string, stList []string) map[string]string {
	var g run.Group
	result := make(map[string]string)
	stFlags := make([]string, len(stList))
	for i, eachFlag := range stList {
		switch eachFlag {
		case "big":
			stFlags[i] = "--bigkeys"
			continue
		case "hot":
			stFlags[i] = "--hotkeys"
			continue
		case "mem":
			stFlags[i] = "--memkeys"
			continue
		}
	}
	for _, each := range stFlags {
		eFlag := each
		g.Add(func() error {
			args := []string{"-h", rh, "-p", rp, "-a", rpd, eFlag}
			cmd := exec.Command("redis-cli", args...)
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			log.Println("executing command:" + cmd.String())
			err := cmd.Run()
			if err != nil {
				gError := fmt.Sprintf("error executing command: %v", err)
				log.Println("error executing command:" + cmd.String())
				log.Println(gError)
				return fmt.Errorf(gError)

			} else {
				result[eFlag[2:5]] = "Result: " + out.String()
			}
			return nil

		},

			func(err error) {
				log.Println("execute cmd quite:" + eFlag[2:])
				return
			},
		)
	}
	err := g.Run()
	if err != nil {
		log.Printf("Error running GetSampleResultConcurrent: %v\n", err)
	}
	log.Printf("result: %s", result)
	return result

}

// GetSampleResultSerial Execute system commands serially and return the result dictionary set
func GetSampleResultSerial(rh string, rp string, rpd string, stList []string) map[string]string {
	result := make(map[string]string)
	stFlags := make([]string, len(stList))
	for i, eachFlag := range stList {
		switch eachFlag {
		case "big":
			stFlags[i] = "--bigkeys"
			continue
		case "hot":
			stFlags[i] = "--hotkeys"
			continue
		case "mem":
			stFlags[i] = "--memkeys"
			continue
		}
	}
	for _, eFlag := range stFlags {
		var out bytes.Buffer
		var stderr bytes.Buffer
		args := []string{"-h", rh, "-p", rp, "-a", rpd, eFlag}
		cmd := exec.Command("redis-cli", args...)
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		log.Println("executing command:" + cmd.String())
		if err != nil {
			log.Println("error executing command:" + cmd.String())
			log.Println(fmt.Sprint(err) + ": " + stderr.String())
			return nil

		} else {
			result[eFlag[2:5]] = "Result: " + out.String()
		}

	}
	log.Printf("result: %s", result)
	return result
}

// ToResultFile Write the results to the specified file
func ToResultFile(res string, fileName string) error {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}
	file, err := os.OpenFile(dir+fmt.Sprintf("/%s_sampleResult", fileName), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	_, err = file.WriteString(res)
	if err != nil {
		return err
	}

	log.Println("Write the results to file success!")
	return nil

}
