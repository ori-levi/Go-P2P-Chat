package main

import (
	"flag"
	"github.com/jroimartin/gocui"
	app "levi.ori/p2p-chat/src/ui"
	widgets "levi.ori/p2p-chat/src/ui/widgets"
	"log"
)

func main() {
	var name string
	flag.StringVar(&name, "name", "", "Your client name")

	var port int
	flag.IntVar(&port, "port", 9090, "local server port for listening")

	flag.Parse()

	if len(name) == 0 {
		log.Fatalln("Name is missing please run with -name <name>")
	}

	mainApp, err := app.NewApp()
	if err != nil {
		log.Panicln(err)
	}
	defer mainApp.Close()

	inputView := widgets.NewInputWidget(name)
	inputView.AddHandler(widgets.KeyHandler{
		Key:     gocui.KeyEnter,
		Handler: onInputChange(inputView.OnValueChange),
	})

	managers := []gocui.Manager{
		widgets.NewHelpWidget(),
		widgets.NewUsersWidget(),
		widgets.NewChatWidget(),
		widgets.NewLogWidget(),
		inputView,
	}

	if err := mainApp.Run(managers...); err != nil {
		log.Panicln(err)
	}

	//
	//inputChannel := make(chan string)
	//logChannel := make(chan string)
	//
	//fmt.Println(inputChannel)
	//fmt.Println(logChannel)
	//
	//serverApp := server.NewServer(name, port, localInterfaceOnly, logChannel)
	//go serverApp.RunServer()
	//
	//clientApp := client.NewClient(name, port, logChannel)
	//go clientApp.Run(serverApp.InChannel, inputChannel)
	//
	//old_code.uiMain(name, logChannel, serverApp.OutChannel, inputChannel)
}
