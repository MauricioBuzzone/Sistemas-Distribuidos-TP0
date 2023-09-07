import struct

LENGHT = 4
BET_TYPE = 'B'
END_TYPE = 'F'
OK_TYPE = '0'
ERR_TYPE = 'E'
WIN_TYPE = 'W'
CHECK_WIN_TYPE = 'Q'

class SocketBroken(Exception):
    pass

def read_all(client_sock, bytes_to_read):
    """
        Recv all n bytes to avoid short read
    """
    data = b''
    while len(data) < bytes_to_read:
        bytes = client_sock.recv(bytes_to_read - len(data))
        if bytes == b'':
            raise SocketBroken()
        data += bytes
    return data

def send_all(client_sock, msg):
    """
        Recv all n bytes to avoid short read
    """
    bytesSended = 0
    while bytesSended < len(msg):
        b = client_sock.send(msg[bytesSended:])
        if b == 0:
            raise SocketBroken()
        bytesSended += b
    return bytesSended



def send_msg(client_sock, msg, type):
    # Send the payload msj
    size_msg = struct.pack('!i', len(msg))
    send_all(client_sock,size_msg)

    # Send the type of message
    type_data = type.encode('utf-8')
    send_all(client_sock,type_data)

    # Send msj
    send_all(client_sock,msg)



def recv_msg(client_sock):
    read_bytes = 0
    fields = []
    
    # Read size message
    len_data = read_all(client_sock, LENGHT)
    len_msj = int.from_bytes(len_data, byteorder='big')

    # Read type message
    type_data = read_all(client_sock, 1)
    type_msj = type_data.decode('utf-8')
    fields.append(type_msj)

    # Read message
    while read_bytes < len_msj:
        len_field_data = read_all(client_sock, LENGHT)
        len_field = int.from_bytes(len_field_data, byteorder='big')
        read_bytes += LENGHT

        field_data = read_all(client_sock, len_field)
        field = field_data.decode('utf-8')
        read_bytes += len_field
        fields.append(field)

    return fields


