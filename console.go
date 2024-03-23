package wavelog

import "fmt"

// outputToTerminal 写入终端
func outputToTerminal(logMsg string) {
	fmt.Println(logMsg)
	//往终端写入
	//_, err := fmt.Fprint(os.Stdout, logMsg)
	//if err != nil {
	//	panic(err)
	//}
}
