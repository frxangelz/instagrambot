// global
package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
	//	"strings"
)

var (
	mypath, executable, tmp_path, logfile string
	flog                                  *os.File
	PathSeparator                         string
	ExeName                               string
	Version                               string
)

func isWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	} else {
		return false
	}
}

func IsFileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func CreateDir(dir string) error {
	var err error

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
	}

	return err
}

func int64ToStr(i int64) string {

	return strconv.FormatInt(i, 10)
}

func intToStr(i int) string {

	return strconv.FormatInt(int64(i), 10)
}

func uint64ToStr(i uint64) string {

	return strconv.FormatUint(i, 10)
}

func uintToStr(i uint) string {

	return strconv.FormatUint(uint64(i), 10)
}

func strToUint64(str string) uint64 {

	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}

	return i
}

func strToInt64(str string) int64 {

	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}

	return i
}

func strToInt(str string) int {

	i, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0
	}

	return int(i)
}

func strToUint(str string) uint {

	i, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0
	}

	return uint(i)
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func strToint64(str string) int64 {

	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}

	return i
}

func Initialize() {

	Version = "1.0"
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	if isWindows() {
		mypath = filepath.Dir(ex) + "\\"
		tmp_path = mypath + "tmp" + "\\"
		PathSeparator = "\\"
	} else {
		mypath = filepath.Dir(ex) + "/"
		tmp_path = mypath + "tmp" + "/"
		PathSeparator = "/"
	}
	executable = ex
	ExeName = filepath.Base(executable)
	logfile = mypath + "logs.txt"

	if IsFileExists(logfile) {

		fi, e := os.Stat(logfile)

		if e == nil {
			// get the size
			size := fi.Size()
			if size > 1024*10 {
				os.Remove(logfile)
			}

		}
	}

	flog, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//flog, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(flog)
}

func getChance(chance int) bool {

	if chance >= 100 {
		return true
	}

	rand.Seed(time.Now().Unix())
	return rand.Intn(100-0)+0 <= chance
}

func CleanUp() {

	flog.Close()
}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func get_images(path string) ([]string, error) {

	images, err := filepath.Glob(path + "*")
	if err != nil {
		return images, err
	}

	return images, err
}

var log_counter int = 0

func mylog(args ...interface{}) {

	if config.debug_to_console {
		fmt.Println(args...)
	}

	if log_counter > config.max_log_lines {

		var err error = nil
		flog.Close()
		os.Remove(logfile)
		// re-opening
		flog, err = os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			//			log.Fatalf("error opening file: %v", err)
		} else {

			log.SetOutput(flog)
		}
		log_counter = 0
	}

	if config.debug_to_file {
		log.Println(args...)
		log_counter++
	}
}

// return a string containing the file name, function name
// and the line number of a specified entry on the call stack
/*
func WhereAmI(depthList ...int) string {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	function, file, line, _ := runtime.Caller(depth)
	return fmt.Sprintf("File: %s  Function: %s Line: %d", chopPath(file), runtime.FuncForPC(function).Name(), line)
}

func WAmI(depthList ...int) string {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	function, _, line, _ := runtime.Caller(depth)
	return fmt.Sprintf("%s::%d", runtime.FuncForPC(function).Name(), line)
}

// return the source filename after the last slash
func chopPath(original string) string {
	i := strings.LastIndex(original, "/")
	if i == -1 {
		return original
	} else {
		return original[i+1:]
	}
}
*/
