package downloader

import (
	"fmt"
	"time"
)

type ThreadManager struct {
	statuses   map[int]bool
	printDebug bool
}

func (tm ThreadManager) New(count int, debug bool) ThreadManager {

	threadStatuses := map[int]bool{}

	for i := 0; i < count; i++ {
		threadStatuses[i] = true
	}

	return ThreadManager{
		statuses:   threadStatuses,
		printDebug: debug,
	}
}

func (tm *ThreadManager) getEmptyThreadId() int {

	for !tm.hasEmptyThread() {
		fmt.Printf("All threads working. Sleep ...\n")
		time.Sleep(3 * time.Second)
	}

	for id, isEmpty := range tm.statuses {
		if isEmpty {
			tm.statuses[id] = false
			return id
		}
	}
	return tm.getEmptyThreadId()
}

func (tm *ThreadManager) hasEmptyThread() bool {

	for _, isEmpty := range tm.statuses {
		if isEmpty {
			return true
		}
	}
	return false
}

func (tm *ThreadManager) allThreadStoped() bool {
	for _, isEmpty := range tm.statuses {
		if !isEmpty {
			return false
		}
	}
	return true
}

func (tm *ThreadManager) processChank(lc LocationChank) {

	threadId := tm.getEmptyThreadId()

	go tm.runThread(threadId, lc)

}

func (tm *ThreadManager) runThread(threadId int, chank LocationChank) {

	threadDone := make(chan int)

	threadFunc := func(thread int, chank LocationChank, threadDone chan int) {
		defer close(threadDone)

		cs := make(chan string)
		go chank.Run(thread, cs)
		for results := range cs {
			if tm.printDebug {
				fmt.Println(results)
			}
		}
		threadDone <- thread
	}

	go threadFunc(threadId, chank, threadDone)

	for threadId := range threadDone {
		if tm.printDebug {
			fmt.Printf("Thread %d done \n", threadId)
		}

	}

	tm.statuses[threadId] = true
}
