def get_trama():
    """
    Solicita al usuario una secuencia de bits y la valida.
    Devuelve la trama como una cadena de '0' y '1'.
    """
    while True:
        trama = input("Introduce the binary plot ")
        if all(bit in '01' for bit in trama) and len(trama) > 0:
            return trama
        print("Error: invalid trama")

if __name__ == "__main__":
    trama_entrada = get_trama()
    print(f"Trama recibida: {trama_entrada}")
