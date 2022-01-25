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

package logrecord

import (
	"fmt"
	"time"
)

// ErrLineRec is a general error which can originate in any of
// the watched services
type ErrLineRec struct {
	service string
	errType string
	ts      time.Time
}

func (e *ErrLineRec) Tags() map[string]string {
	return map[string]string{}
}

func (e *ErrLineRec) Fields() map[string]interface{} {
	return map[string]interface{}{
		"errType": e.errType,
	}
}

func (e *ErrLineRec) Time() time.Time {
	return e.ts
}

func (e *ErrLineRec) Measurement() string {
	return e.service
}

func (e *ErrLineRec) String() string {
	return fmt.Sprintf("ErrLineRec{service: %s, ts: %v, type: %s}", e.service, e.ts, e.errType)
}

func NewErrLineRec(service string, ts time.Time, errType string) *ErrLineRec {
	return &ErrLineRec{
		service: service,
		ts:      ts,
		errType: errType,
	}
}
