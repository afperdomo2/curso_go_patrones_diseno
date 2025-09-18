package main

import "fmt"

// Subject-Observer Pattern - Ejemplo en Go

// 1. Subject

// 1.1 Subject: Definición de la interfaz de sujeto
// Un sujeto puede registrar y notificar observadores
type Subject interface {
	register(observer Observer)
	broadcast()
}

// 1.2 Item: Implementación concreta del sujeto (Subject)
// Item mantiene una lista de observadores y notifica cambios
type Item struct {
	observers []Observer
	name      string
	available bool
}

func NewItem(name string) *Item {
	return &Item{
		name: name,
	}
}

func (i *Item) register(observer Observer) {
	i.observers = append(i.observers, observer)
}

func (i *Item) MarkAsAvailable() {
	fmt.Printf("🔔 El artículo '%s' ahora está disponible\n", i.name)
	i.available = true
	i.broadcast()
}

func (i *Item) broadcast() {
	for _, observer := range i.observers {
		observer.update(i.name)
	}
}

// 2. Observer

// 2.1 Observer: Definición de la interfaz de observador
type Observer interface {
	getId() string
	update(string)
}

// 2.2 EmailClient: Implementación concreta del observador (Observer)
// EmailClient representa un cliente que recibe notificaciones por correo electrónico
type EmailClient struct {
	id    string
	email string
}

func NewEmailClient(id, email string) *EmailClient {
	return &EmailClient{
		id:    id,
		email: email,
	}
}

func (e *EmailClient) getId() string {
	return e.id
}

func (e *EmailClient) update(itemName string) {
	fmt.Printf("📧 Notificación para %s: El artículo '%s' está disponible\n", e.email, itemName)
}

// 2.3 PushClient: Otro tipo de observador que recibe notificaciones push
// PushClient representa un cliente que recibe notificaciones push
type PushClient struct {
	id     string
	device string
}

func NewPushClient(id, device string) *PushClient {
	return &PushClient{
		id:     id,
		device: device,
	}
}

func (p *PushClient) getId() string {
	return p.id
}

func (p *PushClient) update(itemName string) {
	fmt.Printf("📲 Notificación push para %s: El artículo '%s' está disponible\n", p.device, itemName)
}

// 3. Demostración
func main() {
	tarjetaGrafica := NewItem("Tarjeta Gráfica RTX 4090")
	monitorSamsung := NewItem("Monitor Samsung 4K")

	// Crear observadores (clientes) que desean recibir notificaciones
	cliente1 := NewEmailClient("1", "cliente1@example.com")
	cliente2 := NewEmailClient("2", "cliente2@example.com")
	cliente3 := NewPushClient("3", "iPhone de Cliente3")
	cliente4 := NewPushClient("4", "Android de Cliente4")

	// Registrar observadores en el sujeto (artículo)
	tarjetaGrafica.register(cliente1)
	tarjetaGrafica.register(cliente2)
	tarjetaGrafica.register(cliente3)
	tarjetaGrafica.register(cliente4)

	monitorSamsung.register(cliente1)
	monitorSamsung.register(cliente4)

	// Simular que los artículos se vuelven disponibles
	tarjetaGrafica.MarkAsAvailable()
	monitorSamsung.MarkAsAvailable()
}
