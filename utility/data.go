package utility

import (
	"io"
	"os"
)

// Get env.txt data. Each data is seperated by '\n'
func GetEnvData() []string {
	var envData []string

	file, _ := os.Open("env.txt")
	defer file.Close()

	buf := make([]byte, 128)
	cnt, err := file.Read(buf)
	if err != nil && err != io.EOF {
		panic(err)
	}

	var data string
	for i := 0; i < cnt; i++ {
		if buf[i] > 32 {
			data += string(buf[i])
		} else if buf[i] == 10{
			envData = append(envData, data)
			data = ""
		}
	}
	envData = append(envData, data)

	return envData
}

// Update env.txt with updatedEnvData
func UpdateEnvData(updatedEnvData []string) error {
	buf := make([]byte, 128)
	i := 0
	for _, data := range updatedEnvData {
		for j := 0; j<len(data); j++ {
			buf[i] = data[j]
			i++
		}
		buf[i] = 10
		i++
	}

	file, err := os.Create("env.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(buf[:i])
	return err
}