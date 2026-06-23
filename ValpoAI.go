package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

// Constantes de configuración
const N = 5
const META_ENTRENAMIENTO = 5

// --- ESTRUCTURAS PRINCIPALES ---
type Recurso struct {
	sync.Mutex
	id int
}

type Hilo struct {
	id                int
	izquierdo         *Recurso
	derecho           *Recurso
	ciclosCompletados int
}

// --- ESTRUCTURAS DE VISUALIZACIÓN ---
type EstadoInstantaneo struct {
	Mensaje  string
	Recursos [N]int
	Hilos    [N]string
}

// Variables globales para la "cinta de grabación"
var (
	estadoRecursos [N]int
	estadoHilos    [N]string
	historial      []EstadoInstantaneo
	historialMu    sync.Mutex
)

// --- FUNCIONES ---

// logMensaje imprime los eventos en consola en tiempo real
func logMensaje(id int, mensaje string) {
	tiempoActual := time.Now().Format("15:04:05.000")
	fmt.Printf("[%s] Hilo %d: %s\n", tiempoActual, id, mensaje)
}

// registrarEstado guarda una fotografía del clúster de forma silenciosa
func registrarEstado(hiloID int, accion string, recursoID int, estadoHilo string, mensaje string) {
	historialMu.Lock()
	defer historialMu.Unlock()

	if accion == "TOMAR" {
		estadoRecursos[recursoID] = hiloID
	} else if accion == "SOLTAR" {
		estadoRecursos[recursoID] = -1
	}

	if estadoHilo != "" {
		estadoHilos[hiloID] = estadoHilo
	}

	var snap EstadoInstantaneo
	snap.Mensaje = fmt.Sprintf("[Hilo %d] %s", hiloID, mensaje)
	copy(snap.Recursos[:], estadoRecursos[:])
	copy(snap.Hilos[:], estadoHilos[:])

	historial = append(historial, snap)
}

// ejecutar contiene la lógica de las 4 fases
func (h *Hilo) ejecutar(wg *sync.WaitGroup) {
	defer wg.Done()

	for h.ciclosCompletados < META_ENTRENAMIENTO {
		// 1. FASE DE PREPARACIÓN
		logMensaje(h.id, "Fase de Preparación: Limpiando memoria...")
		registrarEstado(h.id, "", -1, "Preparación (Limpiando memoria)", "Inicia fase de preparación")
		time.Sleep(time.Duration(rand.Intn(500)+100) * time.Millisecond)

		// 2. FASE DE ADQUISICIÓN
		logMensaje(h.id, "Fase de Adquisición: Intentando tomar recursos (Izquierda: Gráfica, Derecha: Almacenamiento)...")
		registrarEstado(h.id, "", -1, "Adquisición (Esperando recursos)", "Intentando adquirir recursos")

		var nombrePrimerRecurso, nombreSegundoRecurso string
		primerRecurso := h.izquierdo
		segundoRecurso := h.derecho

		if h.id%2 == 0 {
			nombrePrimerRecurso = "Gráfica (Izquierdo)"
			nombreSegundoRecurso = "Almacenamiento (Derecho)"
		} else {
			primerRecurso = h.derecho
			segundoRecurso = h.izquierdo
			nombrePrimerRecurso = "Almacenamiento (Derecho)"
			nombreSegundoRecurso = "Gráfica (Izquierdo)"
		}

		primerRecurso.Lock()
		logMensaje(h.id, fmt.Sprintf("Aseguró recurso %d: %s.", primerRecurso.id, nombrePrimerRecurso))
		registrarEstado(h.id, "TOMAR", primerRecurso.id, "", fmt.Sprintf("Aseguró recurso %d: %s", primerRecurso.id, nombrePrimerRecurso))

		segundoRecurso.Lock()
		logMensaje(h.id, fmt.Sprintf("Aseguró recurso %d: %s. ¡Ambos recursos obtenidos!", segundoRecurso.id, nombreSegundoRecurso))
		registrarEstado(h.id, "TOMAR", segundoRecurso.id, "Cómputo (Entrenando)", fmt.Sprintf("Aseguró recurso %d: %s. ¡Ambos obtenidos!", segundoRecurso.id, nombreSegundoRecurso))

		// 3. FASE DE CÓMPUTO
		logMensaje(h.id, "Fase de Cómputo: Entrenando modelo de IA...")
		time.Sleep(time.Duration(rand.Intn(800)+200) * time.Millisecond)

		// 4. FASE DE LIBERACIÓN
		logMensaje(h.id, "Fase de Liberación: Soltando recursos...")

		segundoRecurso.Unlock()
		registrarEstado(h.id, "SOLTAR", segundoRecurso.id, "", "Soltando segundo recurso")

		primerRecurso.Unlock()
		registrarEstado(h.id, "SOLTAR", primerRecurso.id, "Liberación completada", "Soltando primer recurso")

		h.ciclosCompletados++
		logMensaje(h.id, fmt.Sprintf(">>> Completó %d de %d entrenamientos. <<<", h.ciclosCompletados, META_ENTRENAMIENTO))
	}
	logMensaje(h.id, "--- HA FINALIZADO SU TRABAJO ---")
	registrarEstado(h.id, "", -1, "FINALIZADO", "Ha terminado sus 5 entrenamientos")
}

// reproducirHistorial dibuja la interfaz paso a paso al terminar
func reproducirHistorial() {
	reader := bufio.NewReader(os.Stdin)

	for i, estado := range historial {
		fmt.Print("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
		fmt.Println("=====================================================")
		fmt.Printf(" REPRODUCCIÓN DEL CLÚSTER - PASO %d / %d\n", i+1, len(historial))
		fmt.Println("=====================================================")
		fmt.Printf(" >>> ACCIÓN: %s\n\n", estado.Mensaje)

		fmt.Println(" [ MAPA FÍSICO DE RECURSOS ]")
		for r := 0; r < N; r++ {
			dueño := "LIBRE"
			if estado.Recursos[r] != -1 {
				dueño = fmt.Sprintf("Ocupado por HILO %d", estado.Recursos[r])
			}
			fmt.Printf("  - Recurso %d : %s\n", r, dueño)
		}

		fmt.Println("\n [ ESTADO LÓGICO DE LOS HILOS ]")
		for h := 0; h < N; h++ {
			fmt.Printf("  - Hilo %d : %s\n", h, estado.Hilos[h])
		}
		fmt.Println("=====================================================")

		fmt.Print("Presiona ENTER para avanzar (o escribe 'q' y ENTER para salir)...")
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(input) == "q" {
			break
		}
	}
	fmt.Println("\n=== FIN DE LA REPRODUCCIÓN VISUAL ===")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Inicializar estado de arreglos para visualización
	for i := 0; i < N; i++ {
		estadoRecursos[i] = -1 // -1 significa libre
		estadoHilos[i] = "Preparación"
	}

	fmt.Println("=== Iniciando Clúster ValpoIA Labs ===")
	fmt.Println("Los hilos procesarán y emitirán logs en tiempo real...")
	fmt.Println("-----------------------------------------------------")

	recursos := make([]*Recurso, N)
	for i := 0; i < N; i++ {
		recursos[i] = &Recurso{id: i}
	}

	hilos := make([]*Hilo, N)
	for i := 0; i < N; i++ {
		hilos[i] = &Hilo{
			id:        i,
			izquierdo: recursos[i],
			derecho:   recursos[(i+1)%N],
		}
	}

	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go hilos[i].ejecutar(&wg)
	}

	wg.Wait()

	fmt.Println("-----------------------------------------------------")
	fmt.Println("=== Todos los hilos han terminado exitosamente. ===")

	// Preguntar al usuario si desea ver la reproducción visual
	fmt.Println("\n¿Deseas ver la reproducción visual de los eventos paso a paso?")
	fmt.Print("Presiona ENTER para iniciar (o escribe 'q' y ENTER para salir)... ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(input) != "q" {
		reproducirHistorial()
	}
}
