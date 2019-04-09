package ROTreader

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
}

func NewRot128Reader(filepath string) *Rot128Reader {
	fd, err := os.Open(filepath)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"filepath": filepath,
		}).Fatal("Couldn't open file")
	}
	rotReader := Rot128Reader{
		reader: bufio.NewReader(fd),
	}
	rotReader.initScanner()
	return &rotReader
}

func (r *Rot128Reader) initScanner() {
	if r.scanner != nil {
		log.Warn("Scanner already initialized.")
		return
	}
	if r.reader == nil {
		log.Fatal("Something deeply wrong happened, the io.Reader is nil")
		return
	}
	s := bufio.NewScanner(r.reader)
	s.Split(scanLinesRot128)
	r.scanner = s
}

func (r *Rot128Reader) Scan() (string, bool) {
	data := []byte{}
	if ok := r.scanner.Scan(); ok {
		data = cipher.Rot128(r.scanner.Bytes())
	}
	return data, ok
}
