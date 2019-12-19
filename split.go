package smirror

import (
	"bufio"
	"github.com/pkg/errors"
	"io"
	"smirror/config"
)

const bufferSize = 1024 * 1024

func splitByLines(scanner *bufio.Scanner, maxLines int, provider func(partition interface{}) io.WriteCloser) error {
	counter := 0
	var err error
	if maxLines == 0 {
		maxLines = 1
	}
	var writer io.WriteCloser
	for scanner.Scan() {
		if writer == nil {
			writer = provider("")
		}
		if err = scanner.Err(); err != nil {
			return err
		}
		if counter > 0 {
			if _, err = writer.Write([]byte{'\n'}); err != nil {
				return err
			}
		}
		if _, err = writer.Write(scanner.Bytes()); err != nil {
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

func splitBySize(scanner *bufio.Scanner, maxSize int, provider func(partition interface{}) io.WriteCloser) error {
	counter := 0
	var err error
	var writer io.WriteCloser

	sizeSoFar := 0
	for scanner.Scan() {
		if writer == nil {
			writer = provider("")
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
			writer = provider("")
		}

		sizeSoFar += len(data)
		if counter > 0 {
			sizeSoFar++ //new line character
			if _, err = writer.Write([]byte{'\n'}); err != nil {
				return err
			}
		}
		if _, err = writer.Write(data); err != nil {
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
func Split(reader io.Reader, writerProvider func(partition interface{}) io.WriteCloser, rule *config.Rule) error {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, bufferSize), 10*bufferSize)
	split := rule.Split

	//TODO split avro
	if split.Partition != nil {
		return splitWithPartition(scanner, split, writerProvider)
	}
	if split.MaxSize > 0 {
		return splitBySize(scanner, split.MaxSize, writerProvider)
	}
	if split.MaxLines == 0 {
		split.MaxLines = 1
	}
	return splitByLines(scanner, split.MaxLines, writerProvider)
}

type partitionedBuffer struct {
	key   interface{}
	lines int
	size  int
	io.WriteCloser
}

func (p *partitionedBuffer) flush(provider func(partition interface{}) io.WriteCloser) error {
	if err := p.Close(); err != nil {
		return err
	}
	p.size = 0
	p.lines = 0
	if provider == nil {
		return nil
	}
	p.WriteCloser = provider(p.key)
	return nil
}

func splitWithPartition(scanner *bufio.Scanner, split *config.Split, provider func(partition interface{}) io.WriteCloser) (err error) {

	var partitions = map[interface{}]*partitionedBuffer{}
	for scanner.Scan() {
		data := scanner.Bytes()
		if err = scanner.Err(); err != nil {
			return err
		}
		key, err := split.Partition.Key(data)
		if err != nil {
			return errors.Wrapf(err, "failed to get key from %s", data)
		}
		partition, ok := partitions[key]
		if !ok {
			partition = &partitionedBuffer{WriteCloser: provider(key), key: key}
			partitions[key] = partition
		}

		if (split.MaxLines > 0 && partition.lines+1 > split.MaxLines) || (split.MaxSize > 0 && partition.size+len(data) > split.MaxSize) {
			if err = partition.flush(provider); err != nil {
				return nil
			}
		}
		if partition.size > 0 {
			partition.size++
			if _, err = partition.Write([]byte{'\n'}); err != nil {
				return err
			}
		}
		if _, err = partition.Write(data); err != nil {
			return err
		}
		partition.lines++
		partition.size += len(data)
	}
	for k := range partitions {
		partition := partitions[k]
		if err = partition.flush(nil); err != nil {
			return err
		}
	}
	return nil
}
