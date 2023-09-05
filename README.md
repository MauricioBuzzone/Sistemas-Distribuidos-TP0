# TP0: Docker + Comunicaciones + Concurrencia

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