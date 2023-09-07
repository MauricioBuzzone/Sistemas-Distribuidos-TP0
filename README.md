# TP0: Docker + Comunicaciones + Concurrencia

### Ejercicio N°8:
Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo. En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

En caso de que el alumno implemente el servidor Python utilizando multithreading, deberán tenerse en cuenta las limitaciones propias del lenguaje.


#### Solución:

Para que el servidor permita procesar mensajes en paralelo, se utilizaron threads que manejen los clientes a medida que se van conectando. Para cada nueva conexión se crea un nuevo thread que ejecuta la función `handle_client_connection`, la cual contiene toda la lógica de la comunicación.

A la hora de querer cerrar el servidor, inmediatamente no se aceptarán más conexiones, pero los clientes que estaban siendo atendidos no serán cerrados forzosamente.

Existen dos elementos  compartidos que necesitan sincronización en el sistema:

1. El acceso al archivo de apuestas del servidor con `load_bets` y `store_bets`, el cual es protegido por un lock.

2. El acceso al registrador de las agencias  (`agencyRegister`) que se encarga de almacenar la información de cuáles son las agencias que ya terminaron de enviar sus apuestas. El acceso a este también se hace por medio de un lock.

El caso más complejo donde se puede ver el manejo de los recursos compartidos es en la función `handle_winners`, donde es necesario acceder a ambos recursos. En este caso, primero se adquieren todos los locks juntos, y luego se van liberando a medida que ya no son necesarios para continuar.

```py
def handle_winners(client_socket, agency_id, agency_register, agency_register_lock, bets_lock):
    # the agency consults about the winners; if all are ready, the winners are sent
    ready = False
    logging.info(f'action: consulta_ganadores | agencia: {agency_id}')
    agency_register_lock.acquire()
    bets_lock.acquire()

    ready = agency_register.finish()
    agency_register_lock.release()  
    if ready:
        data = get_winners(agency_id)
        send_msg(client_socket,data, WIN_TYPE)
    else:
        send_msg(client_socket,b'', CHECK_WIN_TYPE)
    
    bets_lock.release()
    return False

```