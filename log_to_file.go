package wavelog

import (
	"fmt"
	"os"
	"time"
)

// LogToFileByGoChan 写入通道的日志格式
type LogToFileByGoChan struct {
	leveL       int    //错误级别
	logFilePath string //日志存放位置
	logMsg      string //错误的信息
}

// NewLogToFileByGoChan 初始化
func NewLogToFileByGoChan(level int, logFilePath, logMsg string) *LogToFileByGoChan {
	return &LogToFileByGoChan{
		leveL:       level,
		logFilePath: logFilePath,
		logMsg:      logMsg,
	}
}

/*
func : 启动时,开启goroutine处理通道数据
param : ch 只读通道  *LogToFileByGoChan类型
*/
func getChanLogToFile(ch <-chan *LogToFileByGoChan) {
	for {
		select {
		//当通道有值时
		case logData := <-ch:
			writeLogToFile(logData.leveL, logData.logFilePath, logData.logMsg)
		default:
			time.Sleep(time.Millisecond * 500) //当通道无值时，交出cpu控制权
		}
	}
}

//writeLogToFile 写日志到文件
func writeLogToFile(errLevel int, logFilePath, logMsg string, arg ...any) {
	if errLevel >= 0 {

		// 是否切割文件 获取日志文件地址
		filePath := fileIsOpenCut(logFilePath)

		//打开文件 //判断文件是否存在
		f, err := CreatFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Printf("creat file false,err %v \n", err)
			return
		}
		defer f.Close()

		//直接写入字符串数据
		_, err = f.WriteString(logMsg)
		if err != nil {
			fmt.Println("log to file fail,err:", err)
		}
	}
}

/*
是否开启文件切割
*/
func fileIsOpenCut(logFilePath string) (filePath string) {
	if CUT_FILE_TYPE >= 0 {
		//切割文件
		switch CUT_FILE_TYPE {
		case 0: //按照文件大小切割
			filePath = FileCutBySize(logFilePath)
		case 1: //按照时间切割
			filePath = FileCutByTime(logFilePath)
		default:
			filePath = FileCutByTime(logFilePath)
		}
	} else {
		filePath = logFilePath
	}
	return filePath
}
