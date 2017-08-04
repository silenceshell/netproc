package main

import (
	"fmt"
	"github.com/silenceshell/netproc/files"
	flag "github.com/spf13/pflag"
	ui "github.com/gizak/termui"
)

func main() {
	var f *string = flag.String("file", "", "specify the file you want to watch")

	flag.Parse()

	filename := *f
	if filename == "" {
		fmt.Println("Usage: netproc --file [snmp, netstat]")
		return
	}
	filename = "/proc/net/" + filename

	err2 := ui.Init()
	if err2 != nil {
		panic(err2)
	}
	defer ui.Close()

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	files.UIStart(filename)

	//switch file {
	//case "netstat":
	//	files.NetstatUIStart()
	//case "snmp":
	//	files.UIStart()
	//}

	ui.Loop()


	fmt.Println(*f)

}
