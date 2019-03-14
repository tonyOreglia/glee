package pipe

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"syscall"
)

type Pipe struct {
	Reader *bufio.Reader
	Writer *os.File
	Chan   chan string
}

func NewPipe(pipeFile string) *Pipe {
	p := new(Pipe)
	os.Remove(pipeFile)
	err := syscall.Mkfifo(pipeFile, 0666)
	if err != nil {
		log.Fatal("Make named pipe file error:", err)
	}
	p.Chan = make(chan string)

	p.Writer, err = os.OpenFile(pipeFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	fmt.Println("open a named pipe file.")
	file, err := os.OpenFile(pipeFile, os.O_CREATE, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
	}
	p.Reader = bufio.NewReader(file)
	return p
}

func (p *Pipe) StartReader() {
	for {
		line, err := p.Reader.ReadBytes('\n')
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		msg := string(line[0 : len(line)-1]) // drop the newline
		fmt.Printf("Engine recieved command: %s\n", msg)
		p.Chan <- msg
	}
}

func (p *Pipe) Write(msg string) {
	p.Writer.WriteString(fmt.Sprintf("%s\n", msg))
}
