package writer

import (
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/ovh/fossil/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// FileWriter write GTS onto files
type FileWriter struct {
	Metrics    []*core.GTS
	SourceDir  string
	batchCount int64
	batchFull  chan struct{}
	sync.Mutex
}

// NewWriter return an instanciated FileWriter
func NewWriter(dir string) *FileWriter {
	fw := &FileWriter{
		SourceDir:  dir,
		batchFull:  make(chan struct{}),
		batchCount: 0,
		Metrics:    []*core.GTS{},
	}

	go func() {
		tick := time.NewTicker(viper.GetDuration("timeout") * time.Second)

		for range tick.C {
			fw.Lock()
			if fw.batchCount > 0 {
				err := fw.flush()
				if err != nil {
					log.WithFields(log.Fields{"error": err.Error()}).Error("Cannot flush datapoints into file")
				}
			}
			fw.Unlock()
		}
	}()

	return fw
}

func (fw *FileWriter) Write(gts *core.GTS) {
	fw.Lock()
	defer fw.Unlock()
	fw.Metrics = append(fw.Metrics, gts)
	fw.batchCount++

	if fw.batchCount >= int64(viper.GetInt("batch")) {
		err := fw.flush()
		if err != nil {
			log.WithFields(log.Fields{"error": err.Error()}).Error("Cannot flush datapoints into file")
		}
	}
}

func (fw *FileWriter) flush() error {

	now := time.Now()
	fileName := strconv.Itoa(int(now.UnixNano()))
	newFile, err := os.Create(path.Join(fw.SourceDir, fileName+".tmp"))

	if err != nil {
		return err
	}

	for _, gts := range fw.Metrics {
		_, err := newFile.Write(gts.Encode())
		if err != nil {
			log.WithError(err).Error("Failed to write metrics in file")
			return err
		}
	}

	fw.batchCount = 0
	fw.Metrics = []*core.GTS{}

	return os.Rename(path.Join(fw.SourceDir, fileName+".tmp"), path.Join(fw.SourceDir, fileName+".metrics"))
}
