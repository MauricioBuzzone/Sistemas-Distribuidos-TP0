# TP0: Docker + Comunicaciones + Concurrencia

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