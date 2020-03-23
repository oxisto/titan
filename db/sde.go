/*
Copyright 2020 Christian Banse

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package db

import (
	"fmt"
	"os/exec"
	"time"
)

type StaticDataExport struct {
	Version int32
	Server  string
}

func (t StaticDataExport) ID() int32 {
	return t.Version
}

func (t StaticDataExport) ExpiresOn() *time.Time {
	return nil
}

func (t StaticDataExport) SetExpire(time *time.Time) {

}

func (t StaticDataExport) HashKey() string {
	return fmt.Sprintf("sde:%d", t.ID())
}

// RunSDERestoreScript executes the SDE restore script
func RunSDERestoreScript(version int32, server string) {
	log.Infof("Importing SDE %d...", version)

	cmd := exec.Command("./restore.sh", host)
	cmd.Wait()
}
