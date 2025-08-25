package main

import (
	"fmt"
	"sync"
)

var balance int = 100
var mu sync.Mutex // Mutex para proteger el acceso a balance

// Deposit agrega una cantidad específica al balance de manera segura utilizando un mutex.
// amount: La cantidad de dinero a depositar.
// wg: Un puntero a un WaitGroup para sincronizar la finalización de la operación.
// Esta función asegura que el acceso al balance sea seguro en un entorno concurrente,
// bloqueando el mutex antes de modificar el balance y desbloqueándolo después.
func Deposit(amount int, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock() // Bloquea el mutex antes de modificar balance
	balance += amount
	mu.Unlock() // Desbloquea el mutex después de modificar balance
}

func Balance() int {
	mu.Lock()         // Bloquea el mutex antes de leer balance
	defer mu.Unlock() // Desbloquea el mutex después de leer balance
	return balance
}

func main() {
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go Deposit(i*100, &wg)
	}
	wg.Wait()
	fmt.Println("Balance final:", Balance())
}
