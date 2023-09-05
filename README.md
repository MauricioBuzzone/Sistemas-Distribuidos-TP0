# TP0: Docker + Comunicaciones + Concurrencia

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
