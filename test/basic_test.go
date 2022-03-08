package SurfTest

// import (
// 	emptypb "google.golang.org/protobuf/types/known/emptypb"
// 	"os"
// 	"testing"
// 	"time"
// )

// // index.txt should be generated with empty content.
// func TestIndexFile(t *testing.T) {
// 	t.Logf("client1 syncs with an empty base directory.")
// 	cfgPath := "./config_files/3nodes.txt"
// 	test := InitTest(cfgPath, "8080")
// 	defer EndTest(test)
// 	test.Clients[0].SetLeader(test.Context, &emptypb.Empty{})
// 	test.Clients[0].SendHeartbeat(test.Context, &emptypb.Empty{})

// 	worker1 := InitDirectoryWorker("test0", SRC_PATH)
// 	defer worker1.CleanUp()

// 	err := SyncClient("localhost:8080", "test0", BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	workingDir, _ := os.Getwd()

// 	metaFileInfo, err := os.Stat(workingDir + "/test0/" + META_FILENAME)
// 	if err != nil || metaFileInfo.Size() != 0 {
// 		t.Fatalf("index.txt not found.")
// 	}
// }

// // A syncs with a file, index.txt should reflect it.
// func TestSmallSyncOneClient(t *testing.T) {
// 	t.Logf("client1 syncs with a small file.")
// 	cfgPath := "./config_files/3nodes.txt"
// 	test := InitTest(cfgPath, "8080")
// 	defer EndTest(test)
// 	test.Clients[0].SetLeader(test.Context, &emptypb.Empty{})
// 	test.Clients[0].SendHeartbeat(test.Context, &emptypb.Empty{})

// 	worker1 := InitDirectoryWorker("test0", SRC_PATH)
// 	defer worker1.CleanUp()

// 	newFileName := "image.jpg"
// 	err := worker1.AddFile(newFileName)
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	err = SyncClient("localhost:8080", "test0", BLOCK_SIZE)
// 	if err != nil {
// 		t.Log(err)
// 		t.Fatalf("Sync failed")
// 	}

// 	workingDir, _ := os.Getwd()

// 	_, err = os.Stat(workingDir + "/test0/" + META_FILENAME)
// 	if err != nil {
// 		t.Fatalf("Could not find meta file")
// 	}

// 	fileMeta, err := LoadMetaFromMetaFile(workingDir + "/test0/")
// 	if err != nil {
// 		t.Fatalf("Could not load meta file")
// 	}

// 	newFileMeta, ok := fileMeta[newFileName]
// 	if !ok {
// 		t.Fatalf("New file metadata not found in index.txt")
// 	}
// 	if newFileMeta.Version != 1 {
// 		t.Fatalf("Wrong version found in index.txt")
// 	}
// }

// // A syncs twice with an update, check version.
// func TestSyncTwiceOneClient(t *testing.T) {
// 	t.Logf("client1 adds a file. client1 syncs. client1 modifies file. client1 syncs.")
// 	cfgPath := "./config_files/3nodes.txt"
// 	test := InitTest(cfgPath, "8080")
// 	defer EndTest(test)
// 	test.Clients[0].SetLeader(test.Context, &emptypb.Empty{})
// 	test.Clients[0].SendHeartbeat(test.Context, &emptypb.Empty{})

// 	worker1 := InitDirectoryWorker("test0", SRC_PATH)
// 	defer worker1.CleanUp()

// 	newFileName := "calendar.txt"
// 	err := worker1.AddFile(newFileName)
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	//client1 sync
// 	err = SyncClient("localhost:8080", "test0", BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	//update file
// 	err = worker1.UpdateFile(newFileName, "update string")
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	//client1 sync again
// 	err = SyncClient("localhost:8080", "test0", BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	workingDir, _ := os.Getwd()

// 	_, err = os.Stat(workingDir + "/test0/" + META_FILENAME)
// 	if err != nil {
// 		t.Fatalf("Could not find meta file")
// 	}

// 	fileMeta, err := LoadMetaFromMetaFile(workingDir + "/test0/")
// 	if err != nil {
// 		t.Fatalf("Could not read meta file")
// 	}
// 	t.Log(fileMeta)

// 	newFileMeta, ok := fileMeta[newFileName]
// 	if !ok {
// 		t.Fatalf("New file metadata not found in index.txt")
// 	}
// 	if newFileMeta.Version != 2 {
// 		t.Fatalf("Wrong version found in index.txt")
// 	}
// }

// // A syncs twice with delete, check hashlist.
// func TestSyncTwiceWithDeleteOneClient(t *testing.T) {
// 	t.Logf("client1 syncs with a file. client1 deletes file. client1 syncs.")
// 	cfgPath := "./config_files/3nodes.txt"
// 	test := InitTest(cfgPath, "8080")
// 	defer EndTest(test)
// 	test.Clients[0].SetLeader(test.Context, &emptypb.Empty{})
// 	test.Clients[0].SendHeartbeat(test.Context, &emptypb.Empty{})

// 	worker1 := InitDirectoryWorker("test0", SRC_PATH)
// 	defer worker1.CleanUp()

// 	newFileName := "calendar.txt"
// 	err := worker1.AddFile(newFileName)
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	//client1 sync
// 	err = SyncClient("localhost:8080", "test0", BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	//update file
// 	err = worker1.DeleteFile(newFileName)
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	//client1 sync again
// 	err = SyncClient("localhost:8080", "test0", BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	time.Sleep(3 * time.Second)

// 	workingDir, _ := os.Getwd()

// 	_, err = os.Stat(workingDir + "/test0/" + META_FILENAME)
// 	if err != nil {
// 		t.Fatalf("Could not find meta file")
// 	}

// 	fileMeta, err := LoadMetaFromMetaFile(workingDir + "/test0/")
// 	if err != nil {
// 		t.Fatalf("Could not read meta file")
// 	}

// 	t.Log(fileMeta)

// 	newFileMeta, ok := fileMeta[newFileName]
// 	if !ok {
// 		t.Fatalf("New file metadata not found in index.txt")
// 	}
// 	if newFileMeta.Version != 2 {
// 		t.Fatalf("Wrong version found in index.txt")
// 	}
// 	if !(len(newFileMeta.BlockHashList) == 1 && newFileMeta.BlockHashList[0] == TOMBSTONE_HASH) {
// 		t.Fatalf("Hashlist should be Tombstone")
// 	}
// }
