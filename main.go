package main

import (
	"bufio"
	"flag"
	"os"
)

var queue = &Stack{}

func main() {
	flag.Parse()
	target := flag.Arg(0)

	var wout *bufio.Writer
	var werr *bufio.Writer

	queue.Push(target)

	wout = bufio.NewWriterSize(os.Stdout, 1024*1024)
	werr = bufio.NewWriter(os.Stderr)
	defer wout.Flush()
	defer werr.Flush()

	for queue.Size() > 0 {
		read(queue.Pop(), wout, werr)
	}
}

func read(path string, wout *bufio.Writer, werr *bufio.Writer) {
	file, err := os.Open(path)
	if err != nil {
		werr.WriteString(path + " failed: " + err.Error() + "\n")
		return
	}
	defer file.Close()
	items, err := file.ReadDir(-1)
	if err != nil {
		werr.WriteString(path + " failed: " + err.Error() + "\n")
		return
	}
	file.Close()
	for _, item := range items {
		subpath := path + item.Name()
		if item.IsDir() {
			queue.Push(subpath + "/")
		} else {
			wout.WriteString(subpath + "\n")
		}
	}
}

type Stack struct {
	buf []string
}

func (s *Stack) Push(item string) {
	s.buf = append(s.buf, item)
}
func (s *Stack) Pop() string {
	n := len(s.buf) - 1
	item := s.buf[n]
	s.buf = s.buf[:n]
	return item
}

func (s *Stack) Size() int {
	return len(s.buf)
}
