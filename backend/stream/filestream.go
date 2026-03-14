package stream

import "os"

type FileStream struct {
	file *os.File
}

func NewFileStream(p string) (*FileStream, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	return &FileStream{
		file: f,
	}, nil
}

func (s *FileStream) Read(p []byte) (int, error) {
	return s.file.Read(p)
}

func (s *FileStream) Write(p []byte) (int, error) {
	panic("File Stream received a command")
}

func (s *FileStream) Close() error {
	return s.file.Close()
}
