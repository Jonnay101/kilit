package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	port := flag.String("p", "8080", "the port you want to kill")
	flag.Parse()
	task := "lsof"
	cmd := exec.Command(task, "-i", fmt.Sprintf(":%s", *port))
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "COMMAND") {
			continue
		}
		if line == "" {
			break
		}
		words := strings.Fields(line)
		pid := words[1]
		kill := exec.Command("kill", "-9", pid)
		kill.Run()
	}
}
