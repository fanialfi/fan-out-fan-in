package filegeneration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fanialfi/fan-out-fan-in/lib"
)

func GenerateFileConcurrency() {
	err := os.RemoveAll(lib.TempPath)
	if err != nil {
		log.Printf("ERROR remove all directory and file : %s\n", err.Error())
	}

	err = os.MkdirAll(lib.TempPath, os.ModePerm)
	if err != nil {
		log.Printf("ERROR create directory : %s\n", err.Error())
	}

	// pipeline 1 := job distribution
	chanFileIndex := generateFileIndex()

	// pipeline 2 := the main logic (creating file)
	createFileWorker := 20
	chanFileResult := createFile(chanFileIndex, createFileWorker)

	// track and print output
	counterTotal := 0
	counterSuccess := 0
	for fileResult := range chanFileResult {
		if fileResult.Err != nil {
			log.Printf("error creating file %s\n\tstack trace : %s\n", fileResult.FileName, fileResult.Err.Error())
		} else {
			counterSuccess++
		}

		counterTotal++
	}

	log.Printf("%d/%d of file created", counterSuccess, counterTotal)
}

func generateFileIndex() <-chan lib.FileInfo {
	chanOut := make(chan lib.FileInfo)

	go func() {
		for i := 0; i < lib.TotalFile; i++ {
			chanOut <- lib.FileInfo{
				Index:    i,
				FileName: fmt.Sprintf("file-%d.txt", i),
			}
		}
		close(chanOut)
	}()

	return chanOut
}

// sebagai fan-in, fan-out
func createFile(chanIn <-chan lib.FileInfo, numberOfWorkers int) <-chan lib.FileInfo {
	chanOut := make(chan lib.FileInfo)

	// waitGroup to controll the worker
	wg := new(sync.WaitGroup)

	// allocate N of workers
	wg.Add(numberOfWorkers)
	go func() {
		// dispath N worker as fan-out
		for workerIndex := 0; workerIndex < numberOfWorkers; workerIndex++ {
			go func(workerIndex int) {
				// listen to chanIn channel for incoming jobs
				for job := range chanIn {
					// do the jobs (lakukan task)
					filePath := filepath.Join(lib.TempPath, job.FileName)
					content := lib.RandomString(lib.ContentLength)

					err := os.WriteFile(filePath, []byte(content), os.ModePerm)

					// log.Printf("worker %d working on %s file generation", workerIndex, job.FileName)

					// construct the job's result and send it to chanOut as fan-in
					chanOut <- lib.FileInfo{
						FileName:    job.FileName,
						WorkerIndex: workerIndex,
						Err:         err,
						Index:       job.Index,
					}
				}

				// if chanIn is closed and the remaining jobs is finished
				// only then we mark the worker as complete
				wg.Done()
			}(workerIndex)
		}
	}()

	// wait until chanIn closed and then all workers are done
	// because right after that -  we need to close the chanOut channel
	go func() {
		wg.Wait()
		close(chanOut)
	}()

	return chanOut
}
