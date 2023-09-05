# TP0: Docker + Comunicaciones + Concurrencia

### Ejercicio N°6:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento. La cantidad de apuestas dentro de cada _batch_ debe ser configurable. Realizar una implementación genérica, pero elegir un valor por defecto de modo tal que los paquetes no excedan los 8kB. El servidor, por otro lado, deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.


#### Solución:

El cliente ahora lee un archivo de apuestas que debe enviar al servidor. Este envía varias apuestas en un mismo mensaje. La cantidad de apuesta depende de un parámetro de configuración especificado en `config.yaml`. A su vez, los mensaje que se envían no pueden superar los 8kytes. De modo que si se especifica una cantidad de apuestas que supera el tamaño máximo, se enviaran tantas apuesta como entren en estos 8kbytes.

También se incorporó un nuevo tipo de mensaje al protocolo.

```
END_TYPE = 'F' (FIN: Alerta al servidor de que ya se enviaron todos los apuestas del archivo)
```


Actualmente un cliente se conecta con el servidor y envía en _batchs_ la información leída de su archivo correspondiente. Por cada _batch_ el servidor responde con un mensaje de tipo OK y la cantidad de apuestas que fueron enviadas. Una vez el cliente terminó de enviar todas las apuestas, este envía un mensaje de tipo END para indicarle al servidor que la comunicación ya finalizó y ya no hay más apuestas por cargar. Cuando el servidor recibe un mensaje de END, cierra la conexión con el cliente y se pone a escuchar nuevos clientes.