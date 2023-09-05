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