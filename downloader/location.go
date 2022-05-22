package downloader

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Location struct {
	Filename string
	URL      string
}

type LocationChank struct {
	ResultFolder string
	Locations    []Location
}

func (lc LocationChank) New() LocationChank {
	return LocationChank{
		Locations: []Location{},
	}
}

func (lc *LocationChank) Add(l Location) {
	lc.Locations = append(lc.Locations, l)
}

func (lc *LocationChank) Run(thread int, resChan chan string) {
	defer close(resChan)

	for _, loc := range lc.Locations {

		status := []string{loc.URL}

		data, err := download(loc.URL)
		if err != nil {
			status = append(status, "download error", err.Error())
		} else {
			err = os.WriteFile(lc.ResultFolder+loc.Filename, data, 0644)
			if err != nil {
				status = append(status, "write error", err.Error())
			} else {
				status = append(status, "write OK")
			}
		}

		log := strings.Join(status, "\t -- ")

		resChan <- fmt.Sprintf("%d -- %s", thread, log)
	}
}

func download(url string) ([]byte, error) {

	response, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return []byte{}, errors.New("Invalid url")
	}
	return ioutil.ReadAll(response.Body)

}
