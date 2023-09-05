# TP0: Docker + Comunicaciones + Concurrencia

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