package smirror

import (
	"bufio"
	"io"
)

//Split divides reader supplied text by number of specified line
func Split(reader io.Reader, writerProvider func() io.WriteCloser, elementCount int) error {
	scanner := bufio.NewScanner(reader)
	var writer io.WriteCloser
	counter := 0
	var err error
	if elementCount == 0 {
		elementCount = 1
	}
	for scanner.Scan() {
		if writer == nil {
			writer = writerProvider()
		}
		data := scanner.Bytes()
		if err = scanner.Err(); err != nil {
			return err
		}

		if counter > 0 {
			if _, err = writer.Write([]byte{'\n'}); err != nil {
				return err
			}
		}
		if _, err = writer.Write(data); err != nil {
			return err
		}
		counter++
		if counter == elementCount {
			if err := writer.Close(); err != nil {
				return err
			}
			counter = 0
			writer = nil
		}
	}
	if writer != nil {
		return writer.Close()
	}
	return nil
}
