def add_noise(bits, noise_rate):
    '''
        Return a new array with the final value
    '''
    pass

def text_to_bits(str):
    '''
        Return an array of bits
    '''
    final = []
    final.append(format(len(str), '08b'))
    for c in str:
        final.append(format(ord(c), '08b'))
    
    return final

def apply_algorithm(str, algorithm=""):
    pass

if __name__ == '__main__':
    # str = input("Ingresar mensaje: \n")
    # algorithm = input("Ingresar")
    
    # Start Socket connection

    a = text_to_bits("helloWorld")
    print(a)