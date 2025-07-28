from crc32 import crc32_generate_message
from viterbi import ConvolutionalEncoder
import random
import socket
import struct

HOST = "localhost"
PORT = 8081


#########################
# UTILITIES
#########################

def text_to_binary(text):
    """Convert text to binary representation"""
    return ''.join(format(ord(char), '08b') for char in text)

def apply_algorithm(bit_str: str, algorithm: str = "") -> str:
    """
    Aplica el algoritmo seleccionado a la secuencia de bits:
      - 'crc32' o 'crc-32': devuelve bit_str + CRC-32  (detección)
      - 'viterbi': devuelve la secuencia convolucional codificada (corrección, tasa 1/2)
    """
    algo = algorithm.strip().lower()
    if algo in ("crc32", "crc-32"):
        return crc32_generate_message(bit_str)
    elif algo == "viterbi":
        bits = [int(b) for b in bit_str]
        encoder = ConvolutionalEncoder(generators=[0b111, 0b101])
        encoded = encoder.encode(bits)
        return "".join(str(b) for b in encoded)
    else:
        raise ValueError(f"Unsupported algorithm: {algorithm!r}")

def parse_error_rate(rate_str: str) -> float:
    """
    Acepta entradas como '4/100' o '0.04' y devuelve un float entre 0 y 1.
    """
    rate_str = rate_str.strip()
    if "/" in rate_str:
        num, den = rate_str.split("/", 1)
        return float(num) / float(den)
    return float(rate_str)

def add_noise(bit_str: str, noise_rate: float) -> str:
    """
    Invierte cada bit en bit_str con probabilidad noise_rate.
    """
    result = []
    for b in bit_str:
        if random.random() < noise_rate:
            result.append("1" if b == "0" else "0")
        else:
            result.append(b)
    return "".join(result)

##########################
#   SOCKETS
##########################

def send_message(sock, message):
    """Send a message with length prefix"""
    # Encode the message
    message_bytes = message.encode('utf-8')
    # Pack the length as a 4-byte integer (network byte order)
    length_prefix = struct.pack('!I', len(message_bytes))
    # Send length first, then message
    sock.sendall(length_prefix + message_bytes)

def recv_message(sock):
    """Receive a message with length prefix"""
    # First, receive the 4-byte length prefix
    length_bytes = recv_all(sock, 4)
    if not length_bytes:
        return None
    
    # Unpack the length
    message_length = struct.unpack('!I', length_bytes)[0]
    
    # Now receive the actual message
    message_bytes = recv_all(sock, message_length)
    return message_bytes.decode('utf-8')

def recv_all(sock, length):
    """Helper function to receive exactly 'length' bytes"""
    data = b''
    while len(data) < length:
        chunk = sock.recv(length - len(data))
        if not chunk:
            raise ConnectionError("Socket closed unexpectedly")
        data += chunk
    return data

##########################
#   MAIN
##########################

if __name__ == "__main__":

    option = input("Seleccionar algoritmo crc32 (0) viterbi (1): ")
    algorithm = ""
    
    if option == "0" :
        algorithm = "crc32"
    elif option == "1":
        algorithm = "viterbi"
    else:
        print("Invalid option.")
        exit(1)

    raw_rate    = input("Ingrese la tasa de error (ej. '4/100' o '0.04'): ")
    noise_rate  = parse_error_rate(raw_rate)

    while True :
        text = input("Ingresar mensaje ASCII: ")
        binary_msg = text_to_binary(text)
        print("Binary message:", binary_msg)
        encoded = apply_algorithm(binary_msg, algorithm)
        print(f"Encoded message with {algorithm}:", encoded)
        final_msg = add_noise(encoded, noise_rate)
        print(f"Message after applying noise rate of {noise_rate}:", final_msg)

        client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        client_socket.connect((HOST, PORT))

        print(f"Sending message of length: {len(final_msg)}")
        send_message(client_socket, final_msg)

        # Receive response
        response = recv_message(client_socket)
        print(f"Received response: {response[:50]}...")  # Show first 50 chars
        
        client_socket.close()
