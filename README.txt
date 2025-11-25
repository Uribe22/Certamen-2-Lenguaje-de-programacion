Certamen II: Lenguajes de Programación
Link Video explicativo:
https://drive.google.com/drive/folders/15ZJuWG-CgZ2R9qiJSNmYaRM3XijXu3QP?usp=sharing
Link Gitub: 
https://github.com/Uribe22/Certamen-2-Lenguaje-de-programacion.git

1. Descripción General

Este proyecto simula un entorno multi-hilo usando Go y un protocolo de sincronización optimista.
El sistema está compuesto por un Scheduler, varios Workers y un Logger.
La simulación maneja generación de eventos externos, creación de eventos internos, checkpoints, violaciones de causalidad (stragglers) y rollbacks, con registro completo en un archivo CSV.

2. Arquitectura del Sistema

2.1 Scheduler
Genera eventos externos con timestamps crecientes y los distribuye a los Workers. Mantiene su propio LVT.

2.2 Workers
Cada Worker mantiene su LVT, historial de eventos externos y checkpoints.
- Recibe eventos externos, valida causalidad, actualiza el LVT, genera eventos internos y crea checkpoints.
- Si hay un straggler, se ejecuta un rollback, restaurando el estado del checkpoint y reprocesando los eventos.

2.3 Logger
Registra todos los eventos en un archivo `logs.csv` de forma secuencial mediante un canal dedicado.

3. Tipos de Eventos Registrados

- worker_creado
- scheduler_envia_evento
- evento_externo_recibido
- checkpoint
- evento_interno
- inicio_rollback
- fin_rollback

4. Ejecución del Programa

El programa se ejecuta con:

go run main.go <cantidad_workers> <cantidad_eventos> <delay_eventos_max> [<seed>]

Parámetros:
- <cantidad_workers>: Número de Workers
- <cantidad_eventos>: Eventos externos a generar
- <delay_eventos_max>: Incremento máximo para timestamps
- <seed>: Opcional, para reproducibilidad

Ejemplo:
go run main.go 4 50 5

Salida:
- Archivo logs.csv con los eventos registrados
- Registro de ejecución por consola

5. Visualización Offline

El script en Python (visualizacion.py) genera un gráfico que muestra:
- La evolución del LVT de cada Worker.
- Eventos externos, internos y checkpoints.
- Inicio y fin de rollbacks.
Este gráfico ayuda a visualizar el progreso de cada Worker y cómo se manejan los eventos y rollbacks.

6. Análisis de Escalabilidad

Se probó la simulación con diferentes cantidades de Workers (1, 2, 4, 8, 16) y se calculó el speedup comparado con 1 Worker.

Workers	Tiempo Promedio Speedup
1	27.198	         1.00
2	26.6365	         1.02
4	26.1152	         1.04
8	24.2869	         1.12
16	23.3150	         1.17

Resultados:
- Tabla de tiempos
- Gráfico de speedup

7. Consideraciones de Diseño

El sistema se basa en:
- Comunicación por canales para evitar condiciones de carrera
- Logger único para escritura secuencial
- Checkpoints eficientes que almacenan solo lo necesario
- Eventos internos regenerables
- Implementación del enfoque optimista con rollback

8. Archivos Incluidos

main.go            Implementación completa del sistema
documentación.pdf  Análisis de escalabilidad y gráficos
logs.csv           Registro generado por el sistema
visualizacion.py   Script para graficar la ejecución
README.md          Documento actual