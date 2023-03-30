package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	res = make(map[string][]string)
)

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

func calFileMD5(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	md5 := md5.New()
	io.Copy(md5, file)
	md5InBytes := md5.Sum(nil)[:16]
	return hex.EncodeToString(md5InBytes)
}

func main() {

	// root := "/Users/arning/arduino"
	files := make([]string, 0, 5)
	root := "/Users/arning/develop/code/my/golang/clamav-proxy"
	ListFile(root, &files)

	for _, v := range files {
		md5 := calFileMD5(v)
		if file, ok := res[md5]; ok {
			file = append(file, v)
			res[md5] = file
		} else {
			res[md5] = []string{v}
		}
		// fmt.Printf("md5: %v\n", md5)
	}

	for _, v := range res {
		if len(v) > 1 {
			for _, v2 := range v {
				fmt.Printf("v2: %v\n", v2)
			}
			fmt.Println("---------------")
		}
	}

	// list := []string{"sss", "aaa", "ccc", "bbb", "ddd", "bbb"}

	// for _, v := range list {
	// 	md5 := hex.EncodeToString(md5.New().Sum([]byte(v))[:16])
	// 	if e, ok := res[md5]; ok {
	// 		e = append(e, v)
	// 		res[md5] = e
	// 	} else {
	// 		// re := make([]string, 0)
	// 		// re = append(re, v)
	// 		res[md5] = []string{v}
	// 	}
	// }
	// // if v, ok := res["b"]; ok {
	// // 	fmt.Println("into ")
	// // 	v = append(v, "^^^^^")
	// // 	res["b"] = v
	// // }
	// fmt.Printf("res: %v\n", res)

}
