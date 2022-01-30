package main

import (
	"GoWolf/icmpSurvive"
	"GoWolf/portScan"
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	icmp      int
	addr      string
	port      string
	threading int
	fileName  string
	jobNum    int
	finishNum int
	addrs     []string
)

var ini = `

  _____   __          __   _  __ 
 / ____|  \ \        / /  | |/ _|
| |  __  __\ \  /\  / /__ | | |_ 
| | |_ |/ _ \ \/  \/ / _ \| |  _|
| |__| | (_) \  /\  / (_) | | |  
 \_____|\___/ \/  \/ \___/|_|_|  
                                 

`

func init() {
	flag.StringVar(&addr, "a", "", "目标地址")
	flag.StringVar(&fileName, "f", "", "目标地址文件")
	flag.StringVar(&port, "p", "1-100", "目标端口地址, 单个仅输入一个端口, 多个用-分开,默认1-100")
	flag.IntVar(&threading, "t", 5, "设置go程数，默认10个")
	flag.IntVar(&jobNum, "J", 200, "设置工作区缓冲数量, 默认200")
	flag.IntVar(&finishNum, "O", 200, "设置完成区缓冲数量, 默认200")
	flag.IntVar(&icmp, "i", 1, "icmp存活扫描，默认关闭1, 开启0")
	fmt.Printf("%s", ini)
}

func oneAddr() {
	addrs = append(addrs, addr)
	p := portScan.NewScan(threading, port, addrs, jobNum, finishNum)
	p.Start()
	p.Close()
}

func fileAddr() {
	p := portScan.NewScan(threading, port, addrs, jobNum, finishNum)
	p.Start()
	p.Close()
}

func icmpSave() {
	i := icmpSurvive.NewScan(threading, addrs, jobNum, finishNum)
	i.Start()
	i.Close()
}

func fileHandle() {
	file, fileErr := os.Open(fileName)
	if fileErr != nil {
		panic(err("文件打开失败", fileErr))
	}
	content := bufio.NewScanner(file)
	for {
		if !content.Scan() {
			break
		}
		host := strings.TrimSpace(content.Text())
		addrs = append(addrs, host)
	}
}

func err(text string, arg interface{}) error {
	return errors.New(fmt.Sprintf("%s err: %s", text, arg))
}

func main() {
	flag.Parse()
	if a := net.ParseIP(addr); addr != "" && a != nil {
		oneAddr()
	}

	if addr == "" && fileName != "" {
		fileHandle()
	}

	if addr == "" && fileName != "" && icmp == 1 {
		fileAddr()
	}

	if addr == "" && fileName != "" && icmp == 0 {
		icmpSave()
	}
}
