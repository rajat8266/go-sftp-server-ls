package handler

import (
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
)

// FileLister is a custom type for listing file information with support for offset-based access.
type FileLister struct {
	files []os.FileInfo
}

// ListAt copies file information into the provided slice starting from the given offset.
// It returns the number of copied entries and an error if the end of the list is reached.
func (fl *FileLister) ListAt(fileList []os.FileInfo, offset int64) (int, error) {
	if int(offset) >= len(fl.files) {
		return 0, nil // No more files to list
	}

	n := copy(fileList, fl.files[offset:])
	return n, nil
}

// FileLister returns a list of files in the given directory.
func (h *SftpHandler) FileLister(dirPath string) (*FileLister, error) {
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	fileInfos := make([]os.FileInfo, 0, len(dirEntries))
	for _, entry := range dirEntries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		fileInfos = append(fileInfos, info)
	}

	return &FileLister{files: fileInfos}, nil
}

// FileStat returns the file statistics for the given file.
func (h *SftpHandler) FileStat(filename string) (*FileLister, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	return &FileLister{files: []os.FileInfo{stat}}, nil
}

// setFilePath sets the file path for the request by joining the base path, user, and cleaned file path.
func (h *SftpHandler) setFilePath(r *sftp.Request) {
	r.Filepath = filepath.Join(BasePath, h.Username, filepath.Clean(r.Filepath))
}

// fetch opens a file with the specified mode.
func (h *SftpHandler) fetch(path string, mode int) (*os.File, error) {
	return os.OpenFile(path, mode, 0777)
}
