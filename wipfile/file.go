package wipfile

import (
	"bufio"
	"os"
)

type file struct {
	filename string
}

func (f *file) AppendLine(line []byte) error {
	fh, err := f.openAppendable()
	if err != nil {
		return err
	}
	defer fh.Close()
	w := bufio.NewWriter(fh)
	_, err = w.Write(line)
	err = w.WriteByte('\n')
	if err != nil {
		return err
	}
	return w.Flush()
}

func (f *file) Lines(onLine func([]byte) error) error {
	fh, err := f.openReadable()
	if err != nil {
		return err
	}
	defer fh.Close()
	lines := bufio.NewScanner(fh)
	lines.Split(bufio.ScanLines)
	for lines.Scan() {
		err = onLine(lines.Bytes())
		if err != nil {
			return err
		}
	}
	return lines.Err()
}

func (f *file) openAppendable() (*os.File, error) {
	return os.OpenFile(f.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
}

func (f *file) openReadable() (*os.File, error) {
	return os.OpenFile(f.filename, os.O_CREATE|os.O_RDONLY, 0600)
}
