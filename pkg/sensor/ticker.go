// Copyright 2017 Capsule8, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sensor

import (
	"time"

	api "github.com/VikasBhumca2006/capsule8/api/v0"

	"github.com/VikasBhumca2006/capsule8/pkg/stream"
)

type ticker struct {
	ctrl     chan interface{}
	data     chan interface{}
	sensor   *Sensor
	filter   *api.TickerEventFilter
	duration time.Duration
	ticker   *stream.Stream
}

func (t *ticker) newTickerEvent(tick time.Time) *api.Event {
	e := t.sensor.NewEvent()
	e.Event = &api.Event_Ticker{
		Ticker: &api.TickerEvent{
			Seconds:     tick.Unix(),
			Nanoseconds: tick.UnixNano(),
		},
	}

	return e
}

func newTickerSource(sensor *Sensor, filter *api.TickerEventFilter) (*stream.Stream, error) {
	// Each call to New creates a new session with the Sensor. It is the
	// Sensor's responsibility to handle all of its sessions in the most
	// high-performance way possible. For example, a Sensor may install
	// kernel probes for the union of all sessions, but then demux the
	// results through individual goroutines forwarding events over
	// their own channels.

	duration := time.Duration(filter.Interval)
	t := &ticker{
		ctrl:     make(chan interface{}),
		data:     make(chan interface{}),
		sensor:   sensor,
		filter:   filter,
		duration: duration,
		ticker:   stream.Ticker(duration),
	}

	go func() {
		for {
			select {
			case _, ok := <-t.ctrl:
				if !ok {
					close(t.data)
					return
				}

			case e, ok := <-t.ticker.Data:
				if ok {
					tick := e.(time.Time)
					ev := t.newTickerEvent(tick)
					t.data <- ev

				} else {
					return
				}
			}
		}
	}()

	return &stream.Stream{
		Ctrl: t.ctrl,
		Data: t.data,
	}, nil
}
