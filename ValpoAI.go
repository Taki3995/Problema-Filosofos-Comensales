package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Constantes de configuración
const N = 5
const META_ENTRENAMIENTO = 5

// Recurso representa el hardware compartido (GPU + NVMe)
type Recurso struct {
	sync.Mutex // Mutex garantiza la exclusión mutua
	id         int
}

// Hilo representa la unidad de cómputo de IA
type Hilo struct {
	id                int
	izquierdo         *Recurso
	derecho           *Recurso
	ciclosCompletados int
}

// Función auxiliar para imprimir eventos con milisegundos
func logMensaje(id int, mensaje string) {
	// Formato de hora incluyendo milisegundos (.000)
	tiempoActual := time.Now().Format("15:04:05.000")
	fmt.Printf("[%s] Hilo %d: %s\n", tiempoActual, id, mensaje)
}

// ejecutar contiene la lógica de las 4 fases solicitadas
func (h *Hilo) ejecutar(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		// 1. FASE DE PREPARACIÓN
		logMensaje(h.id, "Fase de Preparación: Limpiando memoria...")
		// Simula tiempo de procesamiento local sin usar recursos
		time.Sleep(time.Duration(rand.Intn(500)+100) * time.Millisecond)

		// 2. FASE DE ADQUISICIÓN
		logMensaje(h.id, "Fase de Adquisición: Intentando tomar recursos...")

		// Estrategia Asimétrica (Zurdos y Diestros)
		primerRecurso := h.izquierdo
		segundoRecurso := h.derecho

		if h.id%2 != 0 {
			// Si es un hilo impar, invierte el orden (primero derecha, luego izquierda)
			primerRecurso = h.derecho
			segundoRecurso = h.izquierdo
		}

		// Adquiere el primer recurso (Espera pacientemente en segundo plano si está ocupado)
		primerRecurso.Lock()
		logMensaje(h.id, fmt.Sprintf("Aseguró recurso %d.", primerRecurso.id))

		// Adquiere el segundo recurso
		segundoRecurso.Lock()
		logMensaje(h.id, fmt.Sprintf("Aseguró recurso %d. ¡Ambos recursos obtenidos!", segundoRecurso.id))

		// 3. FASE DE CÓMPUTO
		logMensaje(h.id, "Fase de Cómputo: Entrenando modelo de IA...")
		// Simula el tiempo de entrenamiento de forma aleatoria
		time.Sleep(time.Duration(rand.Intn(800)+200) * time.Millisecond)

		// 4. FASE DE LIBERACIÓN
		logMensaje(h.id, "Fase de Liberación: Soltando recursos...")
		// Suelta los recursos para los vecinos
		segundoRecurso.Unlock()
		primerRecurso.Unlock()

		h.ciclosCompletados++

		// Verificación de la meta impuesta
		if h.ciclosCompletados == META_ENTRENAMIENTO {
			logMensaje(h.id, ">>> ¡ALCANZÓ LA META DE 5 ENTRENAMIENTOS! <<<")
		}
	}
}

func main() {
	// Inicializar la semilla para los tiempos aleatorios
	rand.Seed(time.Now().UnixNano())

	fmt.Println("=== Iniciando Clúster ValpoIA Labs ===")

	// Crear los N recursos compartidos
	recursos := make([]*Recurso, N)
	for i := 0; i < N; i++ {
		recursos[i] = &Recurso{id: i}
	}

	// Crear los N hilos y asignarles sus recursos adyacentes lógicos
	hilos := make([]*Hilo, N)
	for i := 0; i < N; i++ {
		hilos[i] = &Hilo{
			id:        i,
			izquierdo: recursos[i],
			derecho:   recursos[(i+1)%N], // Topología circular
		}
	}

	// WaitGroup para mantener el programa principal corriendo
	var wg sync.WaitGroup

	// Iniciar los hilos de forma concurrente
	for i := 0; i < N; i++ {
		wg.Add(1)
		go hilos[i].ejecutar(&wg)
	}

	// Bloquea el hilo principal para que el programa no termine instantáneamente
	wg.Wait()
}
