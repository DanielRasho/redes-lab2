def get_binary_input():
    while True:
        user_input = input("Enter the binary sequence to encode: ")
        if user_input and all(bit in '01' for bit in user_input):
            return [int(bit) for bit in user_input]
        print("Error: Only '0' and '1' characters are allowed.")

class ConvolutionalEncoder:
    def __init__(self, generators):
        """Initialize with generator polynomials in binary form."""
        self.generators = generators
        self.memory = max(g.bit_length() - 1 for g in generators)
        self.state = [0] * self.memory  # Initialize shift register

    def encode(self, bits):
        """Encodes a sequence of bits using convolutional coding."""
        encoded_bits = []
        for bit in bits:
            # Current window = new bit + previous state
            window = [bit] + self.state
            
            # Calculate output bits for each generator
            for generator in self.generators:
                # Sum bits where generator mask is 1
                encoded_bit = sum(bit for i, bit in enumerate(window) if (generator >> i) & 1) % 2
                encoded_bits.append(encoded_bit)
            
            # Update state (shift right)
            self.state = window[:-1]
        
        return encoded_bits

if __name__ == "__main__":
    input_bits = get_binary_input()

    encoder = ConvolutionalEncoder(generators=[0b111, 0b101])

    encoded_sequence = encoder.encode(input_bits)
    print("Encoded sequence:", ''.join(map(str, encoded_sequence)))