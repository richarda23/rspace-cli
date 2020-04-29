package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"github.com/dustin/go-humanize"
	"strings"
)
// holds information about scanned local files, and whether they have been processed yet
type scannedFileInfo struct {
	Path string
	Info os.FileInfo
	Uploaded bool
}

// validates that file paths entered as command line arguments are readable
func validateInputFilePaths ( paths []string ) {
	for _, filePath := range paths {
		var err error
		filePath, err = filepath.Abs(filePath)
		_, err = os.Stat(filePath)
		if err != nil {
			exitWithErr(err)
		}
	}
}

// filters scannedFileInfo
type acceptFileFilter func (info *scannedFileInfo) bool

// default acceptFileFilter than accepts all files
func acceptAll () acceptFileFilter {
	return func( info *scannedFileInfo) bool {
		return true
	}
}
// default acceptFileFilter than rejects all files
func rejectAll () acceptFileFilter {
	return func( info *scannedFileInfo) bool {
		return false
	}
}

func acceptMsDoc () acceptFileFilter {
	return func( info *scannedFileInfo) bool {
		ext := strings.ToLower(filepath.Ext(info.Path))
		return ext == ".doc" || ext == ".odt" || ext == ".docx"
	}
}

func scanFiles (paths []string, recurse bool, accept acceptFileFilter) []*scannedFileInfo {
	var filesToUpload []*scannedFileInfo = make([]*scannedFileInfo,0)
	for _, filePath := range paths {
		messageStdErr("processing " + filePath)
		filePath, _ = filepath.Abs(filePath)
		fileInfo, _ := os.Stat(filePath)
		if fileInfo.IsDir() {
			messageStdErr("Scanning for files in " +  fileInfo.Name())
			if recurse {
				filepath.Walk(filePath, visit(&filesToUpload) )
			} else {
				readSingleDir(filePath, &filesToUpload)
			}
		} else {
			info,_:= os.Stat(filePath)
			filesToUpload = append(filesToUpload, &scannedFileInfo{filePath,info,false})
		}
	}
	// now filter
	var rc []*scannedFileInfo = make([]*scannedFileInfo,0)
	for _,v := range filesToUpload {
		if accept(v) {
			rc = append(rc, v)
		}
	}
	return rc
}

func sumFileSize(toUpload []*scannedFileInfo) uint64 {
	var sum int64 = 0
	for _,v :=range toUpload {
		sum += v.Info.Size()
	}
	return uint64(sum)
}

func sumFileSizeHuman(toUpload []*scannedFileInfo) string {
	return humanizeBytes(sumFileSize(toUpload))
}

func humanizeBytes(byteSize uint64) string {
	return humanize.Bytes(byteSize)
}

// reads non . files from a single folder
func readSingleDir(filePath string, files *[]*scannedFileInfo) {

	fileInfos,_:= ioutil.ReadDir(filePath)
	for _,inf:=range fileInfos {
		if !inf.IsDir() && !isDot(inf) {
				path := filePath + string(os.PathSeparator) +inf.Name()
				*files = append(*files, &scannedFileInfo{path,inf,false})
		}
	}
 }
func visit (files *[]*scannedFileInfo) filepath.WalkFunc {
	return func  (path string, info os.FileInfo, err error) error {
		// always   ignore '.' folders, don't descend

		if info.IsDir() && isDot(info) {
			messageStdErr("Skipping .folder " + path)
			return filepath.SkipDir
		}
		// always add non . files
		if !info.IsDir() && !isDot(info) {
			*files = append(*files, &scannedFileInfo{path,info,false})
			return nil
		}
		return nil
	}
}
func isDot(info os.FileInfo) bool {
	//return filepath.Base(info.Name())[0] == '.'
	match,_ :=  regexp.MatchString("^\\.[A-Za-z0-9\\-_]+", info.Name())
	return match
}
