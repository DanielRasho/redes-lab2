#!/usr/bin/env python3
import random
import time
import argparse
import socket
import struct
import string

import pandas as pd
import matplotlib.pyplot as plt

# Importamos vuestros helpers de Emisor/main.py
from main import apply_algorithm, parse_error_rate, add_noise, send_message, recv_message

HOST = "localhost"
PORT = 8081

# Usaremos sólo caracteres ASCII imprimibles (espacio .. ~)
PRINTABLE = string.ascii_letters + string.digits + string.punctuation + " "

def run_trial(msg_str: str, algo: str, error_rate: float):
    """
    - msg_str: cadena de ASCII imprimible
    - algo: "crc32" o "viterbi"
    - error_rate: float entre 0 y 1
    """
    raw_bytes = msg_str.encode('ascii')
    bits = "".join(f"{b:08b}" for b in raw_bytes)
    coded = apply_algorithm(bits, algo)
    noisy = add_noise(coded, error_rate)

    # envío por socket
    t0 = time.time()
    sock = socket.socket()
    sock.connect((HOST, PORT))
    send_message(sock, noisy)
    response = recv_message(sock)
    sock.close()
    elapsed = time.time() - t0

    # parseamos la parte decodificada
    decoded = ""
    if response is None:
        decoded = ""
    else:
        # Primero, buscamos la comilla de Viterbi
        if '"' in response:
            decoded = response.split('"')[-2]
        # Sino, si es CRC32 “No error detected; decoded: foo”
        elif "decoded:" in response:
            decoded = response.split("decoded:")[-1].strip()
        # Sino, nada
    ok = (decoded == msg_str)
    return ok, len(coded), elapsed

def main():
    p = argparse.ArgumentParser()
    p.add_argument('--algorithm', choices=['crc32','viterbi'], required=True)
    p.add_argument('--sizes', nargs='+', type=int,
                   default=[16,128,1024,4096])
    p.add_argument('--rates', nargs='+', type=float,
                   default=[0.0,0.01,0.02,0.05,0.1])
    p.add_argument('--trials', type=int, default=200)
    p.add_argument('--output', type=str, default='results.csv')
    args = p.parse_args()

    records = []
    for size in args.sizes:
        for p_error in args.rates:
            succ = 0
            sum_over = 0.0
            sum_time = 0.0
            for _ in range(args.trials):
                # Generar mensaje ASCII imprimible
                msg = "".join(random.choice(PRINTABLE) for _ in range(size))
                ok, coded_len, t = run_trial(msg, args.algorithm, p_error)
                succ += int(ok)
                sum_over += (coded_len - 8*size)/(8*size)
                sum_time += t
            records.append({
                'algorithm': args.algorithm,
                'size_bytes': size,
                'error_rate': p_error,
                'success_rate': succ/args.trials,
                'avg_overhead': sum_over/args.trials,
                'avg_time_ms': sum_time/args.trials*1000
            })
            print(f"{args.algorithm:>7} | size={size:4d} B | p={p_error:.3f} "
                  f"→ éxito={succ/args.trials:.3f}, overhead={(sum_over/args.trials):.3f}")

    df = pd.DataFrame(records)
    df.to_csv(args.output, index=False)
    print(f"\nResultados guardados en {args.output}")

    # Gráfica éxito vs error
    for size in df['size_bytes'].unique():
        sub = df[df['size_bytes']==size]
        plt.figure()
        plt.plot(sub['error_rate'], sub['success_rate'],
                 marker='o', label=args.algorithm)
        plt.title(f'Éxito vs Error ({args.algorithm}, size={size} B)')
        plt.xlabel('Tasa de error')
        plt.ylabel('Tasa de éxito')
        plt.ylim(-0.05,1.05)
        plt.grid(True)
        plt.legend()
        fn = f'success_{args.algorithm}_{size}B.png'
        plt.savefig(fn)
        plt.close()
        print(f"  → {fn}")

if __name__ == '__main__':
    main()
