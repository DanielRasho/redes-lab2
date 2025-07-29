package main

import (
	"fmt"
	"strings"
)

// readEncodedSequence convierte el input en un slice de 0/1.
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

/*
checkViterbi decodifica con Viterbi **y** devuelve:
  - decoded: slice de bits corregidos (la ruta ganadora)
  - errPos: posiciones donde había discrepancias
  - err:    non-nil si bestMetric>0, con mensaje que incluye conteo y posiciones

Firma nueva:
*/
func checkViterbi(input string) (decoded []int, errPos []int, err error) {
	received := readEncodedSequence(input)
	fmt.Printf("Received %d bits.\n", len(received))

	generators := []int{0x7, 0x5}
	if len(received)%len(generators) != 0 {
		return nil, nil, fmt.Errorf(
			"invalid input length: expected multiple of %d bits, got %d",
			len(generators), len(received),
		)
	}

	memory := 2
	numStates := 1 << memory
	INF := len(received) + 1

	// métricas, caminos y posiciones por estado
	pathMetrics := make([]int, numStates)
	survivors := make([][]int, numStates)
	posHist := make([][]int, numStates)
	for s := 0; s < numStates; s++ {
		if s == 0 {
			pathMetrics[s] = 0
		} else {
			pathMetrics[s] = INF
		}
	}

	globalIndex := 0

	for i := 0; i < len(received); i += len(generators) {
		obs := make([]int, len(generators))
		for j := range generators {
			obs[j] = received[i+j]
		}

		nextMetrics := make([]int, numStates)
		nextSurvivors := make([][]int, numStates)
		nextPosHist := make([][]int, numStates)
		for s := 0; s < numStates; s++ {
			nextMetrics[s] = INF
		}

		for s := 0; s < numStates; s++ {
			pm := pathMetrics[s]
			if pm < INF {
				for ib := 0; ib <= 1; ib++ {
					// registro de desplazamiento
					sr := make([]int, memory+1)
					sr[0] = ib
					for k := 0; k < memory; k++ {
						sr[k+1] = (s >> (memory - 1 - k)) & 1
					}

					// salida esperada
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

					// distancia de Hamming y posiciones de error
					d := 0
					errs := []int{}
					for k := range exp {
						if exp[k] != obs[k] {
							d++
							errs = append(errs, globalIndex+k)
						}
					}

					ns := (s >> 1) | (ib << (memory - 1))
					nm := pm + d

					if nm < nextMetrics[ns] {
						nextMetrics[ns] = nm
						// copiar camino y anexar ib
						nextSurvivors[ns] = append([]int(nil), survivors[s]...)
						nextSurvivors[ns] = append(nextSurvivors[ns], ib)
						// copiar historial de posiciones y anexar errores
						nextPosHist[ns] = append([]int(nil), posHist[s]...)
						nextPosHist[ns] = append(nextPosHist[ns], errs...)
					}
				}
			}
		}

		pathMetrics = nextMetrics
		survivors = nextSurvivors
		posHist = nextPosHist
		globalIndex += len(generators)
	}

	// escoger estado final
	bestState := 0
	bestMetric := pathMetrics[0]
	for s := 1; s < numStates; s++ {
		if pathMetrics[s] < bestMetric {
			bestState = s
			bestMetric = pathMetrics[s]
		}
	}

	decoded = survivors[bestState]
	errPos = posHist[bestState]

	fmt.Print("Decoded sequence: ")
	for _, b := range decoded {
		fmt.Print(b)
	}
	fmt.Println()

	if bestMetric > 0 {
		return decoded, errPos, fmt.Errorf(
			"sequence corrupted: detected %d bit errors at positions %v",
			bestMetric, errPos,
		)
	}
	return decoded, nil, nil
}
