package SurfTest

import (
	"os"
	"testing"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
	//	"time"
)

// A syncs with a file. B syncs. B gets the file.
func TestSyncSmallTwoClient(t *testing.T) {
	t.Logf("client1 syncs with a file. client2 syncs.")

	cfgPath := "./config_files/3nodes.txt"
	test := InitTest(cfgPath, "8080")
	defer EndTest(test)
	test.Clients[0].SetLeader(test.Context, &emptypb.Empty{})

	worker1 := InitDirectoryWorker("test0", SRC_PATH)
	worker2 := InitDirectoryWorker("test1", SRC_PATH)
	defer worker1.CleanUp()
	defer worker2.CleanUp()

	//client1 adds file
	newFileName := "image.jpg"
	err := worker1.AddFile(newFileName)
	if err != nil {
		t.FailNow()
	}

	//client1 syncs
	err = SyncClient("localhost:8080", "test0", BLOCK_SIZE)
	if err != nil {
		t.Fatalf("Sync failed")
	}

	//client2 syncs
	err = SyncClient("localhost:8080", "test1", BLOCK_SIZE)
	if err != nil {
		t.Fatalf("Sync failed")
	}

	workingDir, _ := os.Getwd()

	_, err = os.Stat(workingDir + "/test1/" + META_FILENAME)
	if err != nil {
		t.Fatalf("Could not find meta file")
	}

	fileMeta2, err := LoadMetaFromMetaFile(workingDir + "/test1/")
	if err != nil {
		t.Fatalf("Could not load meta file")
	}

	newFileMeta, ok := fileMeta2[newFileName]
	if !ok {
		t.Fatalf("New file metadata not found in client2's index.txt")
	}
	if newFileMeta.Version != 1 {
		t.Fatalf("Wrong version found in client2's index.txt")
	}

	c, e := SameFile(workingDir+"/test1/image.jpg", workingDir+"/test0/image.jpg")
	if e != nil {
		t.Fatalf("Could not read files in client base dirs.")
	}
	if !c {
		t.Fatalf("Files do not match at client1 and client2.")
	}
}

// A syncs with a large file. B syncs. B gets the file.
// func TestSyncLargeTwoClient(t *testing.T) {
// 	t.Logf("client1 syncs with a large file. client2 syncs.")

// 	cfgPath := "./config_files/3nodes.txt"
// 	test := InitTest(cfgPath, "8080")
// 	defer EndTest(test)
// 	test.Clients[0].SetLeader(test.Context, &emptypb.Empty{})

// 	worker1 := InitDirectoryWorker("test0", SRC_PATH)
// 	worker2 := InitDirectoryWorker("test1", SRC_PATH)
// 	defer worker1.CleanUp()
// 	defer worker2.CleanUp()

// 	//client1 adds file
// 	newFileName := "large_file.txt"
// 	err := worker1.AddFile(newFileName)
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

// 	workingDir, _ := os.Getwd()

// 	_, err = os.Stat(workingDir + "/test1/" + META_FILENAME)
// 	if err != nil {
// 		t.Fatalf("Could not find meta file")
// 	}

// 	fileMeta2, err := LoadMetaFromMetaFile(workingDir + "/test1/")
// 	if err != nil {
// 		t.Fatalf("Could not load meta file")
// 	}

// 	newFileMeta, ok := fileMeta2[newFileName]
// 	if !ok {
// 		t.Fatalf("New file metadata not found in client2's index.txt")
// 	}
// 	if newFileMeta.Version != 1 {
// 		t.Fatalf("Wrong version found in client2's index.txt")
// 	}

// 	c, e := SameFile(workingDir+"/test1/large_file.txt", workingDir+"/test0/large_file.txt")
// 	if e != nil {
// 		t.Fatalf("Could not read files in client base dirs.")
// 	}
// 	if !c {
// 		t.Fatalf("Files do not match at client1 and client2.")
// 	}
// }

// // A syncs with a file. B syncs with different file. A syncs again. Both should have both files.
// func TestSyncTwoClientsDisjointFiles(t *testing.T) {
// 	t.Logf("client1 syncs with file1. client2 syncs with file2. client1 syncs again.")
// 	cfgPath := "./config_files/3nodes.txt"
// 	test := InitTest(cfgPath, "8080")
// 	defer EndTest(test)
// 	test.Clients[0].SetLeader(test.Context, &emptypb.Empty{})

// 	worker1 := InitDirectoryWorker("test0", SRC_PATH)
// 	worker2 := InitDirectoryWorker("test1", SRC_PATH)
// 	defer worker1.CleanUp()
// 	defer worker2.CleanUp()

// 	//clients add different files
// 	file1 := "multi_file1.txt"
// 	file2 := "multi_file2.txt"
// 	err := worker1.AddFile(file1)
// 	if err != nil {
// 		t.FailNow()
// 	}
// 	err = worker2.AddFile(file2)
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
// 	if len(fileMeta1) != 2 {
// 		t.Fatalf("wrong number of entries in client1 meta file")
// 	}
// 	if fileMeta1[file2].Version != 1 {
// 		t.Fatalf("Wrong version for file2 in client1 metadata.")
// 	}

// 	c, e := SameFile(workingDir+"/test0/multi_file2.txt", workingDir+"/test1/multi_file2.txt")
// 	if e != nil {
// 		t.Fatalf("Could not read files in client base dirs.")
// 	}
// 	if !c {
// 		t.Fatalf("file2 does not match at client1 and client2.")
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
// 	if len(fileMeta2) != 2 {
// 		t.Fatalf("Wrong number of entries in client2 meta file")
// 	}
// 	if fileMeta2[file1].Version != 1 {
// 		t.Fatalf("Wrong version for file1 in client2 metadata.")
// 	}

// 	c, e = SameFile(workingDir+"/test0/multi_file1.txt", workingDir+"/test1/multi_file1.txt")
// 	if e != nil {
// 		t.Fatalf("Could not read files in client base dirs.")
// 	}
// 	if !c {
// 		t.Fatalf("file1 does not match at client1 and client2.")
// 	}
// }

// // A syncs, modifies, syncs. B syncs. B should get all updates and correct version.
// func TestSyncTwoClientsOneUpdates(t *testing.T) {
// 	t.Logf("client1 syncs with a file, modifies, syncs again. client2 syncs.")

// 	cfgPath := "./config_files/3nodes.txt"
// 	test := InitTest(cfgPath, "8080")
// 	defer EndTest(test)
// 	test.Clients[0].SetLeader(test.Context, &emptypb.Empty{})

// 	worker1 := InitDirectoryWorker("test0", SRC_PATH)
// 	worker2 := InitDirectoryWorker("test1", SRC_PATH)
// 	defer worker1.CleanUp()
// 	defer worker2.CleanUp()

// 	//client1 adds file
// 	newFileName := "multi_file1.txt"
// 	err := worker1.AddFile(newFileName)
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	//client1 syncs
// 	err = SyncClient("localhost:8080", "test0", BLOCK_SIZE)
// 	if err != nil {
// 		t.Fatalf("Sync failed")
// 	}

// 	//client1 updates
// 	err = worker1.UpdateFile(newFileName, "update string")
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

// 	workingDir, _ := os.Getwd()

// 	//check client2
// 	_, err = os.Stat(workingDir + "/test1/" + META_FILENAME)
// 	if err != nil {
// 		t.Fatalf("Could not find meta file")
// 	}

// 	fileMeta2, err := LoadMetaFromMetaFile(workingDir + "/test1/")
// 	if err != nil {
// 		t.Fatalf("Could not load meta file for client2")
// 	}
// 	if len(fileMeta2) != 1 {
// 		t.Fatalf("Wrong number of entries in client2 meta file")
// 	}
// 	if fileMeta2[newFileName].Version != 2 {
// 		t.Fatalf("Wrong version for file in client2 metadata.")
// 	}

// 	c, e := SameFile(workingDir+"/test0/multi_file1.txt", workingDir+"/test1/multi_file1.txt")
// 	if e != nil {
// 		t.Fatalf("Could not read files in client base dirs.")
// 	}
// 	if !c {
// 		t.Fatalf("file1 does not match at client1 and client2.")
// 	}
// }
