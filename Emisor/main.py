from crc32 import crc32_generate_message
from viterbi import ConvolutionalEncoder
import random

def add_noise(bits, noise_rate):
    '''
        Return a new array with the final value
    '''
    final = []
    for s in range(bits):
        result = ""
        for bit in s:
            if random.random() < noise_rate:
                # Flip the bit
                flipped_bit = '1' if bit == '0' else '0'
                result.append(flipped_bit)
            else:
                result.append(bit)
        final.append(''.join(result))
    return final

def text_to_bits(str):
    '''
        Return an array of bits
    '''
    final = []
    final.append(format(len(str), '08b'))
    for c in str:
        final.append(format(ord(c), '08b'))
    
    return final


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

if __name__ == "__main__":

    uwu = text_to_bits("helloWorld")
    print(uwu)
    print(add_noise(uwu, 0.5))

    text      = input("Ingresar mensaje ASCII: ")
    algorithm = input("Seleccionar algoritmo (crc32/viterbi): ")

    bit_stream = text_to_bits(text)
    framed     = apply_algorithm(bit_stream, algorithm)
    print("Trama con integridad/codificación:", framed)

    raw_rate    = input("Ingrese la tasa de error (ej. '4/100' o '0.04'): ")
    noise_rate  = parse_error_rate(raw_rate)
    print(f"Aplicando tasa de error de {noise_rate:.2%}")

    noisy_frame = add_noise(framed, noise_rate)
    print("Trama con ruido           :", noisy_frame)

    errors       = sum(a != b for a, b in zip(framed, noisy_frame))
    ber_measured = errors / len(framed)
    print(f"Tasa de error medida (BER): {ber_measured:.2%}")