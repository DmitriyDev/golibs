# Downloader

# Usage 

### File of links

```
package downloader

const DEBUG = false
const THREADS = 3
const CHUNK_SIZE = 10



func main() {

	dw := Downloader{}.New(THREADS, CHUNK_SIZE, DEBUG)
	dw.FileSource("url_list.txt", "./temp/")

}




    threads := 3
    chunkSize := 10

	dw := Downloader{}.New(threads, chunkSize)
	dw.FileSource("hash.txt", "./temp/")
```
