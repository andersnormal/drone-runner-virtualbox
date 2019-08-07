package run

import (
	"github.com/drone/drone-go/drone"
	log "github.com/sirupsen/logrus"
)

// Stage ...
type Stage struct {
	stage  *drone.Stage
	logger *log.Entry
}

// NewStage ...
func NewStage(stage *drone.Stage, ll *log.Entry) *Stage {
	s := new(Stage)
	s.stage = stage
	s.logger = ll

	return s
}

// Run ...
func (s *Stage) Run() error {
	ll := s.log().WithFields(log.Fields{
		"stage": s.stage.ID,
	})

	ll.Printf("starting run")

	return nil
}

func (s *Stage) log() *log.Entry {
	return s.logger
}
