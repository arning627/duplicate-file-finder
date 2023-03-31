package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	wd, _ = os.Getwd()
	res   = make(map[string][]string)
	wg    sync.WaitGroup
)

type filePojo struct {
	hash string
	path string
}

func ListFile(root string, list *[]string) {
	rootDir, e := ioutil.ReadDir(root)
	if e != nil {
		log.Fatalln(e)
	}

	var builder strings.Builder
	for _, file := range rootDir {
		if file.IsDir() && file.Name() != root {
			ListFile(filepath.Join(root, file.Name()), list)
		} else {
			builder.WriteString(root)
			builder.WriteRune(os.PathSeparator)
			builder.WriteString(file.Name())
			*list = append(*list, builder.String())
			builder.Reset()
		}
	}
}

func calFileMD5(path string, ch chan<- filePojo) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	md5 := md5.New()
	io.Copy(md5, file)
	md5InBytes := md5.Sum(nil)[:16]
	hash := hex.EncodeToString(md5InBytes)
	obj := filePojo{hash, path}
	ch <- obj
	// return hash
}

func main() {
	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)
	log.Println("========start scan==========")
	log.Printf("CPU = %v", cpu)
	start_t := time.Now()
	fileHash := make(chan filePojo, 5)
	files := make([]string, 0, 5)
	var root = wd
	ListFile(root, &files)

	l := len(files)

	log.Printf("There are %v files in the current directory\n", l)
	wg.Add(cpu)
	// for _, v := range files {
	// 	go calFileMD5(v, fileHash)
	// }
	step := int(math.Ceil(float64(l / cpu)))
	for i := 0; i < cpu; i++ {
		index := step * i
		endIndex := index + step
		if i == cpu-1 {
			endIndex = l
		}
		go func() {
			defer wg.Done()
			for i := index; i < endIndex; i++ {
				calFileMD5(files[i], fileHash)
			}
		}()
	}

	// go func() {
	// 	for _, v := range files {
	// 		calFileMD5(v, fileHash)
	// 	}
	// }()

	go func() {
		for {
			hash := <-fileHash
			if file, ok := res[hash.hash]; ok {
				file = append(file, hash.path)
				res[hash.hash] = file
			} else {
				res[hash.hash] = []string{hash.path}
			}
		}
	}()

	f, _ := filepath.Abs("duplicate.txt")

	ff, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	defer ff.Close()

	writer := bufio.NewWriter(ff)

	wg.Wait()

	for _, v := range res {
		if len(v) > 1 {
			for _, v2 := range v {
				writer.WriteString(v2 + "\n")
			}
			writer.WriteString("-------------------\n")
			writer.Flush()
		}

	}
	end_t := time.Now()
	sub_t := end_t.Sub(start_t)
	log.Printf("total time : %v\n", sub_t)

}
