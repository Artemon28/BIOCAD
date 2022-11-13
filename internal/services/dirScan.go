package services

import (
	"BIOCAD/internal/repository"
	"BIOCAD/internal/structures"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strconv"
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

		for _, file := range files {
			unitGuidChan := make(chan structures.Device, 100)
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
	tsvReader := csv.NewReader(file)
	tsvReader.Comma = '\t'
	_, err = tsvReader.Read()
	if err != nil {
		log.Println(err.Error())
	} //read header

	for {
		rec, err := tsvReader.Read()
		if err == io.EOF {
			break
		}
		if len(rec) == 0 {
			log.Println("empty row in file")
			continue
		}
		data = convertRecord(rec)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if data.UnitGuid == "" {
			log.Println("Can't recognise this device: " + fmt.Sprintf("%v", data) + " guid is nil")
			continue
		}
		data.Id = 0
		_, err = sd.r.AddDevice(data)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		unitGuidChan <- data
	}
}

func convertRecord(rec []string) structures.Device {
	var dev structures.Device
	if len(rec) > reflect.ValueOf(dev).NumField() {
		rec = rec[:reflect.ValueOf(dev).NumField()]
	}
	for len(rec) < reflect.ValueOf(dev).NumField() {
		rec = append(rec, "")
	}
	dev.Mqtt = rec[1]
	dev.Invid = rec[2]
	dev.UnitGuid = rec[3]
	dev.MsgId = rec[4]
	dev.Text = rec[5]
	dev.Context = rec[6]
	dev.Class = rec[7]
	dev.Level, _ = strconv.Atoi(rec[8])
	dev.Area = rec[9]
	dev.Addr = rec[10]
	dev.Block = rec[11]
	dev.Type = rec[12]
	dev.Bit, _ = strconv.Atoi(rec[13])
	dev.InvertBit, _ = strconv.Atoi(rec[14])
	dev.Id = 0
	return dev
}

func (sd *ScanDirectoryService) MakeReports(unitGuidChan chan structures.Device) {
	if _, err := os.Stat(REPORT_DIRECTORY); os.IsNotExist(err) {
		err := os.Mkdir(REPORT_DIRECTORY, 0755)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
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
	file, _ = os.OpenFile(path.Join(REPORT_DIRECTORY, device.UnitGuid+".doc"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	_, err := file.WriteString(fmt.Sprintf("%v\n", device))
	if err != nil {
		log.Println("Unable to write to file:", err)
		return
	}
}
