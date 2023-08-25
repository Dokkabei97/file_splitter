package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	outputFileBase = "split_log"
	size           = 500 * 1024 * 1024
)

func main() {
	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var currentSize int64 = 0
	var fileCount int = 1

	outputFile, err := os.Create(fmt.Sprintf("%s_%d.log", outputFileBase, fileCount))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	writer := bufio.NewWriter(outputFile)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("[Error] 파일 읽기 실패: %s\n", err)
			return
		}

		_, err = writer.Write(line)
		if err != nil {
			fmt.Printf("[Error] 파일 쓰기 실패: %s\n", err)
			return
		}
		currentSize += int64(len(line))

		if currentSize >= size {
			writer.Flush()
			outputFile.Close()

			fileCount++
			outputFile, err = os.Create(fmt.Sprintf("%s_%d.log", outputFileBase, fileCount))
			if err != nil {
				fmt.Printf("[Error] 파일 생성 실패: %s\n", err)
				return
			}
			writer = bufio.NewWriter(outputFile)
			currentSize = 0
		}
	}

	writer.Flush()
	outputFile.Close()
	fmt.Printf("파일 분할 완료! 총 %d개의 파일이 생성되었습니다.\n", fileCount)
}
