import pandas as pd
import matplotlib.pyplot as plt

# Leer logs
df = pd.read_csv("logs.csv")
df = df[df["worker"] >= 0]   # solo workers válidos

# Construir orden de ejecución LOCAL por worker
df["orden_local"] = df.groupby("worker").cumcount()

plt.figure(figsize=(15,7))

# Colores por tipo de evento
coloresEventos = {
    "evento_externo_recibido": "blue",
    "evento_interno": "green",
    "checkpoint": "orange",
    "inicio_rollback": "red",
    "fin_rollback": "purple"
}

# Marcadores por tipo de evento
marcadoresEventos = {
    "evento_externo_recibido": "o",
    "evento_interno": "o",
    "checkpoint": "D",
    "inicio_rollback": "v",
    "fin_rollback": "^"
}

zOrderEventos = {
    "evento_externo_recibido": 3,
    "evento_interno": 2,
    "checkpoint": 4,
    "inicio_rollback": 6,
    "fin_rollback": 5
}

# Paleta de colores para workers
coloresWorkers = ["#1f77b4", "#ff7f0e", "#2ca02c",
                 "#d62728", "#9467bd", "#8c564b",
                 "#e377c2", "#7f7f7f", "#bcbd22", "#17becf"]

workers = sorted(df["worker"].unique())

for i, worker in enumerate(workers):
    subset = df[df["worker"] == worker]
    colorWorker = coloresWorkers[i % len(coloresWorkers)]

    # Dibujar línea punteada del worker
    plt.plot(
        subset["orden_local"],
        subset["lvt"],
        color=colorWorker,
        linestyle="--",
        linewidth=1.2,
        alpha=0.8,
        zorder=1
    )

    # Dibujar eventos
    for tipo in coloresEventos:
        puntos = subset[subset["tipo"] == tipo]
        if len(puntos) > 0:
            plt.scatter(
                puntos["orden_local"],
                puntos["lvt"],
                color=coloresEventos[tipo],
                marker=marcadoresEventos[tipo],
                s=50,
                zorder=zOrderEventos[tipo]
            )

    # Final del worker
    last = subset.iloc[-1]
    plt.scatter(
        last["orden_local"],
        last["lvt"],
        color=colorWorker,
        marker="^",
        s=120,
        zorder=20
    )

# Etiquetas gráfico
plt.title("Líneas de procesamiento por Worker")
plt.xlabel("Orden de Ejecución Local")
plt.ylabel("LVT")
plt.grid(True, linestyle="--", alpha=0.3)

# Leyenda de eventos
legend_items = [
    plt.Line2D([0], [0], color="black", linestyle="--", lw=1.2, label="Línea Worker"),
    plt.Line2D([0], [0], color="blue", marker="o", lw=0, markersize=7, label="Evento Externo"),
    plt.Line2D([0], [0], color="green", marker="o", lw=0, markersize=7, label="Evento Interno"),
    plt.Line2D([0], [0], color="orange", marker="D", lw=0, markersize=7, label="Checkpoint"),
    plt.Line2D([0], [0], color="red", marker="v", lw=0, markersize=7, label="Inicio Rollback"),
    plt.Line2D([0], [0], color="purple", marker="^", lw=0, markersize=7, label="Fin Rollback"),
]
plt.legend(handles=legend_items, loc="upper left")

plt.show()
