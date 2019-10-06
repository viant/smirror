package smirror

import (
	"bufio"
	"io"
	"smirror/config"
)

const bufferSize = 1024 * 1024

func splitByLines(scanner *bufio.Scanner, maxLines int, provider func() io.WriteCloser, rewriter *Rewriter) error {
	counter := 0
	var err error
	if maxLines == 0 {
		maxLines = 1
	}
	var writer io.WriteCloser
	for scanner.Scan() {
		if writer == nil {
			writer = provider()
		}
		if err = scanner.Err(); err != nil {
			return err
		}
		if counter > 0 {
			if _, err = writer.Write([]byte{'\n'}); err != nil {
				return err
			}
		}
		if err = rewriter.Write(writer, scanner.Bytes()); err != nil {
			return err
		}
		counter++
		if counter == maxLines {
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

func splitBySize(scanner *bufio.Scanner, maxSize int, provider func() io.WriteCloser, rewriter *Rewriter) error {
	counter := 0
	var err error
	var writer io.WriteCloser

	sizeSoFar := 0
	for scanner.Scan() {
		if writer == nil {
			writer = provider()
		}
		data := scanner.Bytes()
		if err = scanner.Err(); err != nil {
			return err
		}

		if sizeSoFar+(len(data)+1) >= maxSize && sizeSoFar > 0 {
			if err := writer.Close(); err != nil {
				return err
			}
			sizeSoFar = 0
			counter = 0
			writer = provider()
		}

		sizeSoFar += len(data)
		if counter > 0 {
			sizeSoFar++ //new line character
			if _, err = writer.Write([]byte{'\n'}); err != nil {
				return err
			}
		}
		if err = rewriter.Write(writer, data); err != nil {
			return err
		}
		counter++
	}
	if writer != nil && sizeSoFar > 0 {
		return writer.Close()
	}
	return nil
}

//Split divides reader supplied text by number of specified line
func Split(reader io.Reader, writerProvider func() io.WriteCloser, split *config.Split, rewriter *Rewriter) error {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, bufferSize), 10*bufferSize)
	if split.MaxSize > 0 {
		return splitBySize(scanner, split.MaxSize, writerProvider, rewriter)
	}
	if split.MaxLines == 0 {
		split.MaxLines = 1
	}
	return splitByLines(scanner, split.MaxLines, writerProvider, rewriter)
}
