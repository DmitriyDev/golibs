package downloader

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Downloader struct {
	counter       int
	total         int
	chankLimit    int
	lastChankId   int
	threadManager ThreadManager
	chunk         LocationChank
	debug         bool
}

func (d Downloader) New(threadCount int, chankLimit int, printDebug bool) Downloader {

	return Downloader{
		counter:       0,
		total:         0,
		chankLimit:    chankLimit,
		lastChankId:   0,
		threadManager: ThreadManager{}.New(threadCount, printDebug),
		debug:         printDebug,
	}

}

func (d *Downloader) reset() {

	d.counter = 0
	d.total = 0
	d.lastChankId = 0

}

func (d *Downloader) ListSource(urlList []string, resultFolder string) {

	d.reset()
	d.total = len(urlList)
	d.chunk = LocationChank{}.New(resultFolder)

	for _, url := range urlList {
		d.counter++

		filename := filepath.Base(url)

		if d.fileExists(resultFolder, filename) {
			if d.debug {
				fmt.Println(url + " \t  --- Exists")
			}

			continue
		}

		d.chunk.Add(Location{
			Filename: filename,
			URL:      url,
		})

		if d.counter%d.chankLimit == 0 || d.counter >= d.total {

			d.threadManager.processChank(d.chunk)
			d.lastChankId++
			d.printLog()
			d.chunk = LocationChank{}.New(resultFolder)

		}

	}

	for !d.threadManager.allThreadStoped() {
		fmt.Printf("Wait all threads fihished. Sleep ...\n")
		time.Sleep(3 * time.Second)
	}
}

func (d *Downloader) FileSource(sourceFile string, resultFolder string) {

	d.reset()
	d.chunk = LocationChank{}.New(resultFolder)

	total, err := d.getUrlCount(sourceFile)
	if err != nil {
		panic(err)
	}
	d.total = total

	f, err := os.Open(sourceFile)

	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	defer f.Close()

	for fileScanner.Scan() {
		d.counter++

		url := fileScanner.Text()
		filename := filepath.Base(url)

		if d.fileExists(resultFolder, filename) {
			if d.debug {
				fmt.Println(url + " \t  --- Exists")
			}

			continue
		}

		d.chunk.Add(Location{
			Filename: filename,
			URL:      url,
		})

		if d.counter%d.chankLimit == 0 || d.counter >= d.total {
			d.threadManager.processChank(d.chunk)
			d.lastChankId++
			d.printLog()
			d.chunk = LocationChank{}.New(resultFolder)
		}

	}

	for !d.threadManager.allThreadStoped() {
		fmt.Printf("Wait all threads fihished. Sleep ...\n")
		time.Sleep(3 * time.Second)
	}
}

func (d *Downloader) getUrlCount(sourceFile string) (int, error) {

	f, err := os.Open(sourceFile)

	defer f.Close()

	if err != nil {
		return 0, err
	}

	r := bufio.NewReader(f)

	var count int
	const lineBreak = '\n'

	buf := make([]byte, bufio.MaxScanTokenSize)

	for {
		bufferSize, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}

		var buffPosition int
		for {
			i := bytes.IndexByte(buf[buffPosition:], lineBreak)
			if i == -1 || bufferSize == buffPosition {
				break
			}
			buffPosition += i + 1
			count++
		}
		if err == io.EOF {
			break
		}
	}

	return count, nil
}

func (d *Downloader) fileExists(folder string, filename string) bool {
	if _, err := os.Stat(folder + filename); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		panic(err)
	}
}

func (d *Downloader) printLog() {
	log := []string{
		fmt.Sprintf("Current chank %d", d.lastChankId),
		fmt.Sprintf("Total %d", d.total),
		fmt.Sprintf("Total Done %d", d.counter),
		fmt.Sprintf("Percent Done %.3f", d.percentDone()) + " %",
	}

	fmt.Println(strings.Join(log, "\t"))
}

func (d *Downloader) percentDone() float64 {
	return float64(d.counter*100) / float64(d.total)
}
