package main

import (
	"context"
	"flag"
	"gitcom.com/valentin-carl/f4s/pkg/loadgenerator"
	"github.com/VividCortex/multitick"
	"github.com/fatih/color"
	"log"
	"sync"
	"time"
)

var (
	dst  = flag.String("dst", "", "target URL")
	dur  = flag.Duration("t", time.Minute, "experiment duration")
	nps  = flag.Int("n", 1, "number of parallel streams")
	size = flag.Int("size", 1024, "message size in bytes")
)

func main() {

	// TODO
	//  at the moment, there is a set of workers that send messages in regular intervals
	//  this might have to be adjusted depending on what will be measured in the experiments (e.g., throughput)
	// TODO
	//  adjust what to workers do to desired load pattern/use case

	log.Print(color.GreenString("starting load generator"))

	flag.Parse()
	log.Printf("load generator input: %s %v %d %d", *dst, *dur, *nps, *size)

	mt := multitick.NewTicker(time.Second, 0)
	defer mt.Stop()

	ctxAll, cancel := context.WithTimeout(context.Background(), time.Duration(*dur))
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < *nps; i++ {
		ctx := context.WithValue(ctxAll, loadgenerator.Key0, struct {
			id   int
			seed int64
			size int
		}{
			i,
			int64(i),
			*size,
		})
		go func() {
			wg.Add(1)
			loadgenerator.GenerateLoad(ctx, mt.Subscribe(), dst)
			wg.Done()
		}()
	}

	<-ctxAll.Done()
	wg.Wait()
	log.Print(color.GreenString("load generator done"))
}
