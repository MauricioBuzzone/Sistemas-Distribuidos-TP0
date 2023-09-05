# TP0: Docker + Comunicaciones + Concurrencia

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