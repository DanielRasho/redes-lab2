local bit = require("bit")

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

local received = read_encoded_sequence()
print("Received " .. #received .. " bits.")

local generators = {0x7, 0x5}
local memory     = 2
local num_states = 2 ^ memory
local INF        = math.huge

local path_metrics = {}
local survivors    = {}
for state = 0, num_states - 1 do
    path_metrics[state] = (state == 0) and 0 or INF
    survivors[state]    = {}
end

local bit = require("bit")

for i = 1, #received, #generators do
    local obs = {}
    for j = 1, #generators do
        obs[j] = received[i + j - 1]
    end

    local next_metrics   = {}
    local next_survivors = {}
    for s = 0, num_states - 1 do
        next_metrics[s]   = INF
        next_survivors[s] = {}
    end

    for s = 0, num_states - 1 do
        local pm = path_metrics[s]
        if pm < INF then
            for ib = 0, 1 do
                local sr = { ib }
                for k = 1, memory do
                    sr[k+1] = bit.band(bit.rshift(s, memory - k), 1)
                end

                local exp = {}
                for gi, g in ipairs(generators) do
                    local x = 0
                    for k = 1, #sr do
                        if bit.band(g, bit.lshift(1, k-1)) ~= 0 then
                            x = bit.bxor(x, sr[k])
                        end
                    end
                    exp[gi] = x
                end

                local d = 0
                for k = 1, #exp do
                    if exp[k] ~= obs[k] then d = d + 1 end
                end

                local ns = bit.bor(bit.rshift(s, 1), bit.lshift(ib, memory - 1))
                local nm = pm + d

                if nm < next_metrics[ns] then
                    next_metrics[ns]   = nm
                    next_survivors[ns] = { unpack(survivors[s]) }
                    table.insert(next_survivors[ns], ib)
                end
            end
        end
    end

    path_metrics = next_metrics
    survivors    = next_survivors
end

local best_state, best_metric = 0, path_metrics[0]
for s = 1, num_states - 1 do
    if path_metrics[s] < best_metric then
        best_state, best_metric = s, path_metrics[s]
    end
end

local decoded = survivors[best_state]
print("Decoded sequence:", table.concat(decoded))

