package SurfTest

import (
	context "context"
	"cse224/proj5/pkg/surfstore"
	"strings"
	"testing"
	"time"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func TestRaftSetLeader(t *testing.T) {
	//Setup
	cfgPath := "./config_files/3nodes.txt"
	test := InitTest(cfgPath, "8080")
	defer EndTest(test)

	// TEST
	leaderIdx := 0
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})

	// heartbeat
	for _, server := range test.Clients {
		server.SendHeartbeat(test.Context, &emptypb.Empty{})
	}

	for idx, server := range test.Clients {
		// all should have the leaders term
		state, _ := server.GetInternalState(test.Context, &emptypb.Empty{})
		if state.Term != int64(1) {
			t.Logf("Server %d should be in term %d", idx, 1)
			t.Fail()
		}
		if idx == leaderIdx {
			// server should be the leader
			if !state.IsLeader {
				t.Logf("Server %d should be the leader", idx)
				t.Fail()
			}
		} else {
			// server should not be the leader
			if state.IsLeader {
				t.Logf("Server %d should not be the leader", idx)
				t.Fail()
			}
		}
	}

	leaderIdx = 2
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})

	// heartbeat
	for _, server := range test.Clients {
		server.SendHeartbeat(test.Context, &emptypb.Empty{})
	}

	for idx, server := range test.Clients {
		// all should have the leaders term
		state, _ := server.GetInternalState(test.Context, &emptypb.Empty{})
		if state.Term != int64(2) {
			t.Logf("Server should be in term %d", 2)
			t.Fail()
		}
		if idx == leaderIdx {
			// server should be the leader
			if !state.IsLeader {
				t.Logf("Server %d should be the leader", idx)
				t.Fail()
			}
		} else {
			// server should not be the leader
			if state.IsLeader {
				t.Logf("Server %d should not be the leader", idx)
				t.Fail()
			}
		}
	}
}

func TestRaftFollowersGetUpdates(t *testing.T) {
	//Setup
	cfgPath := "./config_files/3nodes.txt"
	test := InitTest(cfgPath, "8080")
	defer EndTest(test)

	// TEST
	leaderIdx := 0
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	filemeta1 := &surfstore.FileMetaData{
		Filename:      "testFile1",
		Version:       1,
		BlockHashList: nil,
	}

	test.Clients[leaderIdx].UpdateFile(test.Context, filemeta1)
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	goldenMeta := surfstore.NewMetaStore("")
	goldenMeta.UpdateFile(test.Context, filemeta1)
	goldenLog := make([]*surfstore.UpdateOperation, 0)
	goldenLog = append(goldenLog, &surfstore.UpdateOperation{
		Term:         1,
		FileMetaData: filemeta1,
	})

	for _, server := range test.Clients {
		state, _ := server.GetInternalState(test.Context, &emptypb.Empty{})
		if !SameLog(goldenLog, state.Log) {
			t.Log("Logs do not match")
			t.Fail()
		}
		if !SameMeta(goldenMeta.FileMetaMap, state.MetaMap.FileInfoMap) {
			t.Log("MetaStore state is not correct")
			t.Fail()
		}
	}
}

func TestRaftServerIsCrashable(t *testing.T) {
	t.Log("a request is sent to a crashed server")
	//Setup
	cfgPath := "./config_files/3nodes.txt"
	test := InitTest(cfgPath, "8080")
	defer EndTest(test)

	// TEST
	leaderIdx := 0
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	test.Clients[leaderIdx].Crash(test.Context, &emptypb.Empty{})

	filemeta1 := &surfstore.FileMetaData{
		Filename:      "testFile1",
		Version:       1,
		BlockHashList: nil,
	}

	_, err := test.Clients[leaderIdx].UpdateFile(test.Context, filemeta1)
	if err == nil || !strings.Contains(err.Error(), surfstore.ERR_SERVER_CRASHED.Error()) {
		t.Log("Server should return ERR_SERVER_CRASHED")
		t.Fail()
	}
}

func TestRaftBlockWhenMajorityDown(t *testing.T) {
	t.Log("leader1 gets a request when the majority of the cluster is crashed.")
	//Setup
	cfgPath := "./config_files/3nodes.txt"
	test := InitTest(cfgPath, "8080")
	defer EndTest(test)

	// TEST
	leaderIdx := 0
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	// crash all but the leader
	for idx, server := range test.Clients {
		if idx != leaderIdx {
			server.Crash(test.Context, &emptypb.Empty{})
		}
	}

	filemeta1 := &surfstore.FileMetaData{
		Filename:      "testFile1",
		Version:       1,
		BlockHashList: nil,
	}

	blockChan := make(chan bool)

	go func() {
		_, _ = test.Clients[leaderIdx].UpdateFile(context.Background(), filemeta1)
		blockChan <- false
	}()

	go func() {
		<-time.NewTimer(5 * time.Second).C
		blockChan <- true
	}()

	if !(<-blockChan) {
		t.Log("failed!")
		t.Fail()
	}
}

func TestRaftRecoverable(t *testing.T) {
	t.Log("leader1 gets a request while all other nodes are crashed. the crashed nodes recover.")
	//Setup
	cfgPath := "./config_files/3nodes.txt"
	test := InitTest(cfgPath, "8080")
	defer EndTest(test)

	// TEST
	leaderIdx := 0
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	// crash all but the leader
	for idx, server := range test.Clients {
		if idx != leaderIdx {
			server.Crash(test.Context, &emptypb.Empty{})
		}
	}

	filemeta1 := &surfstore.FileMetaData{
		Filename:      "testFile1",
		Version:       1,
		BlockHashList: nil,
	}

	reqComplete := make(chan bool)
	go func() {
		_, _ = test.Clients[leaderIdx].UpdateFile(context.Background(), filemeta1)
		reqComplete <- true
	}()

	go func() {
		time.Sleep(1 * time.Second)
		for idx, server := range test.Clients {
			if idx != leaderIdx {
				server.Restore(test.Context, &emptypb.Empty{})
			}
		}
		time.Sleep(5 * time.Second)
		reqComplete <- false
	}()

	time.Sleep(2 * time.Second)

	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	completed := <-reqComplete
	if !completed {
		t.Fail()
	}

	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	goldenMeta := surfstore.NewMetaStore("")
	goldenMeta.UpdateFile(test.Context, filemeta1)
	goldenLog := make([]*surfstore.UpdateOperation, 0)
	goldenLog = append(goldenLog, &surfstore.UpdateOperation{
		Term:         1,
		FileMetaData: filemeta1,
	})

	for idx, server := range test.Clients {
		state, _ := server.GetInternalState(test.Context, &emptypb.Empty{})
		if !SameLog(goldenLog, state.Log) {
			t.Log("Server ", idx, " does not have the correct log")
			t.Fail()
		}
		if !SameMeta(goldenMeta.FileMetaMap, state.MetaMap.FileInfoMap) {
			t.Log("Server ", idx, " does not have the correct state machine.")
			t.Fail()
		}
	}
}

func TestRaftUpdateTwice(t *testing.T) {
	t.Log("leader1 gets a request. leader1 gets another request.")
	cfgPath := "./config_files/3nodes.txt"
	test := InitTest(cfgPath, "8080")
	defer EndTest(test)

	leaderIdx := 0
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	filemeta1 := &surfstore.FileMetaData{
		Filename:      "testFile1",
		Version:       1,
		BlockHashList: nil,
	}
	goldenMeta := surfstore.NewMetaStore("")
	goldenMeta.UpdateFile(test.Context, filemeta1)

	version, err := test.Clients[leaderIdx].UpdateFile(context.Background(), filemeta1)
	// since we call UpdateFile again before checking anything we may not need this here.
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	if err != nil {
		t.Fail()
	}
	if version.Version != 1 {
		t.Fail()
	}

	filemeta1.Version = 2

	version, err = test.Clients[leaderIdx].UpdateFile(context.Background(), filemeta1)
	goldenMeta.UpdateFile(test.Context, filemeta1)

	if err != nil {
		t.Fail()
	}
	if version.Version != 2 {
		t.Fail()
	}

	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	goldenLog := make([]*surfstore.UpdateOperation, 0)
	goldenLog = append(goldenLog, &surfstore.UpdateOperation{
		Term: 1,
		FileMetaData: &surfstore.FileMetaData{
			Filename:      "testFile1",
			Version:       1,
			BlockHashList: nil,
		},
	})
	goldenLog = append(goldenLog, &surfstore.UpdateOperation{
		Term:         1,
		FileMetaData: filemeta1,
	})

	for idx, server := range test.Clients {
		state, _ := server.GetInternalState(test.Context, &emptypb.Empty{})
		if !SameLog(goldenLog, state.Log) {
			t.Log("Server ", idx, " does not have the correct log")
			t.Fail()
		}
		if !SameMeta(goldenMeta.FileMetaMap, state.MetaMap.FileInfoMap) {
			t.Log("Server ", idx, " does not have the correct state machine.")
			t.Fail()
		}
	}
}

func TestRaftLogsCorrectlyOverwritten(t *testing.T) {
	t.Log("leader1 gets several requests while all other nodes are crashed. leader1 crashes. all other nodes are restored. leader2 gets a request. leader1 is restored.")
	//Setup
	cfgPath := "./config_files/3nodes.txt"
	test := InitTest(cfgPath, "8080")
	defer EndTest(test)

	// TEST
	// A is leader
	leaderIdx := 0
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	// B and C are crashed
	for idx, server := range test.Clients {
		if idx != leaderIdx {
			server.Crash(test.Context, &emptypb.Empty{})
		}
	}

	filemeta1 := &surfstore.FileMetaData{
		Filename:      "testFile1",
		Version:       1,
		BlockHashList: nil,
	}

	filemeta2 := &surfstore.FileMetaData{
		Filename:      "testFile2",
		Version:       1,
		BlockHashList: nil,
	}

	filemeta3 := &surfstore.FileMetaData{
		Filename:      "testFile3",
		Version:       1,
		BlockHashList: nil,
	}

	// A gets several requests which it is not able to commit
	errChan := make(chan error, 2)
	go func(b *surfstore.FileMetaData) {
		_, err := test.Clients[leaderIdx].UpdateFile(context.Background(), b)
		errChan <- err
	}(filemeta1)

	go func(b *surfstore.FileMetaData) {
		_, err := test.Clients[leaderIdx].UpdateFile(context.Background(), b)
		errChan <- err
	}(filemeta2)

	// tiny wait in case those goroutines take a little to get scheduled.
	time.Sleep(100 * time.Millisecond)
	state, _ := test.Clients[0].GetInternalState(test.Context, &emptypb.Empty{})
	meta := state.MetaMap.FileInfoMap
	if !SameMeta(meta, make(map[string]*surfstore.FileMetaData)) {
		t.Fail()
	}
	log := state.Log
	// It's possible the entries could be in either order, so just check the length
	if len(log) != 2 {
		t.Log("Leader should have 2 log entries")
		t.Fail()
	}

	// A crashes
	test.Clients[leaderIdx].Crash(test.Context, &emptypb.Empty{})

	// B and C come back
	for idx, server := range test.Clients {
		if idx != leaderIdx {
			server.Restore(test.Context, &emptypb.Empty{})
		}
	}

	// B is leader
	leaderIdx = 1
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})
	// shouldn't really need this here
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	// B gets a request and commits to B and C
	_, err := test.Clients[leaderIdx].UpdateFile(context.Background(), filemeta3)
	if err != nil {
		t.Fail()
	}

	// A comes back
	test.Clients[0].Restore(test.Context, &emptypb.Empty{})

	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	goldenMeta := surfstore.NewMetaStore("")
	goldenMeta.UpdateFile(test.Context, filemeta3)
	goldenLog := make([]*surfstore.UpdateOperation, 0)
	goldenLog = append(goldenLog, &surfstore.UpdateOperation{
		Term:         2,
		FileMetaData: filemeta3,
	})

	// A's logs are correctly overwritten
	for idx, server := range test.Clients {
		state, _ := server.GetInternalState(test.Context, &emptypb.Empty{})
		if !SameLog(goldenLog, state.Log) {
			t.Log("Server ", idx, " does not have the correct log")
			t.Fail()
		}
		if !SameMeta(goldenMeta.FileMetaMap, state.MetaMap.FileInfoMap) {
			t.Log("Server ", idx, " does not have the correct state machine.")
			t.Fail()
		}
	}
}

func TestRaftLogsConsistent(t *testing.T) {
	t.Log("leader1 gets a request while a minority of the cluster is down. leader1 crashes. the other crashed nodes are restored. leader2 gets a request. leader1 is restored.")
	// A B C
	cfgPath := "./config_files/3nodes.txt"
	test := InitTest(cfgPath, "8080")
	defer EndTest(test)

	filemeta1 := &surfstore.FileMetaData{
		Filename:      "testFile1",
		Version:       1,
		BlockHashList: nil,
	}

	filemeta2 := &surfstore.FileMetaData{
		Filename:      "testFile2",
		Version:       1,
		BlockHashList: nil,
	}

	A := 0
	B := 1
	C := 2

	// A is leader
	leaderIdx := A
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	// C is crashed
	test.Clients[C].Crash(test.Context, &emptypb.Empty{})

	// A gets an entry and commits to A and B
	version, err := test.Clients[leaderIdx].UpdateFile(test.Context, filemeta1)
	if version.Version != 1 || err != nil {
		t.Log("Fault1")
		t.Fail()
	}

	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})
	// check A and B logs?

	// A crashes
	test.Clients[A].Crash(test.Context, &emptypb.Empty{})

	// B becomes leader
	leaderIdx = B
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	// C is restored
	test.Clients[C].Restore(test.Context, &emptypb.Empty{})

	// B get an entry and commits to that and the first entry to C
	version, err = test.Clients[B].UpdateFile(test.Context, filemeta2)
	if version.Version != 1 || err != nil {
		t.Log("Fault2")
		t.Fail()
	}

	// A is restored
	test.Clients[A].Restore(test.Context, &emptypb.Empty{})

	// B commits the second entry to A
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	goldenMeta := surfstore.NewMetaStore("")
	goldenMeta.UpdateFile(test.Context, filemeta1)
	goldenMeta.UpdateFile(test.Context, filemeta2)
	goldenLog := make([]*surfstore.UpdateOperation, 0)
	goldenLog = append(goldenLog, &surfstore.UpdateOperation{
		Term:         1,
		FileMetaData: filemeta1,
	})
	goldenLog = append(goldenLog, &surfstore.UpdateOperation{
		Term:         2,
		FileMetaData: filemeta2,
	})

	for idx, server := range test.Clients {
		state, _ := server.GetInternalState(test.Context, &emptypb.Empty{})
		if !SameLog(goldenLog, state.Log) {
			t.Log("Server ", idx, " does not have the correct log")
			t.Fail()
		}
		if !SameMeta(goldenMeta.FileMetaMap, state.MetaMap.FileInfoMap) {
			t.Log("Server ", idx, " does not have the correct state machine.")
			t.Fail()
		}
	}
}

func TestRaftNewLeaderPushesUpdates(t *testing.T) {
	t.Log("leader1 gets a request while the majority of the cluster is down. leader1 crashes. the other nodes come back.")
	// A B C
	cfgPath := "./config_files/5nodes.txt"
	test := InitTest(cfgPath, "8080")
	defer EndTest(test)

	filemeta1 := &surfstore.FileMetaData{
		Filename:      "testFile1",
		Version:       1,
		BlockHashList: nil,
	}

	A := 0
	B := 1
	C := 2
	//D := 3
	E := 4

	// A is leader
	leaderIdx := A
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	// C D E are crashed
	for i := C; i <= E; i++ {
		test.Clients[i].Crash(test.Context, &emptypb.Empty{})
	}

	// A gets an entry and pushes to A and B
	go func() {
		// This should block though and not actually get committed
		test.Clients[A].UpdateFile(test.Context, filemeta1)
	}()

	// this should be plenty of time for A to push to B
	time.Sleep(100 * time.Millisecond)

	// check A and B logs?

	// A crashes
	test.Clients[A].Crash(test.Context, &emptypb.Empty{})

	// C D E come back
	for i := C; i <= E; i++ {
		test.Clients[i].Restore(test.Context, &emptypb.Empty{})
	}

	// B becomes leader
	leaderIdx = B
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	// B commits the second entry to A
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	// Leaders should not commit entries that were created in previous terms.
	goldenMeta := surfstore.NewMetaStore("")
	// Every node should have the entry in it's log though
	goldenLog := make([]*surfstore.UpdateOperation, 0)
	goldenLog = append(goldenLog, &surfstore.UpdateOperation{
		Term:         1,
		FileMetaData: filemeta1,
	})

	for idx, server := range test.Clients {
		state, _ := server.GetInternalState(test.Context, &emptypb.Empty{})
		if !SameLog(goldenLog, state.Log) {
			t.Log("Server ", idx, " does not have the correct log")
			t.Fail()
		}
		if !SameMeta(goldenMeta.FileMetaMap, state.MetaMap.FileInfoMap) {
			t.Log("Server ", idx, " does not have the correct state machine.")
			t.Fail()
		}
	}
}
