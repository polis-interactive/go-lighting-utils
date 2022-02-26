package main

import (
	"github.com/polis-interactive/go-lighting-utils/pkg/graphicsShader"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type shaderRunner struct {
	wg         *sync.WaitGroup
	shutdowns  chan struct{}
	mu         *sync.RWMutex
	gs         *graphicsShader.GraphicsShader
	ud         graphicsShader.UniformDict
	lastShader int
	startTime  time.Time
}

func (sr *shaderRunner) start() {
	sr.mu = &sync.RWMutex{}
	sr.ud = make(graphicsShader.UniformDict)
	sr.wg = &sync.WaitGroup{}
	sr.shutdowns = make(chan struct{})
	sr.wg.Add(1)
	sr.lastShader = 0
	go sr.runShader()
}

func (sr *shaderRunner) runShader() {

	updateTicker := time.NewTicker(33 * time.Millisecond)
	changeShaderTicker := time.NewTicker(5 * time.Second)

	sr.setupShader()

	defer func(sr *shaderRunner, t1 *time.Ticker, t2 *time.Ticker) {
		t1.Stop()
		t2.Stop()
		sr.gs.Cleanup()
		sr.wg.Done()
	}(sr, updateTicker, changeShaderTicker)

	sr.startTime = time.Now()

	for {
		select {
		case _, ok := <-sr.shutdowns:
			if !ok {
				return
			}
		case <-updateTicker.C:
			sr.doRunShader()
			break
		case <-changeShaderTicker.C:
			sr.doChangeShader()
			break
		}
	}
}

func (sr *shaderRunner) setupShader() {
	gs, err := graphicsShader.NewGraphicsShader(
		"lighting-utils", 800, 600, sr.ud, sr.mu,
	)
	if err != nil {
		panic(err)
	}

	err = gs.AttachShader(graphicsShader.ShaderIdentifier{
		Key:      graphicsShader.ShaderKey(rune(0)),
		Filename: "basic",
	})
	if err != nil {
		panic(err)
	}

	err = gs.AttachShader(graphicsShader.ShaderIdentifier{
		Key:      graphicsShader.ShaderKey(rune(1)),
		Filename: "slate-1",
	})
	if err != nil {
		panic(err)
	}
	sr.gs = gs
}

func (sr *shaderRunner) doChangeShader() {
	log.Println("SWITCHING :D")
	sr.lastShader = (sr.lastShader + 1) % 2
	err := sr.gs.SetShader(graphicsShader.ShaderKey(rune(sr.lastShader)))
	if err != nil {
		panic(err)
	}
}

func (sr *shaderRunner) doRunShader() {

	duration := time.Since(sr.startTime)
	sr.ud["time"] = float32(duration.Seconds())

	err := sr.gs.ReloadShader()
	if err != nil {
		panic(err)
	}

	err = sr.gs.RunShader()
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func (sr *shaderRunner) stop() {
	close(sr.shutdowns)
	sr.wg.Wait()
}

func main() {
	sr := &shaderRunner{}
	sr.start()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	sr.stop()
}
