//go:build nats
// +build nats

package leaderelection

import (
	"sync"
	"testing"

	"github.com/kanguki/log"
)

func TestNatsLe(t *testing.T) {
	count := 3
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			node, _ := NewNatsLe("test", count, []string{"nats://127.0.0.1:4222"})
			leader, err := node.WhoIsTheLeader(5)
			if err != nil {
				log.Log("%v", err)
			}
			log.Log("node %s on term %v has leader %s", node.Node.Id(), node.Node.CurrentTerm(), leader)
			log.Log("%s: AmITheLeader? %v", node.Node.Id(), node.AmITheLeader(5))
			// for {
			// 	select {
			// 	case <-time.After(time.Second):
			// 		log.Log("%s: AmITheLeader? %v", node.Node.Id(), node.AmITheLeader(5))
			// 		log.Log(node.Node.Leader())
			// 	}
			// }
			node.Node.Close()
		}()
	}
	wg.Wait()
}
