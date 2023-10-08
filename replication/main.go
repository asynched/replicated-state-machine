package replication

import (
	"fmt"
	"time"

	"github.com/asynched/replicated-state-machine/config"
	"github.com/asynched/replicated-state-machine/fsm"
	"github.com/hashicorp/raft"
)

func GetRaft(config config.Config) (*raft.Raft, error) {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	transport, err := raft.NewTCPTransport(addr, nil, 3, time.Second*10, nil)

	if err != nil {
		return nil, err
	}

	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(addr)

	logs := raft.NewInmemStore()
	stable := raft.NewInmemStore()
	snapshot := raft.NewInmemSnapshotStore()
	fsm := fsm.NewStateMachine()

	r, err := raft.NewRaft(raftConfig, fsm, logs, stable, snapshot, transport)

	if err != nil {
		return nil, err
	}

	if !config.Bootstrap {
		return r, nil
	}

	servers := []raft.Server{
		{
			ID:      raft.ServerID(addr),
			Address: transport.LocalAddr(),
		},
	}

	for _, peer := range config.Peers {
		addr := fmt.Sprintf("%s:%d", peer.Host, peer.Port)

		servers = append(servers, raft.Server{
			ID:      raft.ServerID(addr),
			Address: raft.ServerAddress(addr),
		})
	}

	if err := r.BootstrapCluster(raft.Configuration{Servers: servers}).Error(); err != nil {
		return nil, err
	}

	return r, nil
}
