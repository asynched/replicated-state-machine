package fsm

import (
	"encoding/json"
	"io"
	"log"

	"github.com/hashicorp/raft"
)

type Changeset struct {
	Value int `json:"value"`
}

type StateMachine struct {
	Value int
}

func (s *StateMachine) Apply(l *raft.Log) interface{} {
	data := Changeset{}
	_ = json.Unmarshal(l.Data, &data)

	s.Value = data.Value

	log.Println("Value now is:", s.Value)

	return nil
}

func (s *StateMachine) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (s *StateMachine) Restore(rc io.ReadCloser) error {
	return nil
}

func NewStateMachine() *StateMachine {
	return &StateMachine{
		Value: 0,
	}
}
