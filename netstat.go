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
	ui "github.com/gizak/termui"
	"time"
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

func getP(netExtArray []netExt) *ui.Par{

	hight, width := getTerminalSize()
	width -=6
	total := 0
	fmtStringLen := 0
	var _parser string


	t := time.Now()
	header := t.Format("2006-01-02 15:04:05")

	headerFmt := fmt.Sprintf("[%%%ds\r\n](fg-yellow)", (width + len(header))/2)
	_parser += fmt.Sprintf(headerFmt, header)

	for _, v := range netExtArray {
		fmtString := fmt.Sprintf("%25s: %-10d", v.key, v.count)
		fmtStringLen = len(fmtString)
		total += fmtStringLen
		if total > width {
			//fmt.Println()
			_parser += "\r\n"
			total = fmtStringLen
		}
		_parser += fmtString
		//fmt.Printf("%s", fmtString)
	}

	p := ui.NewPar(_parser)
	p.WrapLength = width-2 // this should be at least p.Width - 2
	p.Height = hight
	p.Width = width
	p.Y = 0
	p.X = 2
	p.TextFgColor = ui.ColorGreen
	p.BorderLabel = "TcpExt"
	p.BorderFg = ui.ColorCyan
	//p.Border = false

	return p
}

func uiStart() {
	//var netExtArray []netExt

	err2 := ui.Init()
	if err2 != nil {
		panic(err2)
	}
	defer ui.Close()

	tcpExtArray, _ := getArray()
	ui.Render(getP(tcpExtArray))

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/timer/1s", func(e ui.Event) {
		t := e.Data.(ui.EvtTimer)
		i := t.Count
		if i > 103 {
			ui.StopLoop()
			return
		}
		tcpExtArray, _ := getArray()
		ui.Clear()
		ui.Render(getP(tcpExtArray))
	})
	ui.Loop()


}

func getArray() ([]netExt, []netExt){

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

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return tcpExtArray, ipExtArray
}

func main() {
	//statPath := "/proc/net/netstat"
	//
	//file, err := os.Open(statPath)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer file.Close()
	//
	//var tcpExtArray []netExt
	//var ipExtArray []netExt
	//
	//scanner := bufio.NewScanner(file)
	//for scanner.Scan() {
	//	//fmt.Println(scanner.Text())
	//	subString := scanner.Text()
	//	segTexts := strings.Split(subString, ":")
	//
	//	seg := segTexts[0]
	//	texts := segTexts[1]
	//	switch seg {
	//	case "TcpExt":
	//		tcpExtArray = procText(texts, tcpExtArray)
	//	case "IpExt":
	//		ipExtArray = procText(texts, ipExtArray)
	//	}
	//}
	//
	//
	//printNetExt("TcpExt", tcpExtArray)
	//
	//printNetExt("IpExt", ipExtArray)
	//
	//fmt.Println()

	//for _, v := range ipExtArray {
	//	fmt.Printf("%s: %d\r\n", v.key, v.count)
	//}

	//if err := scanner.Err(); err != nil {
	//	log.Fatal(err)
	//}

	uiStart()
}
