import logging
import struct

from common.protocol import recv_msg, send_msg, SocketBroken
from common.protocol import BET_TYPE, OK_TYPE, ERR_TYPE, END_TYPE, WIN_TYPE,CHECK_WIN_TYPE
from common.betParser import parser_bet, get_winners
from common.utils import store_bets

def handle_client_connection(client_socket, agency_register, agency_register_lock, bets_lock):
    """
    Read message from a specific client socket and closes the socket

    If a problem arises in the communication with the client, the
    client socket will also be closed
    """
    try:
        client_sending = True
        while client_sending:
            if client_socket:
                msg = recv_msg(client_socket)
                client_sending = handle_message(client_socket, msg, agency_register, agency_register_lock, bets_lock)
            else:
                break
    except (SocketBroken,OSError) as e:
        logging.error(f'action: receive_message | result: fail | error: {e}')
    finally:
        if client_socket:
            logging.info(f'action: release_client_socket | result: success')
            client_socket.close()


def handle_message(client_socket, msg, agency_register, agency_register_lock, bets_lock):
    type_msg = msg[0]
    if type_msg == BET_TYPE:
        return handle_bets(client_socket,msg, bets_lock)
    if type_msg == END_TYPE:
       return handle_finish(msg[1], agency_register, agency_register_lock)
    if type_msg == CHECK_WIN_TYPE:
        return handle_winners(client_socket, msg[1], agency_register, agency_register_lock,bets_lock)

    return True

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


def handle_finish(agency_id, agency_register, agency_register_lock):
    # Update AgencyRegister to indicate that agency has finished sending bets.
    logging.info(f'action: finish_stored | agencia: {int(agency_id)}')
    agency_register_lock.acquire()
    agency_register.update(int(agency_id))
    agency_register_lock.release()
    return False


def handle_bets(client_socket, msg,bets_lock):
    try:
        # The bets are parsed from the information received
        bets = parser_bet(msg[1:])

        bets_lock.acquire()
        store_bets(bets)
        bets_lock.release()
        logging.info(f'action: apuestas_almacenadas: result: sucess | amount: {len(bets)}')

        # Message to the customer with the quantity of stored bets.
        data = serialize_amount_bets(bets)
        send_msg(client_socket,data, OK_TYPE)

    except OSError as e:
        logging.error(f'action: error_bets | result: fail | error: {e}')
        # Send message to notify the client
        send_msg(client_socket,b'', ERR_TYPE)

    return True

def serialize_amount_bets(bets):
    data = b''
    amount_data = str(len(bets)).encode('utf-8')
    amount_data_size = struct.pack('!i',len(amount_data))
    data += amount_data_size
    data += amount_data

    return data
