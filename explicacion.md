# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°1:
Modificar la definición del DockerCompose para agregar un nuevo cliente al proyecto.

#### Solución:
Se modifica el archivo `docker-compose-dev.yaml`, agregando un nuevo servicio llamado client2.

```
  client2:
    container_name: client2
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=2
      - CLI_LOG_LEVEL=DEBUG
    networks:
      - testing_net
    depends_on:
      - server
```

### Ejercicio N°1.1:
Definir un script (en el lenguaje deseado) que permita crear una definición de DockerCompose con una cantidad configurable de clientes.

#### Solución:
Se creó un script de bash que permite configurar la cantidad de clientes. Este sobrescribe el archivo de configuración `docker-compose-dev.yaml`.

```
./create-docker-compose-N.sh <N_CLIENTS>
```

Para ejecutarlo, primero se le debe dar permisos.
```
chmod a+x create-docker-compose-N.sh
```

### Ejercicio N°2:
Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera un nuevo build de las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida afuera de la imagen (hint: `docker volumes`).


#### Solución:
Se modifica el DockerCompose, agregando los archivos de configuración como volumenes para el servidor y el cliente respectivamente.


```
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - LOGGING_LEVEL=DEBUG
    volumes:
      - ./server/config.ini:/config.ini
    networks:
      - testing_net
```

### Ejercicio N°3:
Crear un script que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un EchoServer, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado. Netcat no debe ser instalado en la máquina _host_ y no se puede exponer puertos del servidor para realizar la comunicación (hint: `docker network`).


A partir de una imagen de alphine se descarga netcat y se copia el script `netcat.sh`, para luego ejecutarlo.

```Dockerfile
FROM alpine:latest
RUN apk update && apk add netcat-openbsd
COPY netcat/netcat.sh /
RUN chmod +x netcat.sh
ENTRYPOINT ["/netcat.sh"]
```

El script envía un mensaje al servidor con la IP especificada en el archivo de configuración `config.txt` y verifica que el mensaje devuelto por el servidor es igual al que mandó.

```sh
#!/bin/sh

response=$(echo "Hello" | nc $SERVER_IP $SERVER_PORT)

if [ "$response" == "Hello" ]; then
    echo "Server responded by repeating the message"
else
    echo "Server not responding"
fi
```

Para ejecutarlo se agregó un comando al Makefile que levanta la imagen y corre el script.

```Makefile
docker-netcat:
	docker build -f ./netcat/Dockerfile -t netcat-image .
	docker run --rm --network tp0_testing_net --env-file ./netcat/config.txt --name netcat-container netcat-image 
```


### Ejercicio N°4:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).


#### Solución
Para conseguir que tanto el cliente como el servidor tenga una cierre _graceful_, ambos escuchar por un canal un singal SIGTERM.

#### Servidor

En el caso del servidor, cuando se reciba un SIGTEM: Pondrá el booleano `server_on` en false para que no se acepten más conexiones y cierra el `server_socket`. 

```py

def __handle_signal(self, signum, frame):
    """
    Close server socket graceful
    """
    logging.info(f'action: stop_server | result: in_progress | singal {signum}')
    try:
        self._server_on = False
        self._server_socket.shutdown(socket.SHUT_RDWR)
        logging.info(f'action: shutdown_socket | result: success')
        self._server_socket.close()
        logging.info(f'action: release_socket | result: success')

    except OSError as e:  
        logging.error(f'action: stop_server | result: fail | error: {e}')


```

#### Cliente

Se agrega un chanel por el cual se recibe el SIGTERM. Cuando la señal es recibida, se libera el channel y el client socket.


```go
case <-signalChan:
    log.Infof("action: shutdown_detected | result: success | client_id: %v",
		c.config.ID,
	)
	close(signalChan)
	log.Infof("action: release_channel | result: success | client_id: %v",
		c.config.ID,
	)
	c.conn.Close()
	log.Infof("action: release_socket | result: success")			

```


## Parte 2: Repaso de Comunicaciones

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado **Lotería Nacional**. Para la resolución de las mismas deberá utilizarse como base al código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.

### Ejercicio N°5:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Cliente
Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente.

Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Servidor
Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno.
Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Comunicación:
Se deberá implementar un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:
* Definición de un protocolo para el envío de los mensajes.
* Serialización de los datos.
* Correcta separación de responsabilidades entre modelo de dominio y capa de comunicación.
* Correcto empleo de sockets, incluyendo manejo de errores y evitando los fenómenos conocidos como [_short read y short write_](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/).

#### Solución 

#### Cliente

Para que reciba los datos de la apuesta se modifico el Dockfile del usuario.

```Dockfile
FROM busybox:latest
COPY --from=builder /build/bin/client /client
COPY ./client/config.yaml /config.yaml
ENV NOMBRE="Santiago Lionel"
ENV APELLIDO="Lorca"
ENV DOCUMENTO="30904465"
ENV NACIMIENTO="1999-03-17"
ENV NUMERO="7574"
ENTRYPOINT ["/bin/sh"]
```

El cliente inicialmente crea una conexión con el servidor, crea un apuesta con las variables de entorno, la serializa, la envía y espera que el servidor le responda afirmativamente.

```go
c.createClientSocket()

bet := Bet{
ID:            c.config.ID,
FirstName:     os.Getenv("NOMBRE"),
LastName:	   os.Getenv("APELLIDO"),
Document:	   os.Getenv("DOCUMENTO"),
Birthdate:	   os.Getenv("NACIMIENTO"),
Number:        os.Getenv("NUMERO"),
}
data := serializeBet(bet)
sendBet(c.conn, data)

log.Infof("action: esperando_confirmacion | result: in_progress")
msg, err := readMessage(c.conn)
log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
    msg[1],
    msg[2],
)
```

#### Servidor

El servidor recibe el mensaje y llama a `__handle_message(...)` luego de reconstruir la apuesta envía un mensaje de confirmación al cliente.

### Protocolo


El protocolo de comunicación consta de 3 partes: El Tamaño del mensaje, el tipo de mensaje y el payload o carga útil.

+ Tamaño del mensaje: Son los primeros 4 bytes que se envían e indican el tamaño del mensaje (sin contar estos mismos 4 bytes y el byte de tipo)

+ Tipo de mensaje: Es el quinto byte e indica qué tipo de mensaje se está enviado.
```
BET_TYPE = 'B'   (Bet: Mensaje de envío de apuesta)
OK_TYPE = '0'    (Ok:  Mensaje de confirmación de la llegada de la apuesta)
ERR_TYPE = 'E'   (Bet: Mensaje de aviso de algún error con la apuesta)
```

+ Payload: Esta sección es en la que efectivamente se envían los campos de la apuesta. Cada campo es precedido por su longitud, para que el receptor del mensaje pueda saber cuando bytes debe leer antes de terminar de leer dicho campo.

A continuación se muestra la estructura de un mensaje del protocolo.

```
             |   field1     |   field2     |
| size | type| size | data1 | size | data2 |  
|4bytes|1byte|4bytes|Nbytes |4bytes|Mbytes |
```

En esta primera versión el cliente envía un mensaje de tipo `BET` con la información de la apuesta. El servidor la recibe, la parsea a un objeto bet y envía un mensaje de tipo `OK` con el documento y el número de la apuesta para confirmarle al cliente que el mensaje fue recibido. 



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


### Ejercicio N°7:
Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo, no podrá responder consultas por la lista de ganadores.
Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.


#### Solución:

Luego de enviar todas las apuestas, el cliente crea una nueva conexión y consulta por los ganadores de su agencia. Para ello envía un nuevo tipo de mensaje al servidor.

El mensaje CHECK_WIN que envía el cliente, solicita al servidor que le envíe los documentos de los ganadores de su agencia. Ante este mensaje el servidor puede encontrarse en dos estados.

1. El servidor no tiene todas las apuestas de las 5 agencias: De modo que aun no puede hacer el sorteo. En este caso, envía un mensaje al cliente del tipo  CHECK_WIN, sin fecha. De esta forma, el cliente interpreta que el servidor aun no puedo hacer la votación y esperará un tiempo hasta volver a conectarse y consultar. El tiempo que espera entre consultas se duplica por cada consulta "fallida".


2. El servidor ya tiene todas las apuestas y procede a hacer el sorteo. En este caso el servidor llama a función `get_winners(id)` con el ID de la agencia que hizo la consulta (Este es enviado en el mensaje de consulta del cliente) y obtiene todos los DNI ganadores de la respectiva agencia. Los DNI ganadores son enviados al cliente y se finaliza la conexión.


### Ejercicio N°8:
Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo. En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

En caso de que el alumno implemente el servidor Python utilizando multithreading, deberán tenerse en cuenta las limitaciones propias del lenguaje.


#### Solución:

Para que el servidor permita procesar mensajes en paralelo, se utilizaron threads que manejen los clientes a medida que se van conectando. Para cada nueva conexión se crea un nuevo thread que ejecuta la función `handle_client_connection`, la cual contiene toda la lógica de la comunicación.

A la hora de querer cerrar el servidor, inmediatamente no se aceptarán más conexiones, pero los clientes que estaban siendo atendidos no serán cerrados forzosamente.

Para optimizar la utilización de recursos se aprovecha la oportunidad de una nueva conexión para reciclar los thread que ya finalizaron y no que hacer join a todos cuando se quiere cerrar el servidor.

Existen dos elementos  compartidos que necesitan sincronización en el sistema:

1. El acceso al archivo de apuestas del servidor con `load_bets` y `store_bets`, el cual es protegido por un lock.

2. El acceso al registrador de las agencias  (`agencyRegister`) que se encarga de almacenar la información de cuáles son las agencias que ya terminaron de enviar sus apuestas. El acceso a este también se hace por medio de un lock.