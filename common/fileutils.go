package common

import (
	"os"
	"bufio"
)

type line func(string)

func ReadFile(file string, fn line) {

	inFile, err := os.Open(file)
	if err != nil {

	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fn(scanner.Text())
	}
}
