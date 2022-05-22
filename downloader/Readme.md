# Downloader

# Usage 

### Install

```go get github.com/DmitriyDev/golibs/downloader```

### File of links


```
package main


import (
	dwn "github.com/DmitriyDev/golibs/downloader"
)
const DEBUG = false
const THREADS = 3
const CHUNK_SIZE = 10

func main() {

	dw := Downloader{}.New(THREADS, CHUNK_SIZE, DEBUG)
	dw.FileSource("url_list.txt", "./temp/")

}
```
