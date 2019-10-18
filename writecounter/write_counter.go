package writecounter

import "io"

type WriteCounter struct {
	total    int64
	current  int64
	written  int64
	reader   io.Reader
	progress Progress
}

type Progress func(current, total int64) error

func New(reader io.Reader, total int64, progress Progress) *WriteCounter {
	wc := new(WriteCounter)
	wc.reader = reader
	wc.total = total
	wc.progress = progress
	return wc
}

func (wc *WriteCounter) Copy(writer io.Writer) error {
	var err error
	wc.written, err = io.Copy(writer, io.TeeReader(wc.reader, wc))
	return err
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.current += int64(n)
	if wc.progress != nil {
		if err := wc.progress(wc.current, wc.total); err != nil {
			return n, err
		}
	}
	return n, nil
}

func (wc *WriteCounter) Written() int64 {
	return wc.written
}
