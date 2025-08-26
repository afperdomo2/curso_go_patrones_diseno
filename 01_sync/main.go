// NOTE: Diferencias entre RWLock y no usar nada:

// - Lock bloquea lecturas (con RLock) y escrituras (con Lock) de otras goroutines
// - Unlock permite nuevas lecturas (con Rlock) y/o otra escritura (con Lock)
// - RLock bloquea escrituras (Lock) pero no bloquea lecturas (RLock)
// - RUnlock permite nuevas escrituras (y también lecturas, pero por la naturaleza de RLock,
//   estas no se vieron bloqueadas nunca)
//
// En esencia, RLock de RWLock garantiza una secuencia de lecturas en donde el valor que lees
// no se verá alterado por nuevos escritores, a diferencia de no usar nada.

package main

import (
	"fmt"
	"sync"
)

var balance int = 0

// Deposit agrega una cantidad específica al balance de manera segura utilizando un mutex.
// amount: La cantidad de dinero a depositar.
// wg: Un puntero a un WaitGroup para sincronizar la finalización de la operación.
// Esta función asegura que el acceso al balance sea seguro en un entorno concurrente,
// bloqueando el mutex antes de modificar el balance y desbloqueándolo después.
func Deposit(amount int, wg *sync.WaitGroup, mu *sync.RWMutex) {
	defer wg.Done()
	defer mu.Unlock() // Desbloquea el mutex después de modificar balance

	mu.Lock() // Bloquea el mutex antes de modificar balance
	balance += amount
}

func Balance(wg *sync.WaitGroup, mu *sync.RWMutex) int {
	defer wg.Done()
	defer mu.RUnlock() // Desbloquea el mutex después de leer balance

	mu.RLock() // Bloquea el mutex para lectura
	b := balance
	fmt.Println("✅ Current balance is", b)
	return balance
}

// 1 Deposit() -> Escribiendo (Posible condición de carrera)
// N Balance() -> Muchos leyendo (Seguro)
func main() {
	var wg sync.WaitGroup
	var mu sync.RWMutex

	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go Deposit(i, &wg, &mu)

		wg.Add(1)
		go Balance(&wg, &mu)
	}
	wg.Wait()
	fmt.Println("Balance final:", balance)
}
