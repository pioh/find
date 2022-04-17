package main

import (
	"bufio"
	"context"
	"flag"
	"os"
	"sync"

	"github.com/enriquebris/goconcurrentqueue"
)

const parallel = 2

var queue = goconcurrentqueue.NewFIFO()

var wg sync.WaitGroup

func main() {
	flag.Parse()
	target := flag.Arg(0)

	var wout [parallel]*bufio.Writer
	var werr [parallel]*bufio.Writer
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	queue.Enqueue(target)

	for i := 0; i < parallel; i++ {
		wout[i] = bufio.NewWriterSize(os.Stdout, 1024*1024)
		werr[i] = bufio.NewWriter(os.Stderr)
		go worker(ctx, wout[i], werr[i])
	}
	wg.Wait()
	cancel()
	for i := 0; i < parallel; i++ {
		wout[i].Flush()
		werr[i].Flush()
	}
}

func worker(ctx context.Context, wout *bufio.Writer, werr *bufio.Writer) {
	for {
		item, err := queue.DequeueOrWaitForNextElementContext(ctx)
		if err != nil {
			break
		}
		read(item.(string), wout, werr)
	}
}

func read(path string, wout *bufio.Writer, werr *bufio.Writer) {
	defer wg.Done()
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
			wg.Add(1)
			queue.Enqueue(subpath + "/")
		} else {
			wout.WriteString(subpath + "\n")
		}
	}
}
