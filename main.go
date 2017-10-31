// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"fmt"
	"net/http"
	"time"
)

func main() {
	c := make(chan App)
	initApi(c)
	initDb("db")
	go runChecksApp(c)
	Serve(8080)
}

func runChecksApp(c chan App) {
	apps := getAllApps()

	for _, a := range apps {
		go func(a App) {
			runHttpCheck(a, c)
		}(a)
	}

	for a := range c {
		go runHttpCheck(a, c)
	}
}

func runHttpCheck(a App, c chan App) {
	_, err := http.Get(a.URL)

	if err != nil {
		a.Status = "down"
	} else {
		a.Status = "up"
	}

	fmt.Println("App", a.URL, "is", a.Status)

	updateApp(a.Name, a)

	time.Sleep(time.Second * time.Duration(a.PollTime))
	c <- a
}
