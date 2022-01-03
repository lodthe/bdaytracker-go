package filereader

import (
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

type FileReader struct {
	filepath string
	name     string
}

func NewFileReader(assetsPath, fileName string) *FileReader {
	return &FileReader{
		filepath: path.Join(assetsPath, fileName) + fileName,
		name:     fileName,
	}
}

func (s *FileReader) Name() string {
	return s.name
}

func (s *FileReader) Reader() (io.Reader, error) {
	return os.Open(s.filepath)
}

func (s *FileReader) Size() int64 {
	file, err := os.Stat(s.filepath)
	if err != nil {
		logrus.WithField("filepath", s.filepath).WithError(err).Error("failed to get the file size")
		return 0
	}

	return file.Size()
}
