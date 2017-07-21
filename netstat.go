package main

import (
	"fmt"
	//"io/ioutil"
	//"strings"
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"os/exec"
)

type netExt struct {
	count int
	key   string
}

func isKey(text string) bool {
	ret := false
	_, err := strconv.Atoi(text)
	if err != nil {
		ret = true
	}
	return ret
}

func procText(texts string, netExtArray []netExt) []netExt {
	textArray := strings.Split(strings.TrimSpace(texts), " ")
	if netExtArray == nil {
		netExtArray = make([]netExt, len(textArray))
	}
	if isKey(textArray[0]) {
		for i := 0; i < len(textArray); i++ {
			netExtArray[i].key = textArray[i]
		}
	} else {
		for i := 0; i < len(textArray); i++ {
			netExtArray[i].count, _ = strconv.Atoi(textArray[i])
		}
	}

	return netExtArray
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

func printNetExt(header string, netExtArray []netExt) {

	_, length := getTerminalSize()
	total := 0
	fmtStringLen := 0

	_header := "========================="
	header = _header + header + _header

	headerFmt := fmt.Sprintf("%%%ds\r\n", (length + len(header))/2)

	fmt.Printf(headerFmt, header)

	for _, v := range netExtArray {
		fmtString := fmt.Sprintf("%25s: %-10d", v.key, v.count)
		fmtStringLen = len(fmtString)
		total += fmtStringLen
		if total > length {
			fmt.Println()
			total = fmtStringLen
		}
		fmt.Printf("%s", fmtString)
	}
	fmt.Println()
}

func main() {
	statPath := "/proc/net/netstat"

	file, err := os.Open(statPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var tcpExtArray []netExt
	var ipExtArray []netExt

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		subString := scanner.Text()
		segTexts := strings.Split(subString, ":")

		seg := segTexts[0]
		texts := segTexts[1]
		switch seg {
		case "TcpExt":
			tcpExtArray = procText(texts, tcpExtArray)
		case "IpExt":
			ipExtArray = procText(texts, ipExtArray)
		}
	}


	printNetExt("TcpExt", tcpExtArray)

	printNetExt("IpExt", ipExtArray)

	fmt.Println()

	//for _, v := range ipExtArray {
	//	fmt.Printf("%s: %d\r\n", v.key, v.count)
	//}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
