package wavelog

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*
   .go : 自定义的一些文件操作
*/

var DefaultPerm = os.ModeDir | os.ModePerm

// CUT_FILE_BY_DATE_TIME_FORMAT 按日期切割文件的名称时间格式
const CUT_FILE_BY_DATE_TIME_FORMAT = "200601021504"

// CUT_FILE_BY_SIZE_TIME_FORMAT 按文件大小切割文件的名称时间格式
const CUT_FILE_BY_SIZE_TIME_FORMAT = "20060102150405"

// LOG_FILE_MAX_SIZE 文件超过多大时，切割文件
const LOG_FILE_MAX_SIZE = 1 * 1024 // 1 * 1024 * 1024  //1m
// CreateLogFile 创建文件
func CreateLogFile(filePath string) error {
	// 检查目录是否存在，不存在则创建
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	// 创建日志文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

// MkdirFile 创建目录
func MkdirFile(dirPath string, perm os.FileMode) (bool bool, errInfo error) {
	if perm == 0 {
		perm = DefaultPerm
	}
	err := os.Mkdir(dirPath, perm)
	if err != nil {
		return false, err
	}
	return true, errInfo
}

//MkdirAll 创建完整目录

func MkdirAll(dirPath string, perm os.FileMode) (bool bool, errInfo error) {
	if perm == 0 {
		perm = DefaultPerm
	}
	err := os.MkdirAll(dirPath, perm)
	if err != nil {
		return false, err
	}
	return true, errInfo
}

// CreatFile 创建文件
func CreatFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	if flag == 0 {
		flag = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	}
	if perm == 0 {
		flag = 0666
	}
	file, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return file, err
	}
	return file, nil
}

// GetFileConfigByPath  只读 - 获取文件信息 // 文件名称 //文件大小
func GetFileConfigByPath(path string) (fileName string, fileSize int) {
	//f,_ := os.Open(path)
	//fmt.Println(f.Name()) //./zydhlog/warning.txt
	//fmt.Println(f.Stat())
	fo, err := os.Stat(path)
	if err != nil {
		//log.Error("open file fail ,file path:" + path,err)
		return "", 0
	}

	return fo.Name(), int(fo.Size())
	//warning.txt //537  //单位，字节
}

// GetFileNameByPath 获取文件的大小 //单位，字节
func GetFileNameByPath(path string) string {
	fo, err := os.Stat(path)
	if err != nil {
		fmt.Println("open file fail ,file path:"+path, err)
		//log.Error("open file fail ,file path:" + path,err)
		return ""
	}
	return fo.Name()
}

// GetFileSizeByPath 获取文件的大小 单位，字节
func GetFileSizeByPath(path string) int {
	fo, err := os.Stat(path)
	if err != nil {
		//log.Error("open file fail ,file path:" + path,err)
		return 0
	}
	return int(fo.Size())
}

// BufioReadFile 缓冲区读取文件
func BufioReadFile(path string, splitChar byte) {
	if !FileOrDirIsExist(path) {
		fmt.Println("path is not exist")
		return
	}

	if splitChar == 0 {
		splitChar = '\n'
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file) //创建读对象 bufio.NewReader()(带缓冲区的方式打开,适合打开较大的文件)
	for {
		line, err := reader.ReadString(splitChar) //注意是字符，单引号  //\n 一行一行的读取 reader.ReadString()(读取文件)
		if err == io.EOF {                        //io.EOF表示读到文件末尾
			if len(line) != 0 {
				fmt.Println(line)
			}
			fmt.Println("文件读完了")
			break
		}
		if err != nil {
			fmt.Println("read file failed, err:", err)
			return
		}
		fmt.Print(line) //这里不用fmt.Println()
	}
}

func WriteFile(txt, path string, flag int, perm os.FileMode) (bool bool, errInfo error) {
	if flag == 0 {
		flag = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	}
	if perm == 0 {
		flag = 0666
	}
	file, err := os.OpenFile(path, flag, perm)

	if err != nil {
		fmt.Println("open file failed, err:", err)
		return false, err
	}
	defer file.Close()
	//str := "hello 沙河"
	//file.Write([]byte(str))       //写入字节切片数据
	//file.WriteString("hello 小王子") //直接写入字符串数据
	_, err = file.WriteString(txt) //直接写入字符串数据
	if err != nil {
		return false, err
	}
	return true, errInfo
}

// FileOrDirIsExist 判断文件是否存在
func FileOrDirIsExist(path string) (bool bool) {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}
	return false

	//通过stat的isDir还可以判断一个路径是文件夹还是文件
	//stat.IsDir()
}

// DelFileByPath 删除文件
func DelFileByPath(path string) (bool bool, errInfo error) {
	if FileOrDirIsExist(path) {
		err := os.Remove(path)
		if err != nil {
			return false, err
		}
		return true, errInfo
	}
	return true, errInfo
}

// FileRename 文件重命名
func FileRename(oldName, newName string) (bool, error) {
	err := os.Rename(oldName, newName)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetFileLastPath 获取文件路径 ./zydh/log/log.txt --> ./zydh/log/
func GetFileLastPath(path string) string {
	//分隔
	splice := strings.Split(path, "/")
	return strings.Join(splice[:len(splice)-1], "/")
}

// FileCutBySize 切割文件
func FileCutBySize(path string) string {
	//判断文件是否超出大小
	if GetFileSizeByPath(path) >= LOG_FILE_MAX_SIZE {
		//2.备份  xx.log -> xx.log.back202205071155  // warning.txt.bak20220507122541.973
		//在原目录备份                             /      warning.txt             warning.txt.bak20220507122541.973
		err := os.Rename(GetFileLastPath(path)+"/"+GetFileNameByPath(path), getLogFilePathByTime(path, 1))
		if err != nil {
			fmt.Println("wavelog/file.go FileCutBySize  Backup failed")
		}
	}
	return path
}

// FileCutByTime 按日期切割
func FileCutByTime(path string) string {
	//返回路径
	return getLogFilePathByTime(path, 2)
}

// getLogFilePathByTime 获取新的日志文件路径
func getLogFilePathByTime(path string, scene int) string {
	//文件路径和名词
	pName := GetFileNameByPath(path) //warning.txt
	var nowStr, logName string
	if scene == 1 {
		nowStr = time.Now().Format(CUT_FILE_BY_SIZE_TIME_FORMAT) //20220507122541.973
		logName = fmt.Sprintf("%s.bak%s", pName, nowStr)         //拼接一个备份
	} else if scene == 2 {
		nowStr = time.Now().Format(CUT_FILE_BY_DATE_TIME_FORMAT) //20220507122541 //按秒分隔
		logName = fmt.Sprintf("[%s]%s", nowStr, pName)           //拼接一个备份
	}

	return GetFileLastPath(path) + "/" + logName
}
