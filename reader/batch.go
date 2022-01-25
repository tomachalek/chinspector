// Copyright 2021 Tomas Machalek <tomas.machalek@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package reader

import (
	"bufio"
	"log"
	"os"
	"time"

	"chinspector/config"
	"chinspector/logrecord"
)

type FileBatchReader struct {
	processor ChiaLogRecProcessor
	file      *os.File
	inode     int64
	size      int64
}

func (ftw *FileBatchReader) Processor() ChiaLogRecProcessor {
	return ftw.processor
}

// ApplyNewContent calls a provided function to newly added lines
func (ftw *FileBatchReader) ApplyNewContent(onLine func(line string), onDone func(inode int64, seek int64)) error {
	var err error
	ftw.file.Close()
	ftw.file, err = os.Open(ftw.processor.FilePath())
	if err != nil {
		return err
	}
	sc := bufio.NewScanner(ftw.file)
	for sc.Scan() {
		onLine(sc.Text())
	}
	onDone(ftw.inode, ftw.size)
	return nil
}

func RunBatch(conf *config.Props, recProcessor *logrecord.Processor) {
	worklog := NewWorklog()
	err := worklog.Init()
	var rdr *FileBatchReader

	if err != nil {
		log.Print("ERROR: ", err)
		return

	} else {
		currInode, size, _ := getFileProps(recProcessor.FilePath())
		rdr = &FileBatchReader{
			processor: recProcessor,
			inode:     currInode,
			size:      size,
		}
	}

	rdr.Processor().OnCheckStart(time.Now())
	rdr.ApplyNewContent(
		func(v string) {
			rdr.Processor().OnLineRead(v)
		},
		func(inode int64, seek int64) {
			worklog.UpdateFileInfo(inode, seek)
		},
	)
	rdr.Processor().OnCheckStop()
	rdr.Processor().OnQuit()
}
