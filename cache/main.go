// main.go - Ejemplo de cache en memoria para funciones costosas (Fibonacci)
// Autor: afperdomo2
// Fecha: 25 de agosto de 2025

package main

import (
	"fmt"
	"time"
)

// CacheableFunction define el tipo de función que puede ser cacheada.
type CacheableFunction func(key int) (any, error)

// CachedFunctionResult es un tipo que representa el resultado de una función cacheada.
type CachedFunctionResult struct {
	value any   // Valor calculado por la función
	err   error // Error retornado por la función
}

type Memory struct {
	f     CacheableFunction            // Función a cachear
	cache map[int]CachedFunctionResult // Mapa para almacenar resultados cacheados
}

// newMemory inicializa una instancia de Memory con la función a cachear.
func newMemory(f CacheableFunction) *Memory {
	return &Memory{
		f:     f,
		cache: make(map[int]CachedFunctionResult),
	}
}

// Get retorna el valor cacheado para una clave. Si no existe, lo calcula y lo almacena.
func (m *Memory) Get(key int) (any, error) {
	result, isCached := m.cache[key]
	if isCached {
		fmt.Println("[✅Cacheado]")
		return result.value, result.err
	}
	// Calcula el valor y lo almacena en el cache
	result.value, result.err = m.f(key)
	m.cache[key] = result
	fmt.Printf("[⚙️Calculado]\n")
	return result.value, result.err
}

// GetFibonacci adapta la función Fibonacci para el tipo Function.
func GetFibonacci(n int) (any, error) {
	return Fibonacci(n), nil
}

// main ejecuta el ejemplo de cache usando la función Fibonacci.
func main() {
	cache := newMemory(GetFibonacci)
	fibonacciNumbers := []int{35, 40, 44, 40, 45}
	for _, n := range fibonacciNumbers {
		start := time.Now()

		fmt.Printf("\n🔢 Fibonacci de %d... ", n)
		result, err := cache.Get(n)
		if err != nil {
			panic(err)
		}

		fmt.Printf("🔢 Resultado => %v\n", result)
		fmt.Println("⏱️ Time taken:", time.Since(start))
	}
}

// Fibonacci calcula el n-ésimo número de Fibonacci de forma recursiva.
func Fibonacci(n int) int {
	if n <= 0 {
		return 0
	} else if n == 1 {
		return 1
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}
