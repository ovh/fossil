package writer

import (
	"os"
	"path"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ovh/fossil/core"
)

// FileWriter write GTS onto files
type FileWriter struct {
	SourceDir       string
	currentFile     *os.File
	currentFileName string
}

// NewWriter return an instanciated FileWriter
func NewWriter(dir string) *FileWriter {

	f := &FileWriter{
		SourceDir: dir,
	}

	go func() {
		tick := time.NewTicker(5 * time.Second)
		for range tick.C {
			f.fileRotate()
		}
	}()

	f.fileRotate()

	return f
}

func (fw *FileWriter) Write(in chan *core.GTS) {

	for gts := range in {
		log.Debug(gts)
		_, err := fw.currentFile.Write(gts.Encode())
		if err != nil {
			log.WithError(err).Error("Failed to write metrics in file")
		}
	}
}

func (fw *FileWriter) fileRotate() error {

	now := time.Now()
	fileName := strconv.Itoa(int(now.UnixNano()))

	newFile, err := os.Create(path.Join(fw.SourceDir, fileName+".tmp"))
	if err != nil {
		return err
	}

	oldFileName := fw.currentFileName

	fw.currentFile = newFile
	fw.currentFileName = fileName

	if oldFileName != "" {
		return os.Rename(path.Join(fw.SourceDir, oldFileName+".tmp"), path.Join(fw.SourceDir, oldFileName+".metrics"))
	}
	return nil
}
