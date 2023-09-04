import socket
import logging
import signal
import struct
from common.protocol import recv_msg, send_msg
from common.protocol import BET_TYPE, OK_TYPE, ERR_TYPE, END_TYPE
from common.utils import Bet, store_bets
from common.betParser import parser_bet
class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._server_on = True
        signal.signal(signal.SIGINT, self.__handle_signal)
        signal.signal(signal.SIGTERM, self.__handle_signal)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while self._server_on:
            client_sock = self.__accept_new_connection()
            if self._server_on:
                self.__handle_client_connection(client_sock)


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

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            client_sending = True
            while client_sending:
                msg = recv_msg(client_sock)
                client_sending = self.__handle_message(client_sock, msg)
            client_sock.close()
            
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")


    def __handle_message(self, client_sock, msg):
        type_msg = msg[0]
        if type_msg == BET_TYPE:
            try:
                bets = parser_bet(msg[1:])
                store_bets(bets)
                logging.info(f'action: apuestas_almacenadas: result: sucess | amount: {len(bets)}')

                # Send message to notify the client
                data = b''
                amount_data = str(len(bets)).encode('utf-8')
                amount_data_size = struct.pack('!i',len(amount_data))
                data += amount_data_size
                data += amount_data
                send_msg(client_sock,data, OK_TYPE)
                return True

            except OSError as e:
                logging.error("action: error_bets | result: fail | error: {e}")
                # Send message to notify the client
                send_msg(client_sock,b'', ERR_TYPE)

        if type_msg == END_TYPE:
            return False


    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        try:
            logging.info('action: accept_connections | result: in_progress')
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
            return c
        except OSError as e:
            if self._server_on:
                logging.error(f'action: accept_connections | result: fail | error: {e}')
            else:
                logging.info(f'action: stop_accept_connections | result: success')
            return None