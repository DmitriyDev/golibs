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

	dw := dwn.Downloader{}.New(THREADS, CHUNK_SIZE, DEBUG)
	dw.FileSource("url_list.txt", "./temp/")
	// or
	dw.ListSource([]string{"http://example/test1.csv", "http://example/test2.csv", "http://example/test3.csv"}, "./temp/")
	

}
```

#### url_list.txt example

```
http://example/test1.csv
http://example/test2.csv
http://example/test3.csv
http://example/test4.csv
```