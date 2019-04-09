package rot128

import (
	"bufio"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// Rot128Reader implements io.Reader that transforms
type Rot128Reader struct {
	reader  io.Reader
	scanner *bufio.Scanner
	fd      *os.File
}

func (r *Rot128Reader) Scan() (string, bool) {
	data := []byte{}
	ok := r.scanner.Scan()
	if ok {
		data = r.scanner.Bytes()
		Rot128(data)
	}
	return string(data), ok
}

func (r *Rot128Reader) Close() {
	r.fd.Close()
}

func NewRot128Reader(filepath string) *Rot128Reader {
	fd, err := os.Open(filepath)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"filepath": filepath,
		}).Fatal("Couldn't open file")
	}
	reader := bufio.NewReader(fd)
	if reader == nil {
		fd.Close()
		log.Fatal("Something deeply wrong happened, the io.Reader is nil")
	}
	s := bufio.NewScanner(reader)
	s.Split(scanLinesRot128)
	return &Rot128Reader{
		reader:  reader,
		fd:      fd,
		scanner: s,
	}
}
