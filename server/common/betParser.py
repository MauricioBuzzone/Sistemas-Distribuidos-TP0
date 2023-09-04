import logging

from common.utils import Bet

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
