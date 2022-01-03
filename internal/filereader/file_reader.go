package filereader

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type FileReader struct {
	file string
	name string
}

func NewFileReader(fileName string) *FileReader {
	return &FileReader{
		file: "./assets/" + fileName,
		name: fileName,
	}
}

func (s *FileReader) Name() string {
	return s.name
}

func (s *FileReader) Reader() (io.Reader, error) {
	return os.Open(s.file)
}

func (s *FileReader) Size() int64 {
	file, err := os.Stat(s.file)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"file": s.file,
		}).WithError(err).Error("failed to get the file size")
		return 0
	}

	return file.Size()
}
