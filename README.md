# Certamen II: Lenguajes de Programación

## Link Video Explicativo
[Google Drive Folder](https://drive.google.com/drive/folders/15ZJuWG-CgZ2R9qiJSNmYaRM3XijXu3QP?usp=sharing)

## Link GitHub
[Certamen-2-Lenguaje-de-programacion](https://github.com/Uribe22/Certamen-2-Lenguaje-de-programacion.git)

## 1. Descripción General

Este proyecto simula un entorno multi-hilo usando Go y un protocolo de sincronización optimista.
El sistema está compuesto por un Scheduler, varios Workers y un Logger.
La simulación maneja generación de eventos externos, creación de eventos internos, checkpoints, violaciones de causalidad (stragglers) y rollbacks, con registro completo en un archivo CSV.

## 2. Arquitectura del Sistema

### 2.1 Scheduler
Genera eventos externos con timestamps crecientes y los distribuye a los Workers. Mantiene su propio LVT.

### 2.2 Workers
Cada Worker mantiene su LVT, historial de eventos externos y checkpoints.
- Recibe eventos externos, valida causalidad, actualiza el LVT, genera eventos internos y crea checkpoints.
- Si hay un straggler, se ejecuta un rollback, restaurando el estado del checkpoint y reprocesando los eventos.

### 2.3 Logger
Registra todos los eventos en un archivo `logs.csv` de forma secuencial mediante un canal dedicado.

## 3. Tipos de Eventos Registrados

- worker_creado
- scheduler_envia_evento
- evento_externo_recibido
- checkpoint
- evento_interno
- inicio_rollback
- fin_rollback

## 4. Ejecución del Programa

El programa se ejecuta con:

```bash
go run main.go <cantidad_workers> <cantidad_eventos> <delay_eventos_max> [<seed>]
```

### Parámetros:
- `<cantidad_workers>`: Número de Workers
- `<cantidad_eventos>`: Eventos externos a generar
- `<delay_eventos_max>`: Incremento máximo para timestamps
- `<seed>`: Opcional, para reproducibilidad

### Ejemplo:
```bash
go run main.go 4 50 5
```