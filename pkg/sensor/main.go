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
	"github.com/VikasBhumca2006/capsule8/pkg/config"
	"github.com/VikasBhumca2006/capsule8/pkg/services"
	"github.com/golang/glog"
)

// Main is the main entrypoint for the sensor
func Main() {
	manager := services.NewServiceManager()
	if len(config.Global.ProfilingAddr) > 0 {
		service := services.NewProfilingService(
			config.Global.ProfilingAddr)
		manager.RegisterService(service)
	}
	if len(config.Sensor.ServerAddr) > 0 {
		sensor, err := NewSensor()
		if err != nil {
			glog.Fatalf("Could not create sensor: %s", err.Error())
		}
		if err := sensor.Start(); err != nil {
			glog.Fatalf("Could not start sensor: %s", err.Error())
		}
		defer sensor.Stop()
		service := NewTelemetryService(sensor, config.Sensor.ServerAddr)
		manager.RegisterService(service)
	}

	manager.Run()
}
