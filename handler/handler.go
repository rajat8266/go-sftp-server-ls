package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/pkg/sftp"
)

// SftpHandler handles SFTP operations.
type SftpHandler struct {
	mu       sync.Mutex
	Username string
}

// NewSftpHandler initializes and returns SFTP handlers for different file operations.
func NewSftpHandler(ctx context.Context) (*sftp.Handlers, error) {
	handler := &SftpHandler{
		Username: User,
	}

	return &sftp.Handlers{
		FileGet:  handler,
		FilePut:  handler,
		FileCmd:  handler,
		FileList: handler,
	}, nil
}

// Fileread reads a file and returns an io.ReaderAt interface.
func (h *SftpHandler) Fileread(r *sftp.Request) (io.ReaderAt, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.setFilePath(r)

	file, err := h.fetch(r.Filepath, os.O_RDONLY)
	if err != nil {
		log.Printf("error reading file: %v", err)
		return nil, err
	}
	return file, nil
}

// Filewrite creates a file and returns an io.WriterAt interface.
func (h *SftpHandler) Filewrite(r *sftp.Request) (io.WriterAt, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.setFilePath(r)

	file, err := h.fetch(r.Filepath, os.O_CREATE|os.O_RDWR)
	if err != nil {
		log.Printf("error writing file: %v", err)
		return nil, err
	}
	return file, nil
}

// Filecmd performs various file operations.
//
// Supported operations:
// - Remove: Removes a file.
// Unsupported operations:
// - Setstat
// - Rename
// - Mkdir
// - Rmdir
// - Symlink
func (h *SftpHandler) Filecmd(r *sftp.Request) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	switch r.Method {
	case "Remove":
		h.setFilePath(r)
		err := os.Remove(r.Filepath)
		if err != nil {
			log.Printf("error removing file: %v", err)
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported file command: %s", r.Method)
	}
}

// Filelist performs file listing and stat operations.
//
// Supported operations:
// - List: Returns a list of files.
// - Stat: Returns file stats.
// Unsupported operation:
// - Readlink
func (h *SftpHandler) Filelist(r *sftp.Request) (sftp.ListerAt, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.setFilePath(r)

	switch r.Method {
	case "List":
		list, err := h.FileLister(r.Filepath)
		if err != nil {
			log.Printf("error listing files: %v", err)
			return nil, err
		}
		return list, nil
	case "Stat":
		stats, err := h.FileStat(r.Filepath)
		if err != nil {
			log.Printf("error getting file stats: %v", err)
			return nil, err
		}
		return stats, nil
	default:
		return nil, fmt.Errorf("unsupported file list command: %s", r.Method)
	}
}
