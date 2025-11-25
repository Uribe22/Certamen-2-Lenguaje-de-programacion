package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type Evento struct {
	id           int64
	timestamp    int64
	cantProcesos int64
}

type Checkpoint struct {
	lvt       int64
	historial []Evento
}

func LogCSV(worker int, tipo string, idEvento interface{}, lvt int64, desc string) string {
	fecha := time.Now().Format(time.RFC3339)
	return fmt.Sprintf("%s,%d,%s,%v,%d,%s", fecha, worker, tipo, idEvento, lvt, desc)
}

func EjecutarTarea(id int, evento Evento, lvt *int64, logs chan<- string, random *rand.Rand) {
	// Randomizar cantidad de procesos si no está definida
	if evento.cantProcesos == 0 {
		evento.cantProcesos = int64(random.Intn(5) + 1)
	}

	// Ejecutar procesos
	for i := int64(0); i < evento.cantProcesos; i++ {
		*lvt += 2
		idEventoInterno := fmt.Sprintf("%d.%d", evento.id, i+1)

		logs <- LogCSV(
			id,
			"evento_interno",
			idEventoInterno,
			*lvt,
			"fin_evento_interno",
		)
	}
}

func GuardarCheckpoint(id int, lvt int64, historial []Evento, checkpoints *[]Checkpoint, logs chan<- string) {
	*checkpoints = append(*checkpoints, Checkpoint{lvt: lvt, historial: append([]Evento{}, (historial)...)})

	logs <- LogCSV(
		id,
		"checkpoint",
		"-",
		lvt,
		"checkpoint_guardado",
	)

}

func Logger(logs <-chan string, done *sync.WaitGroup) {
	defer done.Done()

	file, err := os.Create("logs.csv")
	if err != nil {
		fmt.Println("Error creando logs.csv:", err)
		return
	}
	defer file.Close()

	file.WriteString("fecha,worker,tipo,id_evento,lvt,descripcion\n")

	for log := range logs {
		fmt.Println(log)
		file.WriteString(log + "\n")
	}
}

func Worker(id int, eventos <-chan Evento, logs chan<- string, seed int64, done *sync.WaitGroup) {
	defer done.Done()
	random := rand.New(rand.NewSource(seed + int64(id)))

	lvt := int64(0)
	historial := []Evento{}
	checkpoints := []Checkpoint{}

	logs <- LogCSV(
		id,
		"worker_creado",
		"-",
		lvt,
		"worker_creado",
	)

	for eventoActivo := range eventos {
		logs <- LogCSV(
			id,
			"evento_externo_recibido",
			eventoActivo.id,
			eventoActivo.timestamp,
			"recibido_evento_externo",
		)

		if eventoActivo.timestamp > lvt { // Respeta causalidad
			lvt = eventoActivo.timestamp
			historial = append(historial, eventoActivo)

			GuardarCheckpoint(id, lvt, historial, &checkpoints, logs)
			EjecutarTarea(id, Evento{eventoActivo.id, lvt, 0}, &lvt, logs, random)

		} else { // Violación de causalidad
			logs <- LogCSV(
				id,
				"inicio_rollback",
				eventoActivo.id,
				lvt,
				"violacion_de_causalidad",
			)

			// Ejecutar rollback
			var checkpoint_rollback Checkpoint
			encontrado := false
			for i := len(checkpoints) - 1; i >= 0; i-- {
				if checkpoints[i].lvt <= eventoActivo.timestamp {
					checkpoint_rollback = checkpoints[i]
					encontrado = true

					// Eliminar checkpoints posteriores al seleccionado (erróneos)
					checkpoints = checkpoints[:i+1]
					break
				}
			}

			// Verificar checkpoint y restaurar estado, conservando eventos posteriores para reprocesar
			var eventosPosteriores []Evento
			if encontrado {
				eventosPosteriores = historial[len(checkpoint_rollback.historial):]
				lvt = checkpoint_rollback.lvt
				historial = append([]Evento{}, checkpoint_rollback.historial...)

				logs <- LogCSV(
					id,
					"fin_rollback",
					"-",
					lvt,
					"estado_restaurado",
				)

			} else {
				eventosPosteriores = historial
				lvt = 0
				historial = []Evento{}

				logs <- LogCSV(
					id,
					"rollback_estado_inicial",
					"-",
					lvt,
					"sin_checkpoint_valido",
				)

			}

			// Reprocesar eventos posteriores
			eventosReprocesar := append([]Evento{eventoActivo}, eventosPosteriores...)
			for _, reprocesar := range eventosReprocesar {
				EjecutarTarea(id, reprocesar, &lvt, logs, random)
			}
		}
	}
}

func Scheduler(cantidadWorkers int, cantidadEventos int, delayEventosMax int, seed int64, done *sync.WaitGroup) {
	defer done.Done()

	logs := make(chan string, 100)
	var wgWorkers sync.WaitGroup
	var wgLogger sync.WaitGroup

	// Inicializar logger
	wgLogger.Add(1)
	go Logger(logs, &wgLogger)

	// Iniciar workers
	canalesWorkers := make([]chan Evento, cantidadWorkers)
	wgWorkers.Add(cantidadWorkers)

	for i := 0; i < cantidadWorkers; i++ {
		canalesWorkers[i] = make(chan Evento, 100)
		go Worker(i, canalesWorkers[i], logs, seed, &wgWorkers)
	}

	random := rand.New(rand.NewSource(seed))

	// Asignación de eventos a workers
	lvt := int64(0)
	for i := 0; i < cantidadEventos; i++ {
		incremento := random.Intn(delayEventosMax) + 1
		lvt += int64(incremento)

		evento := Evento{
			id:        int64(i),
			timestamp: lvt,
		}

		workerSeleccionado := i % cantidadWorkers
		logs <- LogCSV(
			-1,
			"scheduler_envia_evento",
			evento.id,
			evento.timestamp,
			"envio_evento_externo",
		)
		canalesWorkers[workerSeleccionado] <- evento
	}

	// Cerrar canales de workers y esperar a que terminen
	for i := 0; i < cantidadWorkers; i++ {
		close(canalesWorkers[i])
	}
	wgWorkers.Wait()

	// Cerrar canal de logger y esperar a que termine
	close(logs)
	wgLogger.Wait()
}

func main() {
	// Comprobar argumentos
	if len(os.Args) != 4 && len(os.Args) != 5 {
		fmt.Println("Uso: go run main.go <cantidad_workers> <cantidad_eventos> <delay_eventos_max> [<seed>]")
		return
	}

	// Parsear argumentos
	cantidadWorkers, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil || cantidadWorkers <= 0 {
		fmt.Println("<cantidad_workers> debe ser un entero mayor a 0")
		return
	}

	cantidadEventos, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil || cantidadEventos <= 0 {
		fmt.Println("<cantidad_eventos> debe ser un entero mayor a 0")
		return
	}

	delayMax, err := strconv.ParseInt(os.Args[3], 10, 64)
	if err != nil || delayMax <= 0 {
		fmt.Println("<delay_eventos_max> debe ser un entero mayor a 0")
		return
	}

	// Inicializar generador de números aleatorios
	var seed int64
	if len(os.Args) == 5 {
		seed, err = strconv.ParseInt(os.Args[4], 10, 64)
		if err != nil {
			fmt.Println("<seed> debe ser un entero")
			return
		}

	} else {
		seed = time.Now().UnixNano()
	}

	var wg sync.WaitGroup
	wg.Add(1)

	//Esperar a que termine el scheduler
	go Scheduler(int(cantidadWorkers), int(cantidadEventos), int(delayMax), seed, &wg)
	wg.Wait()
}
