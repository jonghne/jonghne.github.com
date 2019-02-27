package main

import (
	"os/exec"
	"bytes"
	"fmt"
	"bufio"
	"io"
	"strings"
	"time"
	"io/ioutil"
	"path/filepath"
	"chainstack-core/common/util"
	"runtime"
)

//阻塞式的执行外部shell命令的函数,等待执行完毕并返回标准输出
func exec_shell(s string) (string, error){
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", s)

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()
	checkErr(err)


	return out.String(), err
}

func execCommand(commandName string, params []string) bool {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command(commandName, params...)

	//显示运行的命令
	fmt.Println(cmd.Args)
	//StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return false
	}

	cmd.Start()
	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}

	//阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
	cmd.Wait()
	return true
}

//不需要执行命令的结果与成功与否，执行命令马上就返回
func exec_shell_no_result(command string) {
	//处理启动参数，通过空格分离 如：setsid /home/luojing/gotest/src/test_main/iwatch/test/while_little &
	//command_name_and_args := strings.FieldsFunc(command, splite_command)
	command_name_and_args := strings.Split(command, ",")

	cmd := exec.Command(command_name_and_args[0], command_name_and_args[1:]...)
	//开始执行c包含的命令，但并不会等待该命令完成即返回
	err := cmd.Start()
	if err != nil {
		fmt.Printf("%v: exec command:%v error:%v\n", get_time(), command, err)
	}
	fmt.Printf("Waiting for command:%v to finish...\n", command)
	//阻塞等待fork出的子进程执行的结果，和cmd.Start()配合使用[不等待回收资源，会导致fork出执行shell命令的子进程变为僵尸进程]
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("%v: Command finished with error: %v\n", get_time(), err)
	}
	return
}

func get_time() string {
	return time.Now().Format("2019年03月05日")
}

//错误处理函数
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func main() {
	//execCommand("ls", []string{"-al"})
	//exec_shell_no_result("ls,-la")

	_, e := ioutil.ReadFile(filepath.Join("/home/qydev/idtest"))
	fmt.Println(e)

	id := util.StringifyJsonToBytes([]interface{}{"0x2432424", "0x325252a"})

	ioutil.WriteFile(filepath.Join("/home/qydev/idtest"), id, 0644)

	conf, e := ioutil.ReadFile(filepath.Join("/home/qydev/idtest"))

	if e == nil {
		var str []string
		err := util.ParseJsonFromBytes(conf, &str)
		if err == nil {
			fmt.Println(str)
		} else {
			fmt.Println(err)
		}
	}
	var stat runtime.MemStats
	runtime.ReadMemStats(&stat)
	fmt.Println(runtime.NumGoroutine(), stat)
	var s string = "sljflsd"
	fmt.Println(s)
}