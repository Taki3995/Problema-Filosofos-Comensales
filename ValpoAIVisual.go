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

const N = 5
const META_ENTRENAMIENTO = 5

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
	Recursos [N]int    // Almacena qué HiloID tiene el recurso, -1 si está libre
	Hilos    [N]string // Almacena qué está haciendo cada hilo
}

var (
	estadoRecursos [N]int
	estadoHilos    [N]string
	historial      []EstadoInstantaneo
	historialMu    sync.Mutex
)

// registrarEstado guarda una fotografía del momento exacto del clúster
func registrarEstado(hiloID int, accion string, recursoID int, estadoHilo string, mensaje string) {
	historialMu.Lock()
	defer historialMu.Unlock()

	// Actualizamos el estado de la posesión de hardware
	if accion == "TOMAR" {
		estadoRecursos[recursoID] = hiloID
	} else if accion == "SOLTAR" {
		estadoRecursos[recursoID] = -1
	}

	// Actualizamos la actividad del hilo
	if estadoHilo != "" {
		estadoHilos[hiloID] = estadoHilo
	}

	// Tomamos la "fotografía"
	var snap EstadoInstantaneo
	snap.Mensaje = fmt.Sprintf("[Hilo %d] %s", hiloID, mensaje)
	copy(snap.Recursos[:], estadoRecursos[:])
	copy(snap.Hilos[:], estadoHilos[:])

	historial = append(historial, snap)
}

func (h *Hilo) ejecutar(wg *sync.WaitGroup) {
	defer wg.Done()

	for h.ciclosCompletados < META_ENTRENAMIENTO {
		// 1. FASE DE PREPARACIÓN
		registrarEstado(h.id, "", -1, "Preparación (Limpiando memoria)", "Inicia fase de preparación")
		time.Sleep(time.Duration(rand.Intn(50)+10) * time.Millisecond)

		// 2. FASE DE ADQUISICIÓN
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
		registrarEstado(h.id, "TOMAR", primerRecurso.id, "", fmt.Sprintf("Aseguró recurso %d: %s", primerRecurso.id, nombrePrimerRecurso))

		segundoRecurso.Lock()
		registrarEstado(h.id, "TOMAR", segundoRecurso.id, "Cómputo (Entrenando)", fmt.Sprintf("Aseguró recurso %d: %s. ¡Ambos obtenidos!", segundoRecurso.id, nombreSegundoRecurso))

		// 3. FASE DE CÓMPUTO
		time.Sleep(time.Duration(rand.Intn(80)+20) * time.Millisecond)

		// 4. FASE DE LIBERACIÓN
		segundoRecurso.Unlock()
		registrarEstado(h.id, "SOLTAR", segundoRecurso.id, "", "Soltando segundo recurso")

		primerRecurso.Unlock()
		registrarEstado(h.id, "SOLTAR", primerRecurso.id, "Liberación completada", "Soltando primer recurso")

		h.ciclosCompletados++
	}
	registrarEstado(h.id, "", -1, "FINALIZADO", "Ha terminado sus 5 entrenamientos")
}

// reproducirHistorial dibuja la interfaz paso a paso
func reproducirHistorial() {
	reader := bufio.NewReader(os.Stdin)

	for i, estado := range historial {
		// Limpiar consola usando saltos de línea largos para ser seguro en todo sistema
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

	fmt.Println("Calculando simulación de ValpoIA Labs en segundo plano... Por favor espera.")

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
	
	// Una vez terminan todos los cálculos concurrentes, activamos la vista
	reproducirHistorial()
}