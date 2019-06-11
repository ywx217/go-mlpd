package main

import "os"

// ProgressFile file with read progress
type ProgressFile struct {
	file      *os.File
	readSize  int64
	totalSize int64
}

// OpenProgressFile opens an existing file
func OpenProgressFile(path string) (*ProgressFile, error) {
	finfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &ProgressFile{
		file:      f,
		readSize:  0,
		totalSize: finfo.Size(),
	}, nil
}

// Read reads up to len(b) bytes from the File.
func (f *ProgressFile) Read(b []byte) (n int, err error) {
	n, e := f.file.Read(b)
	f.readSize += int64(n)
	return n, e
}

// Close closes file
func (f *ProgressFile) Close() error {
	return f.file.Close()
}

// ReadSize gets read size for the file
func (f *ProgressFile) ReadSize() int64 {
	return f.readSize
}

// TotalSize gets total size for the file
func (f *ProgressFile) TotalSize() int64 {
	return f.totalSize
}

// Progress gets progress in 0.0~1.0 range
func (f *ProgressFile) Progress() float64 {
	return float64(f.readSize) / float64(f.totalSize)
}
