-- Reads a line of '0'/'1' chars and returns a table of numbers
local function read_encoded_sequence()
    io.write("Enter the encoded bit sequence: ")
    local line = io.read("*l")
    local bits = {}
    for c in line:gmatch(".") do
        if c == "0" or c == "1" then
            table.insert(bits, tonumber(c))
        end
    end
    return bits
end

-- Main
local received = read_encoded_sequence()
print("Received " .. #received .. " bits.")
-- Next: initialize Viterbi structures (path metrics, survivor paths) and start the decoding loop.
