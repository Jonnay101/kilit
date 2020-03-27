package main

import (
	"bufio"
	"bytes"
	"flag"
	"log"
	"os/exec"
	"strings"
)

func main() {
	portToKill := flag.String("p", ":8080", "the port you want to kill")
	flag.Parse()
	// execute the command then fetch the combined standard out and error
	// the standard out should be a table representation of the list of open files
	lsofCommand := exec.Command("lsof", "-i", *portToKill)
	lsofTable, err := lsofCommand.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	lsofTableReader := bytes.NewReader(lsofTable)
	scanner := bufio.NewScanner(lsofTableReader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() { // for each row of the scanned table
		lsofTableRow := scanner.Text()
		if strings.Contains(lsofTableRow, "COMMAND") {
			continue // ignore the title row of the table
		}
		if lsofTableRow == "" {
			break // if row is empty, cease iterating
		}
		lsofTableCells := strings.Fields(lsofTableRow)
		pid := lsofTableCells[1]
		killCommand := exec.Command("kill", "-9", pid)
		killCommand.Run()
	}
}
