/*
 * Copyright 2018 Shanghai Junzheng Network Technology Co.,Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *	   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dump

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/hellobike/amazonriver/conf"
	"github.com/hellobike/amazonriver/handler"
)

// Dumper dump database
type Dumper struct {
	pgDump string
	sub    *conf.Subscribe
}

// New create a Dumper
func New(pgDump string, sub *conf.Subscribe) *Dumper {
	if pgDump == "" {
		pgDump = "pg_dump"
	}
	path, _ := exec.LookPath(pgDump)
	return &Dumper{pgDump: path, sub: sub}
}

// Dump database with snapshot, parse sql then write to handler
func (d *Dumper) Dump(snapshotID string, h handler.Handler) error {

	if d.pgDump == "" {
		return nil
	}

	args := make([]string, 0, 16)

	// Common args
	args = append(args, fmt.Sprintf("--host=%s", d.sub.PGConnConf.Host))
	args = append(args, fmt.Sprintf("--port=%d", d.sub.PGConnConf.Port))

	args = append(args, fmt.Sprintf("--username=%s", d.sub.PGConnConf.User))

	args = append(args, d.sub.PGConnConf.Database)
	args = append(args, "--data-only")
	args = append(args, "--column-inserts")

	args = append(args, fmt.Sprintf("--schema=%s", d.sub.PGConnConf.Schema))
	for _, rule := range d.sub.Rules {
		args = append(args, fmt.Sprintf(`--table=%s`, rule.Table))
	}
	args = append(args, fmt.Sprintf("--snapshot=%s", snapshotID))

	cmd := exec.Command(d.pgDump, args...)
	if d.sub.PGConnConf.Password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", d.sub.PGConnConf.Password))
	}
	r, w := io.Pipe()

	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	errCh := make(chan error)
	parser := newParser(r)
	go func() {
		err := parser.parse(h)
		errCh <- err
	}()

	err := cmd.Run()
	w.CloseWithError(err)

	err = <-errCh
	return err
}
