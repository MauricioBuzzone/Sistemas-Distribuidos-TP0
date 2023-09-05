import logging
import struct

from common.utils import Bet, load_bets, has_won

def parser_bet(bets_attribute):
    bets = []
    agency = bets_attribute[0]
    i = 1
    while i <= len(bets_attribute)-5:
        firstName = bets_attribute[i]
        lastName = bets_attribute[i+1]
        document = bets_attribute[i+2]
        birthdate = bets_attribute[i+3]
        number = bets_attribute[i+4]
        i+=5

        bet = Bet(agency,firstName,lastName,document,birthdate,number)
        bets.append(bet)

    return bets

def parse_document(document):
    bytes = b''
    doc_data = document.encode('utf-8')
    doc_data_size = struct.pack('!i',len(doc_data))
    bytes += doc_data_size
    bytes += doc_data
    return bytes


def get_winners(agency_id):
    winners_data = b''
    for bet in load_bets():
        if has_won(bet) and int(bet.agency) == int(agency_id):
            winners_data += parse_document(bet.document)

    return winners_data