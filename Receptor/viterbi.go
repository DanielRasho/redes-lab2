package main

import (
	"fmt"
	"strings"
)

func readEncodedSequence(input string) []int {
	line := strings.TrimSpace(input)
	bits := []int{}
	for _, c := range line {
		if c == '0' || c == '1' {
			bits = append(bits, int(c-'0'))
		}
	}
	return bits
}

// TODO: Return descriptive error messages where input has been corrupted
// through the error object.

/*
Attempts to decodes a binary string encoded with viterbi

	input	binary string (ex: "10000010101") encoded with viterbi
	error	object returning possible error encountered in the message received
*/
func checkViterbi(input string) error {
	// 1. Leer secuencia recibida
	received := readEncodedSequence(input)
	fmt.Printf("Received %d bits.\n", len(received))

	// 2. Generadores y comprobación de longitud
	generators := []int{0x7, 0x5} // polinomios generadores
	if len(received)%len(generators) != 0 {
		return fmt.Errorf(
			"invalid input length: expected a multiple of %d bits, got %d",
			len(generators), len(received),
		)
	}

	// 3. Parámetros del código convolucional
	memory := 2
	numStates := 1 << memory
	INF := len(received) + 1

	// 4. Inicialización de métricas y caminos sobrevivientes
	pathMetrics := make([]int, numStates)
	survivors := make([][]int, numStates)
	for s := 0; s < numStates; s++ {
		if s == 0 {
			pathMetrics[s] = 0
		} else {
			pathMetrics[s] = INF
		}
		survivors[s] = []int{}
	}

	// 5. Viterbi paso a paso
	for i := 0; i < len(received); i += len(generators) {
		obs := make([]int, len(generators))
		for j := range generators {
			obs[j] = received[i+j]
		}

		nextMetrics := make([]int, numStates)
		nextSurvivors := make([][]int, numStates)
		for s := 0; s < numStates; s++ {
			nextMetrics[s] = INF
			nextSurvivors[s] = []int{}
		}

		for s := 0; s < numStates; s++ {
			pm := pathMetrics[s]
			if pm < INF {
				for ib := 0; ib <= 1; ib++ {
					// Construir registro de desplazamiento
					sr := make([]int, memory+1)
					sr[0] = ib
					for k := 0; k < memory; k++ {
						sr[k+1] = (s >> (memory - 1 - k)) & 1
					}
					// Salida esperada
					exp := make([]int, len(generators))
					for gi, g := range generators {
						x := 0
						for idx, bitVal := range sr {
							if (g>>idx)&1 == 1 {
								x ^= bitVal
							}
						}
						exp[gi] = x
					}
					// Distancia de Hamming
					d := 0
					for k := range exp {
						if exp[k] != obs[k] {
							d++
						}
					}
					// Estado siguiente y métrica acumulada
					ns := (s >> 1) | (ib << (memory - 1))
					nm := pm + d
					if nm < nextMetrics[ns] {
						nextMetrics[ns] = nm
						nextSurvivors[ns] = append([]int(nil), survivors[s]...)
						nextSurvivors[ns] = append(nextSurvivors[ns], ib)
					}
				}
			}
		}

		pathMetrics = nextMetrics
		survivors = nextSurvivors
	}

	// 6. Elegir el mejor estado final
	bestState := 0
	bestMetric := pathMetrics[0]
	for s := 1; s < numStates; s++ {
		if pathMetrics[s] < bestMetric {
			bestState = s
			bestMetric = pathMetrics[s]
		}
	}

	// 7. Mostrar secuencia decodificada
	decoded := survivors[bestState]
	fmt.Print("Decoded sequence: ")
	for _, b := range decoded {
		fmt.Print(b)
	}
	fmt.Println()

	// 8. Devolver error si hubo discrepancias en la decodificación
	if bestMetric > 0 {
		return fmt.Errorf("sequence corrupted: detected %d bit errors", bestMetric)
	}
	return nil
}
