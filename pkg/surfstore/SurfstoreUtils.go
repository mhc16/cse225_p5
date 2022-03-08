package surfstore

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// Implement the logic for a client syncing with the server here.
func ClientSync(client RPCClient) {
	baseDir, readErr := ioutil.ReadDir(client.BaseDir)
	if readErr != nil {
		log.Println("Read client base directory error: ", readErr)
	}
	dirFileMap := make(map[string]os.FileInfo)
	for _, fileInfo := range baseDir {
		dirFileMap[fileInfo.Name()] = fileInfo
	}

	// Create the index.txt if not existed
	indexFilePath := client.BaseDir + "/index.txt"
	if _, err := os.Stat(indexFilePath); os.IsNotExist(err) {
		file, _ := os.Create(indexFilePath)
		defer file.Close()
	}

	fileMetaMap, indexMap, _ := LoadMetaFromMetaFile(client.BaseDir)

	// fmt.Println("After LoadMetaFromMetaFile, the LoadMetaFromMetaFile is ", fileMetaMap)

	indexFile, _ := ioutil.ReadFile(indexFilePath)
	indexLines := strings.Split(string(indexFile), "\n")

	clientFileInfoMap := LocalSync(client, fileMetaMap, &indexMap, dirFileMap, &indexLines)
	serverFileInfoMap := make(map[string]*FileMetaData)
	err := client.GetFileInfoMap(&serverFileInfoMap)
	if err != nil {
		log.Println("Cannot get fileInfoMap from the server:", err)
	}

	// Upload updated file to server
	for fileName, fileInfo := range clientFileInfoMap {
		if _, ok := serverFileInfoMap[fileName]; ok {
			clientFileMetaData := fileInfo.FileMetaData
			serverFileMetaData := serverFileInfoMap[fileName]

			// fmt.Println("clientFileMetaData.Version:", clientFileMetaData.Version)
			// fmt.Println("serverFileMetaData.Version:", serverFileMetaData.Version)
			// fmt.Println("fileInfo.Status:", fileInfo.Status)

			if clientFileMetaData.Version == serverFileMetaData.Version && fileInfo.Status == Unchanged {
				continue
			} else if (clientFileMetaData.Version > serverFileMetaData.Version) ||
				(clientFileMetaData.Version == serverFileMetaData.Version && fileInfo.Status == Modified) {
				updateServerFile(client, clientFileMetaData, indexMap, &indexLines, fileInfo.Status)
			} else {
				updateClientFile(client, serverFileMetaData, indexMap, &indexLines)
			}
		} else {
			upload(client, fileInfo.FileMetaData, indexMap, &indexLines)
		}
	}

	// Download new files from server
	for fileName, serverFileMetaData := range serverFileInfoMap {
		if _, ok := clientFileInfoMap[fileName]; !ok {
			if _, okay := indexMap[fileName]; okay {
				deletedFileMetaData := NewFileMetaDataFromConfig(indexLines[indexMap[fileName]])
				if deletedFileMetaData.Version > serverFileMetaData.Version {
					updateServerFile(client, deletedFileMetaData, indexMap, &indexLines, Deleted)
				} else {
					updateClientFile(client, serverFileMetaData, indexMap, &indexLines)
				}
			} else {
				line, err := download(client, fileName, serverFileMetaData)
				if err != nil {
					log.Println("Fail to download file from server: ", err)
				}
				indexLines = append((indexLines), line)
			}
		}
	}

	newIndexFile := ""
	for _, indexLine := range indexLines {
		if indexLine == "" {
			continue
		}
		newIndexFile += indexLine + "\n"
	}

	err = ioutil.WriteFile(indexFilePath, []byte(newIndexFile), 0755)
	if err != nil {
		log.Println("Fail to update index.txt: ", err)
	}

	fmt.Println("---end of ClientSync---")
	// fmt.Println("clientFileInfoMap:", clientFileInfoMap)
	// fmt.Println("newIndexFile:", newIndexFile)
}

func LocalSync(client RPCClient, fileMetaMap map[string]*FileMetaData, indexMap *map[string]int, dirFileMap map[string]os.FileInfo, indexLines *[]string) (clientFileInfoMap map[string]FileInfo) {
	// to check files deleted
	hasBeenDeleted(fileMetaMap, *indexMap, dirFileMap, indexLines)
	clientFileInfoMap = make(map[string]FileInfo)

	fmt.Println("-------begin to local sync--------")
	// fmt.Println("fileMetaMap:", fileMetaMap)
	// fmt.Println("dirFileMap", dirFileMap)

	for filename, f := range dirFileMap {
		if filename == "index.txt" {
			continue
		}

		fmt.Println("filename: ", filename)

		file, err := os.Open(client.BaseDir + "/" + filename)
		if err != nil {
			log.Println("Open file Error: ", err)
		}
		fileSize := f.Size()
		numberBlock := int(math.Ceil(float64(fileSize) / float64(client.BlockSize)))

		fileInfo := &FileInfo{FileMetaData: &FileMetaData{}} // record the current file info
		if fileMetaData, ok := fileMetaMap[filename]; ok {
			// the file exists in both index.txt and base directory
			fmt.Println("exists in both index.txt and base directory")
			isChanged, hashList := getHashList(file, fileMetaData, numberBlock, client.BlockSize)
			fileInfo.FileMetaData.Filename = filename
			fileInfo.FileMetaData.Version = fileMetaData.Version
			hashStr := ""
			for _, hash := range hashList {
				fileInfo.FileMetaData.BlockHashList = append(fileInfo.FileMetaData.BlockHashList, hash)
				hashStr += hash + " "
				// if i != len(hashList)-1 {
				// 	hashStr += " "
				// }
			}

			fmt.Println("!!!this file is changed:", isChanged)
			if isChanged {
				fileInfo.Status = Modified
				(*indexLines)[(*indexMap)[filename]] = filename + "," + strconv.Itoa(int(fileMetaData.Version)) + "," + hashStr + " "
			} else {
				fileInfo.Status = Unchanged
			}

		} else {
			// the file is not recored in index.txt
			fmt.Println("not recored in index.txt")
			fileMetaData := &FileMetaData{}
			_, hashList := getHashList(file, fileMetaData, numberBlock, client.BlockSize)
			fileInfo.FileMetaData.Filename = filename
			fileInfo.FileMetaData.Version = 1
			hashStr := ""
			for _, hash := range hashList {
				fileInfo.FileMetaData.BlockHashList = append(fileInfo.FileMetaData.BlockHashList, hash)
				hashStr += hash + " "
				// if i != len(hashList)-1 {
				// 	hashStr += " "
				// }
			}
			fileInfo.FileMetaData.BlockHashList = hashList
			fileInfo.Status = New

			*indexLines = append((*indexLines), filename+","+strconv.Itoa(int(fileInfo.FileMetaData.Version))+","+hashStr)
			(*indexMap)[filename] = len(*indexLines) - 1
		}

		clientFileInfoMap[filename] = *fileInfo
	}

	fmt.Println("-------the end to local sync--------")
	// fmt.Println("clientFileInfoMap:", clientFileInfoMap)
	return clientFileInfoMap
}

func hasBeenDeleted(fileMetaMap map[string]*FileMetaData, indexMap map[string]int, dirFileMap map[string]os.FileInfo, indexLines *[]string) {
	for filename, fileMetaData := range fileMetaMap {
		// there exists a file recored in index.txt but not in base directory
		if _, ok := dirFileMap[filename]; !ok {
			index := indexMap[filename]
			if len(fileMetaData.BlockHashList) != 1 || fileMetaData.BlockHashList[0] != "0" {
				fileMetaData.Version += 1
			}
			fileMetaData.BlockHashList = []string{"0"}
			(*indexLines)[index] = FileMetaDataToString(fileMetaData)
		}
	}
}

func getHashList(file *os.File, fileMetaData *FileMetaData, numberBlock int, blockSize int) (bool, []string) {
	hashList := make([]string, numberBlock)
	var isChanged bool = false
	for i := 0; i < numberBlock; i++ {
		buf := make([]byte, blockSize)
		n, e := file.Read(buf)
		if e != nil {
			log.Println("read error when getting hashList: ", e)
		}
		buf = buf[:n]

		key := GetBlockHashString(buf)
		hashList[i] = key
		if i >= len(fileMetaData.BlockHashList) || key != fileMetaData.BlockHashList[i] {
			isChanged = true
		}
	}
	if numberBlock != len(fileMetaData.BlockHashList) {
		isChanged = true
	}
	return isChanged, hashList
}

func updateServerFile(client RPCClient, clientFileMetaData *FileMetaData, indexMap map[string]int, indexLines *[]string, status State) {
	if status == Modified {
		clientFileMetaData.Version += 1
		index := indexMap[clientFileMetaData.Filename]
		line := (*indexLines)[index]
		(*indexLines)[index] = line[:strings.Index(line, ",")] + "," + strconv.Itoa(int(clientFileMetaData.Version)) + "," + line[strings.LastIndex(line, ",")+1:]
	}

	err := upload(client, clientFileMetaData, indexMap, indexLines)
	if err != nil {
		log.Println("File to upload file: ", err)
	}
}

func updateClientFile(client RPCClient, serverFileMetaData *FileMetaData, indexMap map[string]int, indexLines *[]string) {

	// fmt.Println("begin to uodate client file----------------")
	// fmt.Println("filename", serverFileMetaData.Filename)

	line, err := download(client, serverFileMetaData.Filename, serverFileMetaData)

	if err != nil {
		log.Println("Download file from server failed: ", err)
	}

	index := indexMap[serverFileMetaData.Filename]
	(*indexLines)[index] = line
}

func upload(client RPCClient, fileMetaData *FileMetaData, indexMap map[string]int, indexLines *[]string) error {

	// fmt.Println("----------------begin upload!!!!!!!!!!!!!")
	// fmt.Println("fileMetaData.filename:", fileMetaData.Filename)
	// fmt.Println("fileMetaData.version:", fileMetaData.Version)

	filePath := client.BaseDir + "/" + fileMetaData.Filename
	if _, e := os.Stat(filePath); os.IsNotExist(e) {
		err := client.UpdateFile(fileMetaData, &fileMetaData.Version)
		if err != nil {
			log.Println("Fail to update file: ", err)
		}
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Open file Error: ", err)
	}

	defer file.Close()

	f, _ := os.Stat(filePath)
	numberBlock := int(math.Ceil(float64(f.Size()) / float64(client.BlockSize)))

	for i := 0; i < numberBlock; i++ {
		var block Block
		block.BlockData = make([]byte, client.BlockSize)
		n, err := file.Read(block.BlockData)
		if err != nil && err != io.EOF {
			log.Println("Read file error: ", err)
		}
		block.BlockSize = int32(n)
		block.BlockData = block.BlockData[:n]

		var blockStoreAddr string
		err = client.GetBlockStoreAddr(&blockStoreAddr)
		if err != nil {
			log.Println("Fail to get block store address ", err)
		}

		var succ bool
		err = client.PutBlock(&block, blockStoreAddr, &succ)
		if err != nil {
			log.Println("Put block failed: ", err)
		}
	}

	// Update file
	err = client.UpdateFile(fileMetaData, &fileMetaData.Version)

	// fmt.Println("After update-------fileMetaData.version:", fileMetaData.Version)

	if err != nil {
		log.Println("Fail to update file:", err)
		serverFileInfoMap := make(map[string]*FileMetaData)

		// fmt.Println("**upload function** before get file info map")
		// PrintMetaMap(serverFileInfoMap)

		client.GetFileInfoMap(&serverFileInfoMap)

		// fmt.Println("**upload function** after get file info map")
		// PrintMetaMap(serverFileInfoMap)

		updateClientFile(client, serverFileInfoMap[fileMetaData.Filename], indexMap, indexLines)
	}
	return err
}

func download(client RPCClient, fileName string, fileMetaData *FileMetaData) (string, error) {
	filePath := client.BaseDir + "/" + fileName
	if _, e := os.Stat(filePath); os.IsNotExist(e) {
		os.Create(filePath)
	} else {
		os.Truncate(filePath, 0)
	}

	if len(fileMetaData.BlockHashList) == 1 && fileMetaData.BlockHashList[0] == "0" {
		err := os.Remove(filePath)
		if err != nil {
			log.Println("Fail to remove file: ", err)
		}
		return FileMetaDataToString(fileMetaData), err
	}

	file, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755) // Add file access mode.
	defer file.Close()

	hashStr := ""
	var err error
	for _, hash := range fileMetaData.BlockHashList {
		var block Block
		var blockStoreAddr string

		err = client.GetBlockStoreAddr(&blockStoreAddr)
		if err != nil {
			log.Println("Fail to get block store address ", err)
		}

		err = client.GetBlock(hash, blockStoreAddr, &block)
		if err != nil {
			log.Println("Fail to get block: ", err)
		}

		data := string(block.BlockData)
		_, err = io.WriteString(file, data)
		if err != nil {
			log.Println("Fail to write file: ", err)
		}

		hashStr += hash + " "
		// if i != len(fileMetaData.BlockHashList)-1 {
		// 	hashStr += " "
		// }
	}

	line := fileMetaData.Filename + "," + strconv.Itoa(int(fileMetaData.Version)) + "," + hashStr
	return line, err
}
