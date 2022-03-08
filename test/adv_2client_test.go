package SurfTest

// import (
// 	emptypb "google.golang.org/protobuf/types/known/emptypb"
// 	"os"
// 	"testing"
// 	//	"time"
// )

// // A creates and syncs with a file. B creates and syncs with same file. A syncs again.
// func TestSyncTwoClientsSameFile(t *testing.T) {
// 	t.Logf("client1 syncs with file1. client2 syncs with file1 (different content). client1 syncs again.")
// 	cfgPath := "./config_files/3nodes.txt"
// 	test := InitTest(cfgPath, "8080")
// 	defer EndTest(test)
// 	test.Clients[0].SetLeader(test.Context, &emptypb.Empty{})
// 	//	serverCmds := InitSurfServers(1)
// 	//
// 	//	ready := make(chan bool)
// 	//	go StartSurfServers(serverCmds, ready)
// 	//	defer KillSurfServers(serverCmds)
// 	//
// 	//	//server started
// 	//	time.Sleep(3 * time.Second)
// 	//	<-ready

// 	worker1 := InitDirectoryWorker("test0", SRC_PATH)
// 	worker2 := InitDirectoryWorker("test1", SRC_PATH)
// 	defer worker1.CleanUp()
// 	defer worker2.CleanUp()

// 	//clients add different files
// 	file1 := "multi_file1.txt"
// 	file2 := "multi_file1.txt"
// 	err := worker1.AddFile(file1)
// 	if err != nil {
// 		t.FailNow()
// 	}
// 	err = worker2.AddFile(file2)
// 	if err != nil {
// 		t.FailNow()
// 	}
// 	err = worker2.UpdateFile(file2, "update text")
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	//client1 syncs
// 	err = SyncClient("localhost:8080", "test0", BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	//client2 syncs
// 	err = SyncClient("localhost:8080", "test1", BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	//client1 syncs
// 	err = SyncClient("localhost:8080", "test0", BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	workingDir, _ := os.Getwd()

// 	//check client1
// 	_, err = os.Stat(workingDir + "/test0/" + META_FILENAME)
// 	if err != nil {
// 		t.Fatalf("Could not find meta file for client1")
// 	}

// 	fileMeta1, err := LoadMetaFromMetaFile(workingDir + "/test0/")
// 	if err != nil {
// 		t.Fatalf("Could not load meta file for client1")
// 	}
// 	if len(fileMeta1) != 1 {
// 		t.Fatalf("Wrong number of entries in client1 meta file")
// 	}
// 	if fileMeta1[file1].Version != 1 {
// 		t.Fatalf("Wrong version for file1 in client1 metadata.")
// 	}

// 	c, e := SameFile(workingDir+"/test0/multi_file1.txt", SRC_PATH+"/multi_file1.txt")
// 	if e != nil {
// 		t.Fatalf("Could not read files in client base dirs.")
// 	}
// 	if !c {
// 		t.Fatalf("file1 should not change at client1")
// 	}

// 	//check client2
// 	_, err = os.Stat(workingDir + "/test1/" + META_FILENAME)
// 	if err != nil {
// 		t.Fatalf("Could not find meta file for client2")
// 	}

// 	fileMeta2, err := LoadMetaFromMetaFile(workingDir + "/test1/")
// 	if err != nil {
// 		t.Fatalf("Could not load meta file for client2")
// 	}
// 	if len(fileMeta2) != 1 {
// 		t.Fatalf("Wrong number of entries in client2 meta file")
// 	}
// 	if fileMeta1[file1].Version != 1 {
// 		t.Fatalf("Wrong version for file1 in client2 metadata.")
// 	}

// 	c, e = SameFile(workingDir+"/test1/multi_file1.txt", SRC_PATH+"/multi_file1.txt")
// 	if e != nil {
// 		t.Fatalf("Could not read files in client base dirs.")
// 	}
// 	if !c {
// 		t.Fatalf("wrong file2 contents at client2")
// 	}
// }

// // A creates and syncs. A updates and syncs. B syncs, updates and syncs. A syncs.
// func TestSyncTwoClientsMultipleUpdates(t *testing.T) {
// 	t.Logf("client1 creates and syncs. client1 modifies and syncs. client2 syncs. client2 modifies and syncs. client1 syncs.")
// 	cfgPath := "./config_files/3nodes.txt"
// 	test := InitTest(cfgPath, "8080")
// 	defer EndTest(test)
// 	test.Clients[0].SetLeader(test.Context, &emptypb.Empty{})

// 	//    serverCmds := InitSurfServers(1)
// 	//
// 	//	ready := make(chan bool)
// 	//	go StartSurfServers(serverCmds, ready)
// 	//	defer KillSurfServers(serverCmds)
// 	//
// 	//	//server started
// 	//	time.Sleep(3 * time.Second)
// 	//	<-ready

// 	worker1 := InitDirectoryWorker("test0", SRC_PATH)
// 	worker2 := InitDirectoryWorker("test1", SRC_PATH)
// 	defer worker1.CleanUp()
// 	defer worker2.CleanUp()

// 	//client1 creates
// 	newFileName := "calendar.txt"
// 	err := worker1.AddFile(newFileName)
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	//client1 syncs
// 	err = SyncClient("localhost:8080", "test0", BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	// client1 updates
// 	err = worker1.UpdateFile(newFileName, "\nhahanew")
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	//client1 syncs
// 	err = SyncClient("localhost:8080", "test0", BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	//client2 syncs
// 	err = SyncClient("localhost:8080", "test1", BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	//client2 updates
// 	err = worker2.TruncateFile(newFileName, 3)
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	//client2 syncs
// 	err = SyncClient("localhost:8080", "test1", BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	//client1 syncs
// 	err = SyncClient("localhost:8080", "test0", BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	workingDir, _ := os.Getwd()

// 	//check client1
// 	_, err = os.Stat(workingDir + "/test0/" + META_FILENAME)
// 	if err != nil {
// 		t.Fatalf("Could not find meta file for client1")
// 	}

// 	fileMeta1, err := LoadMetaFromMetaFile(workingDir + "/test0/")
// 	if err != nil {
// 		t.Fatalf("Could not load meta file for client1")
// 	}
// 	if fileMeta1[newFileName].Version != 3 {
// 		t.Fatalf("Wrong version for file in client1 metadata.")
// 	}

// 	c, e := SameFile(workingDir+"/test0/calendar.txt", workingDir+"/test1/calendar.txt")
// 	if e != nil {
// 		t.Fatalf("Could not read files in client base dirs.")
// 	}
// 	if !c {
// 		t.Fatalf("file does not match at client1 and client2")
// 	}
// }

// // A creates and syncs. B syncs, deletes and syncs. A modifies and syncs. B syncs.
// func TestSyncTwoClientsDeleteConflict(t *testing.T) {
// 	t.Logf("client1 creates and syncs. client2 syncs, deletes and syncs. client1 modifies and syncs. client2 syncs.")
// 	cfgPath := "./config_files/3nodes.txt"
// 	test := InitTest(cfgPath, "8080")
// 	defer EndTest(test)
// 	test.Clients[0].SetLeader(test.Context, &emptypb.Empty{})

// 	//    serverCmds := InitSurfServers(1)
// 	//
// 	//	ready := make(chan bool)
// 	//	go StartSurfServers(serverCmds, ready)
// 	//	defer KillSurfServers(serverCmds)
// 	//
// 	//	// Waiting for server to bootup. you might want to extend it based on yor implementation.
// 	//	// A better implementation would use a ready channel to respond back.
// 	//	time.Sleep(3 * time.Second)
// 	//	<-ready

// 	worker1 := InitDirectoryWorker("test0", SRC_PATH)
// 	worker2 := InitDirectoryWorker("test1", SRC_PATH)
// 	defer worker1.CleanUp()
// 	defer worker2.CleanUp()

// 	//client1 creates
// 	err := worker1.AddFile("image.jpg")
// 	if err != nil {
// 		t.FailNow()
// 	}
// 	err = worker1.AddFile("calendar.txt")
// 	if err != nil {
// 		t.FailNow()
// 	}
// 	err = worker1.AddFile("calendar3.txt")
// 	if err != nil {
// 		t.FailNow()
// 	}
// 	err = worker1.AddFile("calendar2.txt")
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	//client1 syncs
// 	err = SyncClient("localhost:8080", "test0", DEFAULT_BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	//client2 syncs
// 	err = SyncClient("localhost:8080", "test1", DEFAULT_BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	//client2 deletes
// 	err = worker2.DeleteFile("calendar.txt")
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	//client2 syncs
// 	err = SyncClient("localhost:8080", "test1", DEFAULT_BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	//client1 updates
// 	err = worker1.UpdateFile("calendar.txt", "update text")
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	//client1 syncs
// 	err = SyncClient("localhost:8080", "test0", DEFAULT_BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	//client2 syncs
// 	err = SyncClient("localhost:8080", "test1", DEFAULT_BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	workingDir, _ := os.Getwd()

// 	//check client1
// 	_, err = os.Stat(workingDir + "/test0/" + META_FILENAME)
// 	if err != nil {
// 		t.Fatalf("Could not find meta file for client1")
// 	}

// 	fileMeta1, err := LoadMetaFromMetaFile(workingDir + "/test0/")
// 	if fileMeta1["calendar.txt"].Version != 2 {
// 		t.Fatalf("Wrong version for deleted file in client1 metadata")
// 	}
// 	if !IsTombHashList(fileMeta1["calendar.txt"].BlockHashList) {
// 		t.Fatalf("Deleted file should have tombstone hashlist in client1 metadata")
// 	}

// 	fileList1 := worker1.ListAllFile()
// 	_, calendarExist := fileList1["calendar.txt"]
// 	if calendarExist {
// 		t.Fatalf("deleted file should not exist in client1 basedir")
// 	}

// 	//check client2
// 	_, err = os.Stat(workingDir + "/test1/" + META_FILENAME)
// 	if err != nil {
// 		t.Fatalf("Could not find meta file for client2")
// 	}

// 	fileMeta2, err := LoadMetaFromMetaFile(workingDir + "/test1/")
// 	if fileMeta2["calendar.txt"].Version != 2 {
// 		t.Fatalf("Wrong version for deleted file in client2 metadata")
// 	}
// 	if !IsTombHashList(fileMeta2["calendar.txt"].BlockHashList) {
// 		t.Fatalf("Deleted file should have tombstone hashlist in client2 metadata")
// 	}

// 	fileList2 := worker2.ListAllFile()
// 	_, calendarExist = fileList2["calendar.txt"]
// 	if calendarExist {
// 		t.Fatalf("deleted file should not exist in client2 basedir")
// 	}
// }

// // A creates and syncs. B syncs, deletes and syncs. A syncs, creates and syncs. B syncs.
// func TestSyncTwoClientsDeleteRecreate(t *testing.T) {
// 	t.Logf("client1 creates and syncs. client2 syncs, deletes and syncs. client1 syncs, re-creates and syncs. client2 syncs.")
// 	cfgPath := "./config_files/3nodes.txt"
// 	test := InitTest(cfgPath, "8080")
// 	defer EndTest(test)
// 	test.Clients[0].SetLeader(test.Context, &emptypb.Empty{})

// 	//    serverCmds := InitSurfServers(1)
// 	//
// 	//	ready := make(chan bool)
// 	//	go StartSurfServers(serverCmds, ready)
// 	//	defer KillSurfServers(serverCmds)
// 	//
// 	//	// Waiting for server to bootup. you might want to extend it based on yor implementation.
// 	//	// A better implementation would use a ready channel to respond back.
// 	//	time.Sleep(3 * time.Second)
// 	//	<-ready

// 	worker1 := InitDirectoryWorker("test0", SRC_PATH)
// 	worker2 := InitDirectoryWorker("test1", SRC_PATH)
// 	defer worker1.CleanUp()
// 	defer worker2.CleanUp()

// 	//client1 creates
// 	err := worker1.AddFile("image.jpg")
// 	if err != nil {
// 		t.FailNow()
// 	}
// 	err = worker1.AddFile("calendar.txt")
// 	if err != nil {
// 		t.FailNow()
// 	}
// 	err = worker1.AddFile("calendar3.txt")
// 	if err != nil {
// 		t.FailNow()
// 	}
// 	err = worker1.AddFile("calendar2.txt")
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	//client1 syncs
// 	err = SyncClient("localhost:8080", "test0", DEFAULT_BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	//client2 syncs
// 	err = SyncClient("localhost:8080", "test1", DEFAULT_BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	//client2 deletes
// 	err = worker2.DeleteFile("calendar.txt")
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	//client2 syncs
// 	err = SyncClient("localhost:8080", "test1", DEFAULT_BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	//client1 syncs
// 	err = SyncClient("localhost:8080", "test0", DEFAULT_BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	//client1 re-creates
// 	err = worker1.AddFile("calendar.txt")
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	//client1 syncs
// 	err = SyncClient("localhost:8080", "test0", DEFAULT_BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	//client2 syncs
// 	err = SyncClient("localhost:8080", "test1", DEFAULT_BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	workingDir, _ := os.Getwd()

// 	//check client1
// 	_, err = os.Stat(workingDir + "/test0/" + META_FILENAME)
// 	if err != nil {
// 		t.Fatalf("Could not find meta file for client1")
// 	}

// 	fileMeta1, err := LoadMetaFromMetaFile(workingDir + "/test0/")
// 	if fileMeta1["calendar.txt"].Version != 3 {
// 		t.Fatalf("Wrong version for re-created file in client1 metadata")
// 	}

// 	fileList1 := worker1.ListAllFile()
// 	_, calendarExist := fileList1["calendar.txt"]
// 	if !calendarExist {
// 		t.Fatalf("file was deleted from client1 basedir")
// 	}

// 	//check client2
// 	_, err = os.Stat(workingDir + "/test1/" + META_FILENAME)
// 	if err != nil {
// 		t.Fatalf("Could not find meta file for client2")
// 	}

// 	fileMeta2, err := LoadMetaFromMetaFile(workingDir + "/test1/")
// 	if fileMeta2["calendar.txt"].Version != 3 {
// 		t.Fatalf("Wrong version for deleted file in client2 metadata")
// 	}

// 	fileList2 := worker2.ListAllFile()
// 	_, calendarExist = fileList2["calendar.txt"]
// 	if !calendarExist {
// 		t.Fatalf("file was deleted in client2 basedir")
// 	}

// 	c, e := SameFile(workingDir+"/test0/calendar.txt", workingDir+"/test1/calendar.txt")
// 	if e != nil {
// 		t.Fatalf("Could not read files in client base dirs.")
// 	}
// 	if !c {
// 		t.Fatalf("re-created file does not match at client1 and client2.")
// 	}
// }
