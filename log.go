package wavelog

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// 日志等级
const (
	ERR_LEVEL_NO int = iota
	ERR_LEVEL_DEBUG
	ERR_LEVEL_TRACE
	ERR_LEVEL_INFO
	ERR_LEVEL_WARNING
	ERR_LEVEL_ERROR
	ERR_LEVEL_FATAL
)

// 日志存放路径
const (
	LOG_FILE_PATH_NORMAL  string = "./log/normal.txt"  //普通日志路径
	LOG_FILE_PATH_WARNING string = "./log/warning.txt" //警告日志路径
	LOG_FILE_PATH_ERROR   string = "./log/error.txt"   //错误日志路径
	LOG_FILE_PATH_FATAL   string = "./log/fatal.txt"   //500服务器内部致命错误日志路径
)

// IS_OPEN_OR_CLOSE_LOG 调试等级，关闭打开 0 关闭，不写入 1 打开
const IS_OPEN_OR_CLOSE_LOG = 1

// LOG_TO_SITE 打印位置 0 终端，1 文件 2 通道->文件
const LOG_TO_SITE = 2

// LOG_CHAN_MAX_CAP 日志通道能放的最大队列数
const LOG_CHAN_MAX_CAP = 1000

// CUT_FILE_TYPE 分隔文件的方式 -1 不分割 0 按照文件大小切割 1 按照时间日期切割
const CUT_FILE_TYPE = 0

// LogChan 日志通道
var LogChan chan *LogToFileByGoChan

// 初始化
func init() {
	CreateLogFile(LOG_FILE_PATH_NORMAL)
	CreateLogFile(LOG_FILE_PATH_WARNING)
	CreateLogFile(LOG_FILE_PATH_ERROR)
	CreateLogFile(LOG_FILE_PATH_FATAL)
	//当开启日志 且 是使用go-chan方法时，初始化通道
	if IS_OPEN_OR_CLOSE_LOG == 1 && LOG_TO_SITE == 2 {
		LogChan = make(chan *LogToFileByGoChan, LOG_CHAN_MAX_CAP) //初始化通道
		runtime.GOMAXPROCS(4)
		//开启通道处理
		for i := 0; i <= 3; i++ {
			go getChanLogToFile(LogChan)
		}
	}
}

// LogContent 日志内容
func LogContent(errLevel, layer int, errInfo string, arg ...any) string {
	timeString := time.Now().Format("2006-01-02 15:04:05.000")
	txtMap := map[int]string{
		ERR_LEVEL_NO:      "[NORMAL]",
		ERR_LEVEL_DEBUG:   "[DEBUG]",
		ERR_LEVEL_TRACE:   "[TRACE]",
		ERR_LEVEL_INFO:    "[INFO]",
		ERR_LEVEL_WARNING: "[WARNING]",
		ERR_LEVEL_ERROR:   "[ERROR]",
		ERR_LEVEL_FATAL:   "[FATAL]",
	}

	return timeString + "  " + getLayerCode(layer) + "\n" + txtMap[errLevel] + fmt.Sprintf("%v", errInfo) + " | " + fmt.Sprintf("%v", arg) + "\n"
}

// getLayerCode 获取代码锚点
func getLayerCode(layer int) string {
	//传递参数，可以拿到当前执行程序执行隔了多少层
	pc, file, line, ok := runtime.Caller(layer)
	if !ok {
		return "use runtime.caller() failed," + "no find layer file name," + "no find layer func name"
	}

	//获得报名和文件名
	pkgNameAndFuncName := strings.Split(runtime.FuncForPC(pc).Name(), ".")

	return pkgNameAndFuncName[0] + "/" + path.Base(file) + " -> " + pkgNameAndFuncName[1] + "()" + " | " + strconv.Itoa(line) + "行"
}

// logWriteToScene 写入日志到相应的通道
func logWriteToScene(errLevel int, logFilePath, logMsg string, arg ...any) {
	//是否开启日志
	if IS_OPEN_OR_CLOSE_LOG > 0 {
		logmsg := LogContent(errLevel, 4, logMsg, arg)
		//开启后，日志写到哪的场景
		switch LOG_TO_SITE {
		case 1: //文件 no goroutine
			writeLogToFile(errLevel, logFilePath, logmsg)
		case 2: //发送到通道 ToChan
			select {
			case LogChan <- NewLogToFileByGoChan(errLevel, logFilePath, logmsg): //当通道没满时，走这
			default: //当通道满了，走这，丢弃日志，保证不出现阻塞
			}
		default: //终端
			outputToTerminal(logmsg)
		}
	}
}

//实现方法 - 不使用接口，接口每次都要先实例化一个对应的结构体，错误

func Trace(logInfo string, arg ...any) {
	logWriteToScene(ERR_LEVEL_TRACE, LOG_FILE_PATH_NORMAL, logInfo, arg)
}

func Debug(logInfo string, arg ...any) {
	logWriteToScene(ERR_LEVEL_DEBUG, LOG_FILE_PATH_NORMAL, logInfo, arg)
}

func Warning(logInfo string, arg ...any) {
	logWriteToScene(ERR_LEVEL_WARNING, LOG_FILE_PATH_WARNING, logInfo, arg)
}

func Info(logInfo string, arg ...interface{}) {
	logWriteToScene(ERR_LEVEL_INFO, LOG_FILE_PATH_NORMAL, logInfo, arg)
}

func Error(logInfo string, arg ...any) {
	logWriteToScene(ERR_LEVEL_ERROR, LOG_FILE_PATH_ERROR, logInfo, arg)
}

func Fatal(logInfo string, arg ...any) {
	logWriteToScene(ERR_LEVEL_FATAL, LOG_FILE_PATH_FATAL, logInfo, arg)
}
