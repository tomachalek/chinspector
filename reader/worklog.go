// Copyright 2019 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2019 Institute of the Czech National Corpus,
//                Faculty of Arts, Charles University
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

// WorklogRecord stores inode & seek position of last read operation
type WorklogRecord struct {
	Inode int64 `json:"inode"`
	Seek  int64 `json:"seek"`
}

type updateRequest struct {
	Value WorklogRecord
}

// Worklog provides functions to store/retrieve information about
// file reading operations to be able to continue in case of an
// interruption
type Worklog struct {
	rec         WorklogRecord
	updRequests chan updateRequest
}

// Init initializes the worklog. It must be called before any other
// operation.
func (w *Worklog) Init() error {
	w.updRequests = make(chan updateRequest)
	go func() {
		for req := range w.updRequests {
			w.rec = req.Value
		}
	}()
	return nil
}

// Close cleans up worklog for safe exit
func (w *Worklog) Close() {
	w.rec = WorklogRecord{}
	if w.updRequests != nil {
		close(w.updRequests)
	}
}

// UpdateFileInfo adds individual app reading position info. Please
// note that this does not save the worklog.
func (w *Worklog) UpdateFileInfo(inode int64, seek int64) {
	w.updRequests <- updateRequest{Value: WorklogRecord{Inode: inode, Seek: seek}}
}

// GetData retrieves reading info for a provided app
func (w *Worklog) GetData() WorklogRecord {
	return w.rec
}

// NewWorklog creates a new Worklog instance. Please note that
// Init() must be called before you can begin using the worklog.
func NewWorklog() *Worklog {
	return &Worklog{rec: WorklogRecord{}}
}
