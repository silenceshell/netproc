package files

import (
	"fmt"
	"time"
	"os"
	"bufio"
	"log"
	"strings"
	"os/exec"
	"strconv"
)

type netInfo struct {
	count int
	key   string
}

func getTerminalSize() (high, length int) {
	cmd := exec.Command("stty", "size")
  	cmd.Stdin = os.Stdin
  	out, err := cmd.Output()
	if err != nil {
		high = 40
		length = 80
	} else {
		fmt.Sscanf(string(out), "%d %d", &high, &length)
	}
	return high, length
}

func getPar(netInfo []netInfo) (string, int){
	var _parser string
	_, width := getTerminalSize()
	width -=6
	fmtStringLen := 0
	total := 0
	lines := 1

	for _, v := range netInfo {
		fmtString := fmt.Sprintf("%25s: %-10d", v.key, v.count)
		fmtStringLen = len(fmtString)
		total += fmtStringLen
		if total > width {
			_parser += "\r\n"
			lines += 1
			total = fmtStringLen
		}
		_parser += fmtString
	}
	return _parser, lines
}

func getTime() string {
	var _parser string

	t := time.Now()
	header := t.Format("2006-01-02 15:04:05")

	_, width := getTerminalSize()
	width -=6

	headerFmt := fmt.Sprintf("[%%%ds\r\n](fg-yellow)", (width + len(header))/2)
	_parser = fmt.Sprintf(headerFmt, header)

	return _parser
}

func isKey(text string) bool {
	ret := false
	_, err := strconv.Atoi(text)
	if err != nil {
		ret = true
	}
	return ret
}

func procText(texts string, infoArray []netInfo) []netInfo {
	textArray := strings.Split(strings.TrimSpace(texts), " ")
	if infoArray == nil {
		infoArray = make([]netInfo, len(textArray))
	}
	if isKey(textArray[0]) {
		for i := 0; i < len(textArray); i++ {
			infoArray[i].key = textArray[i]
		}
	} else {
		for i := 0; i < len(textArray); i++ {
			infoArray[i].count, _ = strconv.Atoi(textArray[i])
		}
	}

	return infoArray
}

func getInfoMapMap(file_path string) (map[string][]netInfo){
	file, err := os.Open(file_path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	netInfoMap := make(map[string][]netInfo, 64)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subString := scanner.Text()
		segTexts := strings.Split(subString, ":")

		seg := segTexts[0]
		texts := segTexts[1]

		netInfoMap[seg] = procText(texts, netInfoMap[seg])
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return netInfoMap
}
