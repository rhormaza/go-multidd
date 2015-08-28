package main

import (
    "flag"
	"bytes"
	"fmt"
	"log"
//	"math/rand"
	"os/exec"
//	"time"
    "strconv"
)

var SLEEP_RANGE = 500
var MAX = 1
var actionPtr *string

var args_map map[string]*string

//var DD_FMT = "if=/dev/zero of=/tmp/FOO_%d bs=1M count=10 oflag=direct"

type CmdResult struct {
    Stdout string
    Stderr string
}

func runCmd(tag int) string {

	cmd_args := fmt.Sprintf("dd if=/dev/zero of=/tmp/FOO_%d bs=1M count=10 oflag=direct", tag)
	fmt.Println(cmd_args)

	cmd := exec.Command("sh", "-c", cmd_args)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	fmt.Println("Running cmd %s", cmd_args)
	err := cmd.Run()

	if err != nil {
		log.Fatal("DD:", err)
	}
	return fmt.Sprintf("tag[%d]: stdout: %s, stderr:%s", tag, stdout.String(), stderr.String())
}


func runDD(num int) *CmdResult {
    
	cmd_fmt := "dd if=%s of=%s_%d bs=%s count=%s %s"

	//cmd_args := fmt.Sprintf("dd if=/dev/zero of=/tmp/FOO_%d bs=1M count=10 oflag=direct", tag)
	cmd_args := fmt.Sprintf(cmd_fmt, 
        *args_map["infile"], 
        *args_map["outfile"], 
        num, 
        *args_map["block_size"], 
        *args_map["count"], 
        *args_map["flags"])

    fmt.Println("===================================")
	fmt.Println("Worker", num,"cmd ->", cmd_args)
    fmt.Println("===================================")
	cmd := exec.Command("sh", "-c", cmd_args)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		log.Fatal("DD:", err)
	}
	return &CmdResult{stdout.String(), stderr.String()}
}

func doWork() {
    MAX, _ := strconv.Atoi(*args_map["procs"])
	out := make(chan *CmdResult, MAX)
	for i := 0; i < MAX; i++ {
		go func(ii int) {
			out <- runDD(ii)
		}(i)
	}
	for i := 0; i < MAX; i++ {
		select {
		case o := <-out:
            fmt.Println("========= Command output  =========")
			fmt.Printf(o.Stderr)
            fmt.Println("===================================")
		}
	}

}

func parseCmdline() {
    args_map = make(map[string]*string)
    args_map["action"] = flag.String("action", "write", "Write or read to the targets. Defaults to write.")
    args_map["block_size"] = flag.String("blocksize", "1M", "Define the block size to write or read.")
    args_map["procs"] = flag.String("procs", "1", "Number of parallel processes to be run.")
    args_map["count"] = flag.String("count", "1", "Number of times that the block will be writen or read.")
    args_map["infile"] = flag.String("infile", "/dev/zero", "Path to the input file. Same than if= in dd(1).")
    args_map["outfile"] = flag.String("outfile", "/var/tmp/.multidd_data", "Path to the output file.")
    args_map["flags"] = flag.String("flags", "oflag=direct", "Flags to be passed to dd(1)")
    //args_map[""] = flag.String("", "", "")
    
    flag.Parse()
}

func printArgsMap() {
    fmt.Println("====== Printing arguments map ======")
    for k, v := range args_map {
        fmt.Println("key:", k, "-> value: ", *v)
    }
    fmt.Println("====================================")
}

func main() {
    parseCmdline()
    printArgsMap()
    doWork()
}
