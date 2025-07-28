def crc32_generate_message(binary_str: str, poly: str = "100000100110000010001110110110111") -> str:
    # Fill up string to get up to 32 bits.
    degree = len(poly) - 1
    padded_msg = binary_str + '0' * degree
    padded_msg = list(padded_msg)
    poly = list(poly)

    for i in range(len(binary_str)):
        if padded_msg[i] == '1':
            for j in range(len(poly)):
                padded_msg[i + j] = str(int(padded_msg[i + j]) ^ int(poly[j]))

    # Adding the remainder as the padding
    crc = ''.join(padded_msg[-degree:])
    full_msg = binary_str + crc
    return full_msg

if __name__ == "__main__":
    message = input("Enter a binary number:")
    msg_with_crc = crc32_generate_message(message)
    print("Message with CRC-32:", msg_with_crc)