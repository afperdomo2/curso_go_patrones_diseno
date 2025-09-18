/*
Patrón de Diseño Adapter - Ejemplo en Go

El patrón Adapter es un patrón de diseño estructural que permite que objetos
con interfaces incompatibles trabajen juntos. Actúa como un puente entre dos
interfaces incompatibles.

Problema que resuelve:
- Tenemos una clase existente (CreditCardPayment) que tiene una interfaz incompatible
- Necesitamos usar esta clase con un sistema que espera una interfaz diferente (IPayment)
- No podemos modificar la clase existente

Solución:
- Crear una clase adaptadora (CreditCardAdapter) que:
 1. Implementa la interfaz esperada (IPayment)
 2. Contiene una instancia de la clase incompatible (CreditCardPayment)
 3. Traduce las llamadas de una interfaz a otra

Ventajas:
- Permite reutilizar código existente sin modificarlo
- Separa la lógica de conversión de interfaces de la lógica de negocio
- Facilita la integración de bibliotecas externas

En este ejemplo:
- IPayment: Interfaz objetivo que esperan los clientes
- CashPayment: Implementación que ya cumple con IPayment
- CreditCardPayment: Clase incompatible que necesita adaptación
- CreditCardAdapter: Adaptador que hace compatible CreditCardPayment con IPayment
*/
package main

import "fmt"

// 1. Definición de la interfaz IPayment

// IPayment define la interfaz objetivo que esperan los clientes
// Todos los métodos de pago deben implementar esta interfaz
type IPayment interface {
	Pay() // Método estándar que todos los pagos deben implementar
}

// CashPayment representa un pago en efectivo que ya es compatible con IPayment
// Esta clase NO necesita adaptador porque ya implementa la interfaz correcta
type CashPayment struct{}

// Pay implementa directamente la interfaz IPayment para pagos en efectivo
func (c CashPayment) Pay() {
	fmt.Println("💰 Pagando con efectivo")
}

// ProcessPayment es una función que puede trabajar con cualquier tipo de pago
// que implemente la interfaz IPayment. Demuestra el polimorfismo.
func ProcessPayment(p IPayment) {
	p.Pay()
}

// 2. Definición de la clase incompatible y el adaptador

// CreditCardPayment representa un sistema de pago con tarjeta de crédito INCOMPATIBLE
// Esta clase tiene una interfaz diferente (requiere userAccountID como parámetro)
// NO puede ser usado directamente con IPayment - necesita un adaptador
type CreditCardPayment struct{}

// Pay es el método original de CreditCardPayment que NO es compatible con IPayment
// Requiere un parámetro userAccountID, mientras que IPayment.Pay() no requiere parámetros
func (CreditCardPayment) Pay(userAccountID int) {
	fmt.Printf("💳 Pagando desde la cuenta de usuario %d usando tarjeta de crédito\n", userAccountID)
}

// CreditCardPaymentAdapter es el ADAPTADOR que hace compatible CreditCardPayment con IPayment
// Implementa el patrón Adapter siguiendo estos principios:
// 1. Implementa la interfaz objetivo (IPayment)
// 2. Contiene una instancia del objeto incompatible (CreditCardPayment)
// 3. Almacena los datos necesarios para la adaptación (UserAccountID)
type CreditCardPaymentAdapter struct {
	CreditCardPayment *CreditCardPayment // La clase incompatible que queremos adaptar
	UserAccountID     int                // Datos adicionales necesarios para la adaptación
}

// Pay implementa la interfaz IPayment en el adaptador
// Esta es la "traducción" que hace que CreditCardPayment sea compatible con IPayment
// El adaptador toma la llamada sin parámetros de IPayment.Pay() y la convierte
// en una llamada con parámetros a CreditCardPayment.Pay(userAccountID)
func (cca CreditCardPaymentAdapter) Pay() {
	cca.CreditCardPayment.Pay(cca.UserAccountID)
}

// 3. Demostración adicional con otro método de pago incompatible

type BankPayment struct{}

func (b BankPayment) Pay(accountNumber string) {
	fmt.Printf("🏦 Pagando desde la cuenta bancaria %s\n", accountNumber)
}

type BankPaymentAdapter struct {
	BankPayment   *BankPayment
	AccountNumber string
}

func (ba BankPaymentAdapter) Pay() {
	ba.BankPayment.Pay(ba.AccountNumber)
}

// main demuestra el uso del patrón Adapter
func main() {
	// 🔄 Ejemplo 1: Usar CashPayment directamente (ya compatible con IPayment)
	fmt.Println("🟢 Procesando pago directo (sin adaptador):")
	cash := &CashPayment{}
	ProcessPayment(cash)

	fmt.Println("\n🔧 Procesando pago con adaptador:")
	// 🔄 Ejemplo 2: Usar CreditCardPayment a través del adaptador
	ccpa := &CreditCardPaymentAdapter{
		CreditCardPayment: &CreditCardPayment{},
		UserAccountID:     12345,
	}
	ProcessPayment(ccpa)

	fmt.Println("\n🔧 Procesando pago bancario con adaptador:")
	// 🔄 Ejemplo 3: Usar BankPayment a través del adaptador
	bpa := &BankPaymentAdapter{
		BankPayment:   &BankPayment{},
		AccountNumber: "987654321",
	}
	ProcessPayment(bpa)
}
