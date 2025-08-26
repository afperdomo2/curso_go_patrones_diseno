package main

import (
	"fmt"
	"sync"
	"time"
)

// CacheItem representa un elemento en el cache con su valor y tiempo de expiración
// Esta estructura encapsula el valor almacenado junto con metadatos básicos
type CacheItem struct {
	Value      any   // El valor que se almacena (puede ser cualquier tipo de dato)
	Expiration int64 // Timestamp de cuando expira (0 significa que nunca expira)
}

// IsExpired verifica si el elemento del cache ha expirado
// Retorna true si el elemento debe considerarse como eliminado
func (item *CacheItem) IsExpired() bool {
	if item.Expiration == 0 {
		return false // Si es 0, nunca expira
	}
	return time.Now().UnixNano() > item.Expiration
}

// SimpleRedisCache implementa un cache básico en memoria similar a Redis
// Usa un mapa simple para almacenar los datos y un mutex para thread-safety
type SimpleRedisCache struct {
	data  map[string]*CacheItem // Mapa que contiene todos los elementos del cache
	mutex sync.RWMutex          // Mutex para permitir acceso concurrente seguro
}

// NewSimpleRedisCache crea y retorna una nueva instancia del cache
// Inicializa el mapa interno para almacenar los datos
func NewSimpleRedisCache() *SimpleRedisCache {
	return &SimpleRedisCache{
		data: make(map[string]*CacheItem),
	}
}

// Set almacena un valor en el cache con una clave específica
// Parámetros:
//   - key: la clave para identificar el elemento
//   - value: el valor a almacenar (puede ser cualquier tipo)
//   - ttl: tiempo de vida del elemento (time.Duration, 0 = nunca expira)
func (c *SimpleRedisCache) Set(key string, value any, ttl time.Duration) {
	// Bloquear para escritura (exclusivo)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var expiration int64
	if ttl > 0 {
		// Calcular el timestamp de expiración
		expiration = time.Now().Add(ttl).UnixNano()
	}

	// Crear el elemento y almacenarlo en el mapa
	c.data[key] = &CacheItem{
		Value:      value,
		Expiration: expiration,
	}

	fmt.Printf("✅ SET '%s' = '%v'", key, value)
	if ttl > 0 {
		fmt.Printf(" (expira en %v)", ttl)
	}
	fmt.Println()
}

// Get recupera un valor del cache usando su clave
// Retorna:
//   - any: el valor almacenado
//   - bool: true si la clave existe y no ha expirado, false en caso contrario
func (c *SimpleRedisCache) Get(key string) (any, bool) {
	// Bloquear para lectura (permite múltiples lectores concurrentes)
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Buscar el elemento en el mapa
	item, exists := c.data[key]
	if !exists {
		fmt.Printf("❌ GET '%s' - Clave no encontrada\n", key)
		return nil, false
	}

	// Verificar si el elemento ha expirado
	if item.IsExpired() {
		fmt.Printf("⏰ GET '%s' - Clave expirada\n", key)
		return nil, false
	}

	fmt.Printf("✅ GET '%s' = '%v'\n", key, item.Value)
	return item.Value, true
}

// Delete elimina un elemento del cache
// Retorna true si el elemento existía y fue eliminado, false si no existía
func (c *SimpleRedisCache) Delete(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Verificar si la clave existe antes de eliminarla
	if _, exists := c.data[key]; exists {
		delete(c.data, key)
		fmt.Printf("🗑️ DELETE '%s' - Eliminado exitosamente\n", key)
		return true
	}

	fmt.Printf("❌ DELETE '%s' - Clave no encontrada\n", key)
	return false
}

// Exists verifica si una clave existe en el cache y no ha expirado
func (c *SimpleRedisCache) Exists(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists || item.IsExpired() {
		fmt.Printf("❌ EXISTS '%s' - No existe o expiró\n", key)
		return false
	}

	fmt.Printf("✅ EXISTS '%s' - Existe\n", key)
	return true
}

// Size retorna el número de elementos actualmente en el cache
func (c *SimpleRedisCache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.data)
}

// demonstrateBasicOperations muestra las operaciones básicas del cache
func demonstrateBasicOperations() {
	fmt.Println("🚀 === DEMOSTRACIÓN BÁSICA DEL CACHE REDIS === 🚀")

	// Crear una nueva instancia del cache
	cache := NewSimpleRedisCache()

	fmt.Println("📝 1. Operaciones de ESCRITURA (SET):")
	fmt.Println("   - Almacenar datos con y sin expiración")

	// Almacenar diferentes tipos de datos
	cache.Set("nombre", "Juan Pérez", 0)                        // String sin expiración
	cache.Set("edad", 25, 0)                                    // Entero sin expiración
	cache.Set("activo", true, 0)                                // Boolean sin expiración
	cache.Set("temporal", "Este valor expirará", 3*time.Second) // String con expiración

	fmt.Printf("\n📊 Tamaño del cache después de SET: %d elementos\n\n", cache.Size())

	fmt.Println("🔍 2. Operaciones de LECTURA (GET):")
	fmt.Println("   - Recuperar datos almacenados")

	// Leer los valores almacenados
	if valor, found := cache.Get("nombre"); found {
		fmt.Printf("   Nombre recuperado: %s\n", valor)
	}

	if valor, found := cache.Get("edad"); found {
		fmt.Printf("   Edad recuperada: %d\n", valor)
	}

	if valor, found := cache.Get("activo"); found {
		fmt.Printf("   Estado activo: %t\n", valor)
	}

	// Intentar leer una clave que no existe
	cache.Get("clave_inexistente")

	fmt.Println("\n⏰ 3. Demostración de EXPIRACIÓN:")
	fmt.Println("   - Verificar comportamiento con TTL")

	// Verificar que el valor temporal existe
	cache.Exists("temporal")

	fmt.Println("   Esperando 4 segundos para que expire...")
	time.Sleep(4 * time.Second)

	// Verificar que el valor ha expirado
	cache.Get("temporal")
	cache.Exists("temporal")

	fmt.Println("\n🗑️ 4. Operaciones de ELIMINACIÓN (DELETE):")
	cache.Delete("edad")
	cache.Delete("clave_inexistente") // Intentar eliminar algo que no existe

	fmt.Printf("\n📊 Tamaño final del cache: %d elementos\n", cache.Size())
}

// demonstrateConcurrency demuestra que el cache es seguro para uso concurrente
func demonstrateConcurrency() {
	fmt.Println("\n🔄 === DEMOSTRACIÓN DE CONCURRENCIA === 🔄")

	cache := NewSimpleRedisCache()
	var wg sync.WaitGroup

	// Función que simula escrituras concurrentes
	writer := func(id int) {
		defer wg.Done()
		for i := range 5 {
			key := fmt.Sprintf("worker_%d_item_%d", id, i)
			value := fmt.Sprintf("valor_%d_%d", id, i)
			cache.Set(key, value, 10*time.Second)
			time.Sleep(100 * time.Millisecond) // Pequeña pausa para simular trabajo real
		}
	}

	// Función que simula lecturas concurrentes
	reader := func(id int) {
		defer wg.Done()
		for i := range 5 {
			key := fmt.Sprintf("worker_%d_item_%d", id%3, i) // Leer de diferentes workers
			cache.Get(key)
			time.Sleep(150 * time.Millisecond)
		}
	}

	fmt.Println("Iniciando 3 workers de escritura y 2 workers de lectura...")

	// Iniciar workers de escritura
	for i := range 3 {
		wg.Add(1)
		go writer(i)
	}

	// Esperar un poco antes de iniciar los lectores
	time.Sleep(200 * time.Millisecond)

	// Iniciar workers de lectura
	for i := range 2 {
		wg.Add(1)
		go reader(i)
	}

	// Esperar a que todos terminen
	wg.Wait()

	fmt.Printf("\n✅ Operaciones concurrentes completadas. Tamaño final: %d elementos\n", cache.Size())
}

// main función principal que ejecuta todas las demostraciones
func main() {
	fmt.Println("🎯 Sistema de Cache Estilo Redis - Versión Educativa")
	fmt.Println("=====================================================")
	fmt.Println()
	fmt.Println("📚 CONCEPTOS IMPLEMENTADOS:")
	fmt.Println("   • Almacenamiento key-value en memoria")
	fmt.Println("   • Thread-safety con sync.RWMutex")
	fmt.Println("   • Expiración automática de elementos (TTL)")
	fmt.Println("   • Operaciones básicas: SET, GET, DELETE, EXISTS")
	fmt.Println()

	// Ejecutar demostración básica
	demonstrateBasicOperations()

	// Ejecutar demostración de concurrencia
	demonstrateConcurrency()

	fmt.Println("\n🎉 ¡Demostración completada!")
	fmt.Println("\n💡 PUNTOS CLAVE APRENDIDOS:")
	fmt.Println("   1. Un cache es un almacén temporal de datos en memoria")
	fmt.Println("   2. Redis usa el patrón key-value para almacenar datos")
	fmt.Println("   3. TTL (Time To Live) permite que los datos expiren automáticamente")
	fmt.Println("   4. Los mutex garantizan operaciones thread-safe en entornos concurrentes")
	fmt.Println("   5. RWMutex permite múltiples lectores simultáneos pero escritores exclusivos")
}
