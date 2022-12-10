package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Not a file")
	}

	// 引数からファイル名を取得する
	sjisFile, err := os.Open(os.Args[1])
	if err != nil {
		panic(1)
	}
	defer sjisFile.Close()

	// ShiftJISのデコーダーを噛ませたReaderを作成し読み込む
	reader := transform.NewReader(sjisFile, japanese.ShiftJIS.NewDecoder())
	ret, err := io.ReadAll(reader)
	if err != nil {
		panic(1)
	}

	_, inFileName := filepath.Split(sjisFile.Name())
	outFile, err := os.Create(inFileName + ".acc")
	if err != nil {
		panic(1)
	}
	writer := transform.NewWriter(outFile, japanese.ShiftJIS.NewEncoder())
	defer outFile.Close()

	//
	fileStr := strings.Split(string(ret), "\r\n")
	fmt.Println(fileStr)

	re1 := regexp.MustCompile(`\s+`)

	for i := 0; i < len(fileStr)-1; i++ {
		if i < 2 {
			writer.Write([]byte(fileStr[i]))
		} else {
			lineTrimmed := strings.TrimSpace(fileStr[i])
			lineProcessed := re1.ReplaceAllString(lineTrimmed, `\s`)
			lineReturned := strings.ReplaceAll(lineProcessed, `\s`, "\r\n")
			writer.Write([]byte(lineReturned))
		}

		if i != len(fileStr)-2 {
			writer.Write([]byte("\r\n"))
		}
	}
}
