package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	chia()
}

// 复制P图文件
func chia() {
	/* 用法
	.\tci.exe --src C:/Users/wang10k/Downloads/bbb --dest c:/chia --all x:/chia,y:/chia,z:/chia --count 2
	*/
	src := flag.String("src", "", "请输入本机P图目录")
	dest := flag.String("dest", "", "请输入复制目的地目录")
	all := flag.String("all", "", "请输入所有复制目的地目录")
	counter := flag.Int("count", 2, "同时存在的标记文件数量")
	flag.Parse() // 解析参数

	// net use Z: \\share\xxx /persistent:yes
	// src = "C:/Users/wang10k/Downloads/bbb"
	// desc := "C:/chia"
	// x1, err := os.Stat(path)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	if matchFile(*dest) {
		log.Println("文件存在,pass")
	} else {
		if matchMarkFileCount(*all, *counter) {
			createMarkFile(*dest)
			copyPlots(*src, *dest)
			delMarkFile(*dest)
		} else {
			log.Println("存储并行数量达到上限")
		}
	}
}

// copyPlots 复制P图文件
func copyPlots(srcDir, destDir string) {
	srcFiles, err := ioutil.ReadDir(srcDir)
	if err != nil {
		log.Fatal("list srcDir file error", err)
	}
	destFiles, err := ioutil.ReadDir(destDir)
	if err != nil {
		log.Fatal("list destDir file error", err)
	}
	srcTmpFiles := make(map[string]bool)
	for _, v1 := range srcFiles {
		srcTmpFiles[v1.Name()] = true
	}
	for _, v1 := range srcFiles {
		for _, v2 := range destFiles {
			if v1.Name() == v2.Name() {
				delete(srcTmpFiles, v1.Name())
			}
		}
	}
	for k := range srcTmpFiles {
		// 注意这里传参需要2个文件而不是目录
		copyFlag := false
		srcFilePath := fmt.Sprintf("%s/%s", srcDir, k)
		destFilePath := fmt.Sprintf("%s/%s", destDir, k)
		if err := Copy(srcFilePath, destFilePath, 1024*1024*100); err != nil {
			log.Fatal("复制文件失败", err)
		} else {
			log.Println("复制文件成功", k)
			log.Println("开始删除文件", k)
			delSrcFile(srcFilePath)
			copyFlag = true
		}
		if copyFlag {
			break
		}
	}
}

// matchMarkFileCount 判断标记文件数量
func matchMarkFileCount(path string, cLimit int) bool {
	allPath := strings.Split(path, ",")
	allFiles := []os.FileInfo{}
	for _, v := range allPath {
		files, err := ioutil.ReadDir(v)
		if err != nil {
			log.Fatal("list file error", err)
		}
		allFiles = append(allFiles, files...)
	}

	counter := 0
	for _, v := range allFiles {
		if strings.HasPrefix(v.Name(), "DSY") {
			counter++
		}
	}
	return counter < cLimit
}

// matchFile 判断文件是否存在
func matchFile(path string) bool {
	// allPath := strings.Split(path, ",")
	// allFiles := []os.FileInfo{}
	// for _, v := range allPath {

	// }
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal("list file error", err)
	}
	// allFiles = append(allFiles, files...)
	mark := false
	for _, v := range files {
		// fmt.Println(k, v.Name(), v.Size()/1024/1024, "M")
		if v.Name() == getHostName() {
			mark = true
			break
		}
	}
	return mark
}

// createMarkFile 创建标记文件
func createMarkFile(path string) {
	filePath := fmt.Sprintf("%s/%s", path, getHostName())
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("创建标记文件", getHostName())
	}
	defer f.Close()
}

// 删除标记文件
func delMarkFile(path string) {
	filePath := fmt.Sprintf("%s/%s", path, getHostName())
	err := os.Remove(filePath)
	if err != nil {
		log.Fatal("删除标记文件失败")
	} else {
		log.Println("删除标记文件成功")
	}
}

// 删除标记文件
func delSrcFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Fatalf("删除%s失败", path)
	} else {
		log.Printf("删除%s成功\n", path)
	}
}

// getHostName 获取主机名
func getHostName() string {
	name, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return name
}

// Copy 复制文件
func Copy(src, dst string, BUFFERSIZE int64) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file.", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()
	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("File %s already exists.", dst)
	}
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	if err != nil {
		panic(err)
	}
	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}
