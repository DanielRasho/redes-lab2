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
	// Leer secuencia recibida
	received := readEncodedSequence(input)
	fmt.Printf("Received %d bits.\n", len(received))

	// Parámetros del código convolucional
	generators := []int{0x7, 0x5} // polinomios generadores
	memory := 2
	numStates := 1 << memory
	// INF más grande que cualquier métrica posible
	INF := len(received) + 1

	// Inicialización de métricas de camino y sobrevivientes
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

	// Recorrer la secuencia de dos en dos (len(generators)=2)
	for i := 0; i < len(received); i += len(generators) {
		// Observaciones en este paso
		obs := make([]int, len(generators))
		for j := range generators {
			obs[j] = received[i+j]
		}

		// Preparar las estructuras para el siguiente paso
		nextMetrics := make([]int, numStates)
		nextSurvivors := make([][]int, numStates)
		for s := 0; s < numStates; s++ {
			nextMetrics[s] = INF
			nextSurvivors[s] = []int{}
		}

		// Para cada estado actual
		for s := 0; s < numStates; s++ {
			pm := pathMetrics[s]
			if pm < INF {
				// Para cada posible bit de entrada (0 o 1)
				for ib := 0; ib <= 1; ib++ {
					// Construir el registro de desplazamiento [ib, memoria bits...]
					sr := make([]int, memory+1)
					sr[0] = ib
					for k := 0; k < memory; k++ {
						sr[k+1] = (s >> (memory - 1 - k)) & 1
					}

					// Calcular bits esperados según los generadores
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

					// Calcular distancia de Hamming (número de diferencias)
					d := 0
					for k := range exp {
						if exp[k] != obs[k] {
							d++
						}
					}

					// Estado siguiente y nueva métrica
					ns := (s >> 1) | (ib << (memory - 1))
					nm := pm + d

					// Actualizar métrica y camino si es mejor
					if nm < nextMetrics[ns] {
						nextMetrics[ns] = nm
						// Copiar camino sobreviviente y agregar ib
						nextSurvivors[ns] = append([]int(nil), survivors[s]...)
						nextSurvivors[ns] = append(nextSurvivors[ns], ib)
					}
				}
			}
		}

		// Avanzar al siguiente paso
		pathMetrics = nextMetrics
		survivors = nextSurvivors
	}

	// Encontrar el estado final con métrica mínima
	bestState := 0
	bestMetric := pathMetrics[0]
	for s := 1; s < numStates; s++ {
		if pathMetrics[s] < bestMetric {
			bestState = s
			bestMetric = pathMetrics[s]
		}
	}

	// Recuperar y mostrar la secuencia decodificada
	decoded := survivors[bestState]
	fmt.Print("Decoded sequence: ")
	for _, b := range decoded {
		fmt.Print(b)
	}
	fmt.Println()

	return nil
}
