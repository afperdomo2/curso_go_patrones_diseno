package main

import (
	"fmt"
	"sync"
	"time"
)

var mu sync.Mutex

type DataBase struct {
	connectionString string
}

func (db *DataBase) Connect() {
	fmt.Println("🔗 Connecting to database...")
	time.Sleep(2 * time.Second)
	fmt.Println("✅ Connected to database!")
}

var instance *DataBase

func GetDataBaseInstance() *DataBase {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		fmt.Printf("🧪 Creating new database instance...\n")
		instance = &DataBase{}
		instance.Connect()
	} else {
		fmt.Printf("🔍 Reusing existing database instance...\n")
	}
	return instance
}

func main() {
	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			GetDataBaseInstance()
		}(i)
	}
	wg.Wait()
	fmt.Println("All goroutines finished.")
}
