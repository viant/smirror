package base

import (
	"context"
	"github.com/viant/afs"
	"github.com/viant/afs/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const replayPathElement = "/_replay_/"

//Replayer represents a asset replayed by re-triggering original even by moving storage operation
type Replayer struct {
	*sync.WaitGroup
	routines   int
	hasError   int32
	errChannel chan error
	closed     int32
	fs         afs.Service
	replayURLs chan string
}

func (d *Replayer) replay(ctx context.Context, sourceURL string) {
	defer d.Done()
	bucket := url.Host(sourceURL)
	destURL := strings.Replace(sourceURL, bucket, bucket+replayPathElement, 1)
	e := d.fs.Move(ctx, sourceURL, destURL)
	if e != nil {
		if atomic.CompareAndSwapInt32(&d.hasError, 0, 1) {
			d.errChannel <- e
		}
		return
	}

	if e := d.fs.Move(ctx, destURL, sourceURL); e != nil {
		if atomic.CompareAndSwapInt32(&d.hasError, 0, 1) {
			d.errChannel <- e
		}
	}
}

//Schedule scheduler replay
func (d *Replayer) Schedule(URL string) {
	if URL == "" {
		return
	}
	d.WaitGroup.Add(1)
	d.replayURLs <- URL
}

//Wait for all replay completion
func (d *Replayer) Wait() (err error) {
	time.Sleep(30 * time.Second)
	d.WaitGroup.Wait()
	atomic.StoreInt32(&d.closed, 1)
	for i := 0; i < d.routines; i++ {
		d.replayURLs <- ""
	}
	if atomic.LoadInt32(&d.hasError) == 1 {
		err = <-d.errChannel
	}
	defer close(d.errChannel)
	defer close(d.replayURLs)
	return err
}

//Run start mover go routines
func (d *Replayer) Run(ctx context.Context, routines int) {
	d.routines = routines
	d.replayURLs = make(chan string, routines)
	for i := 0; i < routines; i++ {
		d.WaitGroup.Add(1)
		go func() {
			d.WaitGroup.Done()
			for atomic.LoadInt32(&d.closed) == 0 {
				replayURL := <-d.replayURLs
				if replayURL == "" {
					return
				}
				d.replay(ctx, replayURL)
			}
		}()
	}
}

//NewMover create a mover
func NewReplayer(fs afs.Service) *Replayer {
	return &Replayer{
		WaitGroup:  &sync.WaitGroup{},
		errChannel: make(chan error, 1),
		fs:         fs,
	}
}
