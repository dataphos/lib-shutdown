// Copyright 2024 Syntio Ltd.
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

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dataphos/lib-shutdown/pkg/graceful"
)

func main() {
	secondsRemaining := 10
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(secondsRemaining)*time.Second)
	defer cancel()

	ctx = graceful.WithSignalShutdown(ctx)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("context cancelled, leaving")
			return
		case <-time.After(1 * time.Second):
			secondsRemaining -= 1
			fmt.Println(secondsRemaining, "seconds remaining")
		}
	}
}
