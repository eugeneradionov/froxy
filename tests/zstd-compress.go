package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/DataDog/zstd"
)

func main() { // nolint
	_, filename, _, _ := runtime.Caller(0) // nolint
	path := filepath.Dir(filename)

	rawDir := fmt.Sprintf("%s/test-files/raw", path)
	zstdDir := fmt.Sprintf("%s/test-files/zstd", path)

	files, err := ioutil.ReadDir(rawDir)
	if err != nil {
		log.Fatalf("failed to read dir %s: %v", rawDir, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := filepath.Join(rawDir, file.Name())

		f, err := os.Open(fileName)
		if err != nil {
			log.Printf("failed to open file %s: %v\n", fileName, err)
		}

		name, ext := fileNameWithExt(fileName)
		outF, err := os.OpenFile(fmt.Sprintf("%s/%s_zstd%s", zstdDir, name, ext), os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			log.Fatalf("failed to open file: %v", err)
		}

		w := zstd.NewWriter(outF)
		_, err = io.Copy(w, f)
		if err != nil {
			log.Fatalf("failed to compress file %s, %v\n", fileName, err)
		}

		err = w.Flush()
		if err != nil {
			log.Printf("failed to flush zstd writer: %v", err)
		}

		err = w.Close()
		if err != nil {
			log.Printf("failed to close zstd writer: %v", err)
		}

		err = outF.Close()
		if err != nil {
			log.Printf("failed to close file %s: %v\n", outF.Name(), err)
		}

		err = f.Close()
		if err != nil {
			log.Printf("failed to close file %s: %v\n", fileName, err)
		}
	}
}

func fileNameWithExt(fileName string) (string, string) { // nolint:gocritic
	base := filepath.Base(fileName)
	ext := filepath.Ext(fileName)

	return base[:len(base)-len(ext)], ext
}
