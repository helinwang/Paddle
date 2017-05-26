package main

/*
#include <string.h>

typedef int reader;
typedef int writer;
*/
import "C"

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unsafe"

	"github.com/PaddlePaddle/Paddle/paddle/go/recordio"
)

const (
	bufferSize = 10
)

var mu sync.Mutex
var handleMap = make(map[C.reader]*reader)
var curHandle C.reader

func addReader(r *reader) C.reader {
	mu.Lock()
	defer mu.Unlock()
	reader := curHandle
	curHandle++
	handleMap[reader] = r
	return reader
}

func getReader(reader C.reader) *reader {
	mu.Lock()
	defer mu.Unlock()
	return handleMap[reader]
}

func removeReader(reader C.reader) *reader {
	mu.Lock()
	defer mu.Unlock()
	r := handleMap[reader]
	delete(handleMap, reader)
	return r
}

func addWriter(r *writer) C.writer {
	mu.Lock()
	defer mu.Unlock()
	reader := curHandle
	curHandle++
	handleMap[reader] = r
	return reader
}

func getReader(reader C.reader) *reader {
	mu.Lock()
	defer mu.Unlock()
	return handleMap[reader]
}

func removeReader(reader C.reader) *reader {
	mu.Lock()
	defer mu.Unlock()
	r := handleMap[reader]
	delete(handleMap, reader)
	return r
}

type reader struct {
	buffer chan []byte
	cancel chan struct{}
}

func read(paths []string, buffer chan<- []byte, cancel chan struct{}) {
	var curFile *os.File
	var curScanner *recordio.Scanner
	var pathIdx int

	var nextFile func() bool
	nextFile = func() bool {
		if pathIdx >= len(paths) {
			return false
		}

		path := paths[pathIdx]
		pathIdx++
		f, err := os.Open(path)
		if err != nil {
			return nextFile()
		}

		idx, err := recordio.LoadIndex(f)
		if err != nil {
			log.Println(err)
			err = f.Close()
			if err != nil {
				log.Println(err)
			}

			return nextFile()
		}

		curFile = f
		curScanner = recordio.NewScanner(f, idx, 0, -1)
		return true
	}

	more := nextFile()
	if !more {
		close(buffer)
		return
	}

	closeFile := func() {
		err := curFile.Close()
		if err != nil {
			log.Println(err)
		}
		curFile = nil
	}

	for {
		for curScanner.Scan() {
			select {
			case buffer <- curScanner.Record():
			case <-cancel:
				close(buffer)
				closeFile()
				return
			}
		}

		if err := curScanner.Error(); err != nil && err != io.EOF {
			log.Println(err)
		}

		closeFile()
		more := nextFile()
		if !more {
			close(buffer)
			return
		}
	}
}

func paddle_new_writer(path *C.char) C.writer {
	return 0
}

//export paddle_new_reader
func paddle_new_reader(path *C.char) C.reader {
	p := C.GoString(path)
	ss := strings.Split(p, ",")
	var paths []string
	for _, s := range ss {
		match, err := filepath.Glob(s)
		if err != nil {
			log.Printf("error applying glob to %s: %v\n", s, err)
			return -1
		}

		paths = append(paths, match...)
	}

	if len(paths) == 0 {
		log.Println("no valid path provided.", p)
		return -1
	}

	buffer := make(chan []byte, bufferSize)
	cancel := make(chan struct{})
	r := &reader{buffer: buffer, cancel: cancel}
	go read(paths, buffer, cancel)
	return addReader(r)
}

//export paddle_reader_next_item
func paddle_reader_next_item(reader C.reader, size *C.int) *C.uchar {
	r := getReader(reader)
	buf, ok := <-r.buffer
	if !ok {
		// channel closed and empty, reached EOF.
		*size = -1
		return (*C.uchar)(unsafe.Pointer(uintptr(0)))
	}

	if len(buf) == 0 {
		// empty item
		*size = 0
		return (*C.uchar)(unsafe.Pointer(uintptr(0)))
	}

	ptr := C.malloc(C.size_t(len(buf)))
	C.memcpy(ptr, unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
	*size = C.int(len(buf))
	return (*C.uchar)(ptr)
}

//export paddle_reader_release
func paddle_reader_release(reader C.reader) {
	r := removeReader(reader)
	close(r.cancel)
}

func main() {} // Required but ignored
