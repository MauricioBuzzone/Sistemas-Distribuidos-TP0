import socket
import logging
import signal
import struct
import threading

from threading import Thread, Lock
from common.protocol import recv_msg, send_msg
from common.protocol import BET_TYPE, OK_TYPE, ERR_TYPE, END_TYPE, WIN_TYPE,CHECK_WIN_TYPE
from common.utils import Bet, store_bets
from common.betParser import parser_bet, get_winners
from common.agencyRegister import AgencyRegister
from common.handleClient import handle_client_connection

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        signal.signal(signal.SIGINT, self.__handle_signal)
        signal.signal(signal.SIGTERM, self.__handle_signal)
        self._server_on = True
        self._agency_register = AgencyRegister(listen_backlog)
        self._agency_register_lock = threading.Lock()
        self._bets_lock = threading.Lock()


    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        workers = []
        while self._server_on:
            client_sock = self.__accept_new_connection()

            worker = threading.Thread(
            target=handle_client_connection, args=(
                client_sock, 
                self._agency_register,
                self._agency_register_lock,
                self._bets_lock))
            worker.start()
            workers.append(worker)

        for worker in workers:
            worker.join()
        logging.info(f'action: join_threads | result: success')


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
                logging.error(f'action: accept_connections | result: fail')
            else:
                logging.info(f'action: stop_accept_connections | result: success')
            return