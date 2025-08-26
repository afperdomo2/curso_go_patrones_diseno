/*
Patrón de Diseño Factory - Ejemplo en Go

El patrón Factory es un patrón de diseño creacional que proporciona una interfaz
para crear objetos en una superclase, pero permite a las subclases alterar el
tipo de objetos que se crearán.

Ventajas:
- Evita el acoplamiento fuerte entre el código cliente y las clases concretas
- Facilita la extensión del código para agregar nuevos tipos de productos
- Centraliza la lógica de creación de objetos
- Permite la reutilización de código

En este ejemplo:
- IProduct es la interfaz común para todos los productos
- Computer es la estructura base que implementa la funcionalidad común
- Laptop y Desktop son productos concretos que extienden Computer
- GetComputerFactory es la función factory que retorna constructores específicos
*/
package main

import "fmt"

// IProduct define la interfaz común para todos los productos que puede crear la factory
// Establece el contrato que deben cumplir todos los productos concretos
type IProduct interface {
	setStock(stock int)
	getStock() int
	setName(name string)
	getName() string
}

// Computer es la estructura base que contiene los campos comunes
// para todos los tipos de computadoras (Laptop y Desktop)
type Computer struct {
	name  string
	stock int
}

func (c *Computer) setStock(stock int) {
	c.stock = stock
}

func (c *Computer) getStock() int {
	return c.stock
}

func (c *Computer) setName(name string) {
	c.name = name
}

func (c *Computer) getName() string {
	return c.name
}

// Laptop representa un producto concreto de tipo laptop
// Utiliza composición para heredar funcionalidad de Computer
type Laptop struct {
	Computer
}

// NewLaptop es el constructor para crear instancias de Laptop
// Retorna una interfaz IProduct para mantener el polimorfismo
func NewLaptop(name string, stock int) IProduct {
	return &Laptop{
		Computer: Computer{
			name:  name,
			stock: stock,
		},
	}
}

// Desktop representa un producto concreto de tipo computadora de escritorio
// Utiliza composición para heredar funcionalidad de Computer
type Desktop struct {
	Computer
}

// NewDesktop es el constructor para crear instancias de Desktop
// Retorna una interfaz IProduct para mantener el polimorfismo
func NewDesktop(name string, stock int) IProduct {
	return &Desktop{
		Computer: Computer{
			name:  name,
			stock: stock,
		},
	}
}

// GetComputerFactory es la función factory principal del patrón
// Retorna una función constructora específica basada en el tipo solicitado
// Parámetros:
//   - ComputerType: string que especifica el tipo de computadora ("laptop" o "desktop")
// Retorna:
//   - Una función constructora específica para el tipo solicitado
//   - Un error si el tipo no es válido
func GetComputerFactory(ComputerType string) (func(name string, stock int) IProduct, error) {
	if ComputerType == "laptop" {
		return NewLaptop, nil
	}
	if ComputerType == "desktop" {
		return NewDesktop, nil
	}
	return nil, fmt.Errorf("❌ Invalid computer type: %s", ComputerType)
}

// printNameAndStock es una función auxiliar para mostrar información del producto
// Demuestra el polimorfismo al trabajar con la interfaz IProduct
func printNameAndStock(product IProduct) {
	fmt.Printf("📦 Product Name: %s, 📊 Stock: %d\n", product.getName(), product.getStock())
}

// main demuestra el uso del patrón Factory
func main() {
	// 1. Obtener la función factory para laptops
	laptopFactory, err := GetComputerFactory("laptop")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 2. Crear un producto laptop usando la factory
	laptop := laptopFactory("MacBook Pro", 10)
	printNameAndStock(laptop)

	// 3. Obtener la función factory para computadoras de escritorio
	desktopFactory, err := GetComputerFactory("desktop")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 4. Crear productos desktop usando la misma factory
	iMacDesktop := desktopFactory("iMac", 5)
	printNameAndStock(iMacDesktop)

	legionDesktop := desktopFactory("Lenovo Legion", 8)
	printNameAndStock(legionDesktop)
}
