package services

import (
	"BIOCAD/internal/repository"
	"BIOCAD/internal/structures"
	"fmt"
	"github.com/dogenzaka/tsv"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const REPORT_DIRECTORY = "reports"

type ScanDirectoryService struct {
	r *repository.Repository
}

func NewScanDirectory(r *repository.Repository) *ScanDirectoryService {
	return &ScanDirectoryService{r: r}
}

func isNew(checked []structures.File, fileName string) bool {
	for _, file := range checked {
		if file.Name == fileName {
			return false
		}
	}
	return true
}

func (sd *ScanDirectoryService) Scan(dirName string, duration time.Duration) {
	for {
		checkedFiles, err := sd.r.GetAllFiles()
		if err != nil {
			log.Println(err.Error())
		}

		files, err := ioutil.ReadDir(dirName)
		if err != nil {
			log.Fatal(err)
		}

		unitGuidChan := make(chan structures.Device, 100)
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".tsv") {
				if isNew(checkedFiles, file.Name()) {
					go sd.ReadFile(dirName, file.Name(), unitGuidChan)
					_, err := sd.r.AddFile(file.Name())
					if err != nil {
						log.Println(err.Error())
						return
					}
					go sd.MakeReports(unitGuidChan)
				}
			}
		}
		time.Sleep(duration)
	}
}

func (sd *ScanDirectoryService) ReadFile(dirName, fileName string, unitGuidChan chan structures.Device) {
	defer close(unitGuidChan)
	file, err := os.Open(path.Join(dirName, fileName))
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer file.Close()

	data := structures.Device{}
	parser, _ := tsv.NewParser(file, &data)

	for {
		eof, err := parser.Next()
		if eof {
			return
		}
		if err != nil {
			return
		}
		if data.UnitGuid == "" {
			return
		}
		data.Id = 0
		_, err = sd.r.AddDevice(data)
		if err != nil {
			log.Println(err.Error())
			return
		}
		unitGuidChan <- data
	}
}

func (sd *ScanDirectoryService) MakeReports(unitGuidChan chan structures.Device) {
	if _, err := os.Stat(REPORT_DIRECTORY); os.IsNotExist(err) {
		os.Mkdir(REPORT_DIRECTORY, 755)
	}
	os.Chdir(REPORT_DIRECTORY)
	defer os.Chdir("..")
	var wg sync.WaitGroup
	for {
		i, ok := <-unitGuidChan
		if i.UnitGuid == "" && !ok {
			break
		}
		wg.Add(1)
		go sd.MakeReportFile(i, &wg)
	}
	wg.Wait()
}

func (sd *ScanDirectoryService) MakeReportFile(device structures.Device, wg *sync.WaitGroup) {
	defer wg.Done()
	var file *os.File
	file, _ = os.OpenFile(device.UnitGuid+".doc", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	_, err := file.WriteString(fmt.Sprintf("%v\n", device))
	if err != nil {
		log.Println("Unable to write to file:", err)
		return
	}
}
