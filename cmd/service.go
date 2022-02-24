package cmd

import (
	"context"
	"github.com/pkg/errors"
	"github.com/viant/smirror"
	"github.com/viant/smirror/cmd/build"
	"github.com/viant/smirror/cmd/mirror"
	"github.com/viant/smirror/cmd/validate"
	"github.com/viant/smirror/contract"
	"github.com/viant/afs"
	"sync/atomic"
)

//Service represents a client service
type Service interface {
	//Build build a rule for cli options
	Build(ctx context.Context, request *build.Request) error
	//Validate check rule either build or with specified URL
	Validate(ctx context.Context, request *validate.Request) error
	//Load start load process for specified source and rule
	Mirror(ctx context.Context, request *mirror.Request) (*mirror.Response, error)
	//Stop stop service
	Stop()
}


type service struct {
	config        *smirror.Config
	mirrorService smirror.Service
	fs            afs.Service
	stopped       int32
	stopChan      chan bool
	requestChan   chan *contract.Request
	responseChan  chan *contract.Response
}

//Stop stops service
func (s *service) Stop() {
	if atomic.CompareAndSwapInt32(&s.stopped, 0, 1) {
		for i := 0; i < 2; i++ {
			s.stopChan <- true
		}
	}
}

//New creates a service
func New(projectID string) (Service, error) {
	ctx := context.Background()
	cfg, err := NewConfig(ctx, projectID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create config")
	}
	tailService, err := smirror.New(ctx, cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create scanFiles service")
	}
	return &service{
		config:        cfg,
		fs:            afs.New(),
		mirrorService: tailService,
		requestChan:   make(chan *contract.Request, processingRoutines),
		responseChan:  make(chan *contract.Response, processingRoutines),
		stopChan:      make(chan bool, 2),
	}, nil
}

