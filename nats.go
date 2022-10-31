package leaderelection

import (
	"fmt"
	"time"

	klog "github.com/kanguki/log"
	"github.com/nats-io/graft"
	"github.com/nats-io/nats.go"
)

// NatsLe is a leader elector based on nats
type NatsLe struct {
	//Node saves state of current node and its leader
	Node *graft.Node
}

// AmITheLeader check if current node is the leader.
//
// If nats server dies, old leader still marks it as leader while other doesnt keep leader's information
func (nle NatsLe) AmITheLeader(timeoutDecide int) bool {
	leader, err := nle.WhoIsTheLeader(timeoutDecide)
	if err != nil {
		klog.Log("WhoIsTheLeader error: %v", err)
	}
	if leader == nle.Node.Id() {
		return true
	}
	return false
}

// WhoIsTheLeader gets leader of the cluster.
//
// After timeout without a leader in cluster, it returns error.
//
// Leader is automatically updated if old leader crashes.
func (nle NatsLe) WhoIsTheLeader(timeout int) (string, error) {
	stop := make(chan bool, 1)
	go func() {
		<-time.After(time.Duration(timeout) * time.Second)
		stop <- true
	}()
	for nle.Node.Leader() == "" {
		select {
		case <-stop:
			return "", fmt.Errorf("timeout electing leader")
		default:
			time.Sleep(time.Second)
		}
	}
	return nle.Node.Leader(), nil
}

//NewNatsLe creates a leader elector with nats. each cluster has only 1 leader
//
// electionName is used in nats engine to separeate with other elections
//
// size is number of candidates in election, should be at least 3.
// if more than size nodes join 1 cluster, there's still only 1 leader.
// leader is decided based on size, not on actual connected nodes,
// i.e if size = 3 and 7 nodes connected, leader is decided by 2 nodes only!
//
// natsQuorum is list of nats servers, e.g. nats://127.0.0.1:4222. reference: https://github.com/nats-io/go-nats/blob/master/example_test.go
func NewNatsLe(electionName string, size int, natsQuorum []string) (NatsLe, error) {
	klog.Log("starting a new nats le named %s using nats quorum %v", electionName, natsQuorum)
	ci := graft.ClusterInfo{Name: electionName, Size: size}
	do := nats.GetDefaultOptions()
	do.Servers = natsQuorum
	rpc, err := graft.NewNatsRpc(&do)
	if err != nil {
		klog.Log("error creating NewNatsRpc: %v", err)
		return NatsLe{}, err
	}
	errChan := make(chan error)
	stateChangeChan := make(chan graft.StateChange)
	handler := graft.NewChanHandler(stateChangeChan, errChan)
	node, err := graft.New(ci, handler, rpc, "/tmp/graft.log")
	if err != nil {
		klog.Log("error joining NewNatsLe: %v", err)
		return NatsLe{}, err
	}
	go func() {
		for {
			select {
			case sc := <-stateChangeChan:
				klog.Debug("node %v's state changed from %v to %v", node.Id(), sc.From, sc.To)
			case err := <-errChan:
				klog.Debug("node %v received error %v", node.Id(), err)
			}
		}
	}()
	return NatsLe{Node: node}, err
}
