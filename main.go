package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"
	"github.com/jacobsa/go-serial/serial"
)

func run(c *cli.Context) {
	options := serial.OpenOptions{
		PortName:        c.String("port-name"),
		BaudRate:        57600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	port, err := serial.Open(options)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	go read(port)
	write(port)
}

func read(r io.Reader) {
	for {
		if _, err := io.Copy(os.Stdout, r); err != nil {
			log.Fatal(err)
		}
	}
}

func write(w io.Writer) {
	for {
		fmt.Printf("\nEnter command: ")
		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		cmd := fmt.Sprintf("%s\r\n", strings.TrimSpace(line))
		if _, err := w.Write([]byte(cmd)); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "lorabeeterminal"
	app.Usage = "terminal to read and write from / to the LoRaBee"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port-name",
			Value: "/dev/tty.usbserial-AH02ZDDI",
		},
	}
	app.Action = run
	app.Run(os.Args)
}
