package logFileRotator

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type Options struct {
	fileName   string
	maxSize    int
	maxAge     int
	maxBackups int
	localTime  bool
	timeFormat string
}

type Engine struct {
	options Options
	size    int64
	file    *os.File
	mu      sync.Mutex
}

func (that *Engine) openExistingOrNew(writeLen int) error {
	info, err := os.Stat(that.options.fileName)
	if os.IsNotExist(err) {
		return that.openNew()
	}
	if err != nil {
		return fmt.Errorf("error getting log file info: %s", err)
	}

	if info.Size()+int64(writeLen) >= int64(that.options.maxSize) {
		return that.rotate()
	}

	file, err := os.OpenFile(that.options.fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return that.openNew()
	}

	that.file = file
	that.size = info.Size()
	return nil
}

func (that *Engine) compress(name string) {
	inFile := name
	outFile := name + ".gz"
	zipHandle, err := os.OpenFile(outFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		that.error("Compress: Opening file:", err)
		return
	}
	defer zipHandle.Close()

	zipWriter, err := gzip.NewWriterLevel(zipHandle, 9)
	if err != nil {
		that.error("Compress: New gzip writer:", err)
		return
	}
	defer zipWriter.Close()

	inReader, err := os.OpenFile(inFile, os.O_RDONLY, 0666)
	if err != nil {
		that.error("Compress: Opening old log file:", err)
		return
	}
	defer inReader.Close()

	_, err = io.Copy(zipWriter, inReader)
	if err != nil {
		that.error("Compress: copy", err)
		return
	}
}

func (that *Engine) removeFile(filename string) {
	err := os.Remove(filename)
	if err != nil {
		that.error("Delete old log fileError :"+filename, err)
		return
	}
}

func (that *Engine) openNew() error {
	err := os.MkdirAll(that.dir(), 0744)
	if err != nil {
		return fmt.Errorf("can't make directories for new logfile: %s", err)
	}

	filename := that.options.fileName
	mode := os.FileMode(0644)
	info, err := os.Stat(filename)
	if err == nil {
		mode = info.Mode()
		newName := that.backupName(filename, that.options.localTime)
		if err := os.Rename(filename, newName); err != nil {
			return fmt.Errorf("can't rename log file: %s", err)
		}

		that.compress(newName)
		that.removeFile(newName)
		if err := chown(filename, info); err != nil {
			return err
		}
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}

	that.file = f
	that.size = 0
	return nil
}

func (that *Engine) Write(p []byte) (n int, err error) {
	that.mu.Lock()
	defer that.mu.Unlock()

	return that.write(p)
}

func (that *Engine) write(p []byte) (n int, err error) {
	writeLen := len(p)
	if writeLen > that.options.maxSize {
		return 0, fmt.Errorf(
			"Record size(%d) exceeds max size of log(%d)", writeLen, that.options.maxSize,
		)
	}

	if that.file == nil {
		if err = that.openExistingOrNew(len(p)); err != nil {
			return 0, err
		}
	}

	if that.size+int64(writeLen) > int64(that.options.maxSize) {
		if err := that.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = that.file.Write(p)
	that.size += int64(n)
	return n, err
}

func (that *Engine) dir() string {
	return filepath.Dir(that.options.fileName)
}

func (that *Engine) Rotate() error {
	that.mu.Lock()
	defer that.mu.Unlock()

	return that.rotate()
}

func (that *Engine) rotate() error {
	if err := that.close(); err != nil {
		return err
	}

	if err := that.openNew(); err != nil {
		return err
	}
	return that.cleanup()
}

func (that *Engine) Close() error {
	that.mu.Lock()
	defer that.mu.Unlock()

	return that.close()
}

func (that *Engine) close() error {
	if that.file == nil {
		return nil
	}

	err := that.file.Close()
	that.file = nil
	return err
}

func (that *Engine) cleanup() error {
	if that.options.maxBackups == 0 && that.options.maxAge == 0 {
		return nil
	}

	files, err := that.oldLogFiles()
	if err != nil {
		return err
	}

	var deletes []logInfo

	if that.options.maxBackups > 0 && that.options.maxBackups < len(files) {
		deletes = files[that.options.maxBackups:]
		files = files[:that.options.maxBackups]
	}
	if that.options.maxAge > 0 {
		diff := time.Duration(int64(24*time.Hour) * int64(that.options.maxAge))

		cutoff := time.Now().Add(-1 * diff)

		for _, f := range files {
			if f.timestamp.Before(cutoff) {
				deletes = append(deletes, f)
			}
		}
	}

	if len(deletes) == 0 {
		return nil
	}

	go deleteAll(that.dir(), deletes)

	return nil
}

func deleteAll(dir string, files []logInfo) {
	for _, f := range files {
		_ = os.Remove(filepath.Join(dir, f.Name()))
	}
}

func (that *Engine) prefixAndSuffix() (prefix, suffix string) {
	filename := filepath.Base(that.options.fileName)
	suffix = filepath.Ext(filename)
	prefix = filename[:len(filename)-len(suffix)] + "-"
	return prefix, suffix
}

func (that *Engine) timeFromName(filename, prefix, ext string) string {
	if !strings.HasPrefix(filename, prefix) {
		return ""
	}

	filename = filename[len(prefix):]
	if strings.HasSuffix(filename, ext+".gz") {
		ext = ext + ".gz"
	}

	if !strings.HasSuffix(filename, ext) {
		return ""
	}

	filename = filename[:len(filename)-len(ext)]
	return filename
}

func (that *Engine) oldLogFiles() ([]logInfo, error) {
	entries, err := os.ReadDir(that.dir())
	if err != nil {
		return nil, fmt.Errorf("can't read log file directory: %s", err)
	}
	var logFiles byFormatTime

	prefix, ext := that.prefixAndSuffix()
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := that.timeFromName(entry.Name(), prefix, ext)
		if name == "" {
			continue
		}

		t, err := time.Parse(that.options.timeFormat, name)
		if err == nil {
			logFiles = append(logFiles, logInfo{timestamp: t, DirEntry: entry})
		}
	}

	sort.Sort(logFiles)

	return logFiles, nil
}

func (that *Engine) error(msg string, err error) {
	fmt.Println(msg, err.Error())
}

func (that *Engine) backupName(name string, local bool) string {
	dir := filepath.Dir(name)
	filename := filepath.Base(name)
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)]
	t := time.Now()
	if !local {
		t = t.UTC()
	}
	timestamp := t.Format(that.options.timeFormat)
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", prefix, timestamp, ext))
}

type logInfo struct {
	timestamp time.Time
	os.DirEntry
}

type byFormatTime []logInfo

func (b byFormatTime) Less(i, j int) bool {
	return b[i].timestamp.After(b[j].timestamp)
}

func (b byFormatTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byFormatTime) Len() int {
	return len(b)
}
