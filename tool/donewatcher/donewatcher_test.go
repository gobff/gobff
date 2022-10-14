package donewatcher

import (
	"fmt"
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	watcher := NewWatcher()
	go func() {
		watcher.Wait([]string{"test1", "test2", "test3"})
		fmt.Println("done:", "test1", "test2", "test3")
	}()
	go func() {
		watcher.Wait([]string{"test1", "test2"})
		fmt.Println("done:", "test1", "test2")
	}()
	go func() {
		watcher.Wait([]string{"test2", "test3"})
		fmt.Println("done:", "test2", "test3")
	}()
	go func() {
		watcher.Wait([]string{"test1", "test3"})
		fmt.Println("done:", "test1", "test3")
	}()
	go func() {
		watcher.Wait([]string{"test1"})
		fmt.Println("done:", "test1")
	}()
	go func() {
		watcher.Wait([]string{"test2"})
		fmt.Println("done:", "test2")
	}()
	go func() {
		watcher.Wait([]string{"test3"})
		fmt.Println("done:", "test3")
	}()

	time.Sleep(5 * time.Second)
	watcher.Done("test1")
	time.Sleep(5 * time.Second)
	watcher.Done("test2")
	time.Sleep(5 * time.Second)
	watcher.Done("test3")
	time.Sleep(5 * time.Second)
}
