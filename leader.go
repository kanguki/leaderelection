package leaderelection

// Le wraps a leader election with the main purpose is to decide
// which node is the leader
type Le interface {
	//AmITheLeader checks if current node is the leader
	AmITheLeader(timeoutDecide int) bool
	//Close cleans resources
	Close()
}
