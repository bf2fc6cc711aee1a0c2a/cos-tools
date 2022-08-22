package cmdutil

import (
	"bufio"
	"github.com/hashicorp/go-multierror"
	"io"
	"os"
)

type OutputWriter struct {
	delegate  io.Writer
	writer    *bufio.Writer
	mustClose bool
}

func NewOutputWriter(delegate io.Writer) (*OutputWriter, error) {
	answer := OutputWriter{
		delegate:  delegate,
		writer:    bufio.NewWriter(delegate),
		mustClose: false,
	}

	return &answer, nil
}

func NewOutputFileWriter(path string) (*OutputWriter, error) {
	delegate, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	answer := OutputWriter{
		delegate:  delegate,
		writer:    bufio.NewWriter(delegate),
		mustClose: true,
	}

	return &answer, nil
}

func (in *OutputWriter) Write(p []byte) (n int, err error) {
	return in.writer.Write(p)
}

func (in *OutputWriter) Close() error {
	var err error

	if ferr := in.writer.Flush(); err != nil {
		err = multierror.Append(err, ferr)
	}

	if in.mustClose {
		if c, ok := in.delegate.(io.Closer); ok {
			if cerr := c.Close(); cerr != nil {
				err = multierror.Append(err, cerr)
			}
		}
	}

	return err
}
