package domain

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

type DataStore struct {
	data []byte
}

func (d *DataStore) OnDataCreation(data []byte) {
	d.data = data
}

func (d *DataStore) Data() []byte {
	return d.data
}

type pidListener interface {
	OnPIDFound(string)
}

func GetAllPIDsFromData(data []byte, lst pidListener) {
	readAllPIDs(splitDataByLine(data), lst)
}

func splitDataByLine(data []byte) []string {
	lines := make([]string, 0)
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Split(bufio.ScanLines)

	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	return lines
}

func readAllPIDs(lines []string, lst pidListener) {
	for _, line := range lines {
		readPIDFromString(line, lst)
	}
}

func readPIDFromString(str string, lst pidListener) {
	defer func() {
		recover() // recover the panic from getWordByIndex but return nothing
	}()
	if strings.Contains(str, "COMMAND") {
		return
	}
	pid := getWordByIndex(str, 1)
	if pid == "" || lst == nil {
		return
	}
	lst.OnPIDFound(pid)
}

// getWordByIndex takes a string and an index int and returns the word at that index
// getWordByIndex will PANIC if the idx is greater than the word length of str
func getWordByIndex(str string, idx int) string {
	if str == "" || idx < 0 {
		return ""
	}

	words := strings.Fields(str)
	if (idx + 1) > len(words) {
		panic(fmt.Errorf("getWordByIndex: the provided index of %d is greater than the word count (%d) in the provided string", idx, len(words)))
	}

	return words[idx]
}
