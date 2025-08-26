package main

import (
	"fmt"
	"sync"
	"time"
)

type Service struct {
	InProgress map[int]bool
	IsPending  map[int][]chan int
	mu         sync.RWMutex
}

func newService() *Service {
	return &Service{
		InProgress: make(map[int]bool),
		IsPending:  make(map[int][]chan int),
	}
}

func (s *Service) Work(job int) {
	s.mu.RLock()

	isJobInProgress := s.InProgress[job]
	if isJobInProgress {
		s.mu.RUnlock()
		response := make(chan int)
		defer close(response)

		s.mu.Lock()
		s.IsPending[job] = append(s.IsPending[job], response)
		s.mu.Unlock()

		fmt.Printf("⏳ Esperando resultado de Fibonacci de %d\n", job)

		resp := <-response

		fmt.Printf("✅ Resultado recibido de Fibonacci de %d: %d\n", job, resp)
		return
	}
	s.mu.RUnlock()

	s.mu.Lock()
	// Si no está en progreso, lo marcamos como tal y comenzamos el trabajo
	s.InProgress[job] = true
	s.mu.Unlock()

	result := ExpensiveFibonacci(job)

	s.mu.RLock()
	pendingWorkers, exists := s.IsPending[job]
	s.mu.RUnlock()

	if exists {
		for _, ch := range pendingWorkers {
			ch <- result
		}
		fmt.Printf("🔔 Notificados a todos los pendientes de Fibonacci de %d\n", job)
	}

	s.mu.Lock()
	s.InProgress[job] = false
	s.IsPending[job] = make([]chan int, 0)
	s.mu.Unlock()
}

// main ejecuta varios trabajos concurrentes usando goroutines y un servicio que gestiona el estado de los trabajos.
// El objetivo es evitar cálculos duplicados y notificar a los clientes cuando el resultado esté disponible.
func main() {
	service := newService()               // Instancia el servicio que gestiona los trabajos concurrentes
	jobs := []int{3, 4, 5, 5, 4, 8, 8, 8} // Lista de trabajos a ejecutar (con repetidos para simular concurrencia)

	var wg sync.WaitGroup // WaitGroup para esperar a que todas las goroutines terminen
	wg.Add(len(jobs))
	for _, job := range jobs {
		go func(j int) {
			defer wg.Done() // Marca la goroutine como finalizada
			service.Work(j) // Ejecuta el trabajo y gestiona la sincronización y notificación
		}(job)
	}
	wg.Wait() // Espera a que todas las goroutines finalicen
}

func ExpensiveFibonacci(n int) int {
	fmt.Printf("⚙️ Calculando Fibonacci de %d...\n", n)
	time.Sleep(5 * time.Second)
	return n
}
