# go-mpv

Toolbox to automate MPV media player. Builds on top of the [mpvipc](https://github.com/dexterlb/mpvipc) by DexterLB.

## Installation

```bash
go get github.com/piitschy/go-mpv
```

## Usage

```go 
package main

import (
	"log"
	"os"

	"github.com/Piitschy/gompv"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <file>")
	}

	session := gompv.NewSession(os.Args[1])
	session.SetOSC(false)
	session.SetNoInputDefaultBindings(true)
	for _, arg := range os.Args[2:] {
		session.AddAudioSource(arg)
	}
	session.AddGlobalAudioFilter("amix")
	err := session.Start()

	if err != nil {
		log.Fatalf("failed to start session: %v", err)
	}

	defer func() {
		if err := session.Stop(); err != nil {
			log.Fatalf("failed to stop session: %v", err)
		}
	}()

	client := session.Client()

	events, stopListening := client.NewEventListener()

	err = client.Pause()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("paused!")

	go func() {
		client.WaitUntilClosed()
		stopListening <- struct{}{}
	}()

	for event := range events {
		log.Printf("received event: %s", event.Name)
	}

	log.Printf("mpv closed socket")
}
```
