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
checkViterbi decodifica con el algoritmo de Viterbi y, en caso de error,
devuelve un fmt.Errorf que incluye tanto la métrica como las posiciones
de bit donde hubo discrepancias.

	input    cadena de bits codificada (ej. "10010011...")
	error    si bestMetric>0: "sequence corrupted: detected X bit errors at positions [p1 p2 ...]"
*/func checkViterbi(input string) error {
	// 1. Leer secuencia recibida
	received := readEncodedSequence(input)
	fmt.Printf("Received %d bits.\n", len(received))

	// 2. Generadores y comprobación de longitud
	generators := []int{0x7, 0x5} // polinomios generadores (111, 101)
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

	// 4. Inicialización de métricas, caminos sobrevivientes y posiciones de error por estado
	pathMetrics := make([]int, numStates)
	survivors := make([][]int, numStates)
	posHist := make([][]int, numStates)
	for s := 0; s < numStates; s++ {
		if s == 0 {
			pathMetrics[s] = 0
		} else {
			pathMetrics[s] = INF
		}
		survivors[s] = nil
		posHist[s] = nil
	}

	// Índice global de bit (se incrementa de a len(generators))
	globalIndex := 0

	// 5. Viterbi paso a paso
	for i := 0; i < len(received); i += len(generators) {
		// Observaciones en este bloque de len(generators) bits
		obs := make([]int, len(generators))
		for j := range generators {
			obs[j] = received[i+j]
		}

		// Preparar estructuras para el siguiente paso
		nextMetrics := make([]int, numStates)
		nextSurvivors := make([][]int, numStates)
		nextPosHist := make([][]int, numStates)
		for s := 0; s < numStates; s++ {
			nextMetrics[s] = INF
			nextSurvivors[s] = nil
			nextPosHist[s] = nil
		}

		// Para cada estado actual
		for s := 0; s < numStates; s++ {
			pm := pathMetrics[s]
			if pm < INF {
				// Para cada bit de entrada posible (0 o 1)
				for ib := 0; ib <= 1; ib++ {
					// 5.1 Construir registro de desplazamiento [ib, memoria bits...]
					sr := make([]int, memory+1)
					sr[0] = ib
					for k := 0; k < memory; k++ {
						sr[k+1] = (s >> (memory - 1 - k)) & 1
					}

					// 5.2 Calcular bits esperados según cada generador
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

					// 5.3 Distancia de Hamming y registrar posiciones de discrepancia
					d := 0
					var errs []int
					for k := range exp {
						if exp[k] != obs[k] {
							d++
							errs = append(errs, globalIndex+k)
						}
					}

					// 5.4 Estado siguiente y métrica acumulada
					ns := (s >> 1) | (ib << (memory - 1))
					nm := pm + d

					// 5.5 Si este camino es mejor, guardarlo
					if nm < nextMetrics[ns] {
						nextMetrics[ns] = nm
						// copiar camino sobreviviente y anexar ib
						nextSurvivors[ns] = append([]int(nil), survivors[s]...)
						nextSurvivors[ns] = append(nextSurvivors[ns], ib)
						// copiar historial de posiciones y anexar discrepancias
						nextPosHist[ns] = append([]int(nil), posHist[s]...)
						nextPosHist[ns] = append(nextPosHist[ns], errs...)
					}
				}
			}
		}

		// Avanzar a la siguiente iteración
		pathMetrics = nextMetrics
		survivors = nextSurvivors
		posHist = nextPosHist
		globalIndex += len(generators)
	}

	// 6. Elegir el estado final con métrica mínima
	bestState := 0
	bestMetric := pathMetrics[0]
	for s := 1; s < numStates; s++ {
		if pathMetrics[s] < bestMetric {
			bestState = s
			bestMetric = pathMetrics[s]
		}
	}

	// 7. Mostrar secuencia decodificada (opcional)
	decoded := survivors[bestState]
	fmt.Print("Decoded sequence: ")
	for _, b := range decoded {
		fmt.Print(b)
	}
	fmt.Println()

	// 8. Si hubo errores, devolverlos con las posiciones correspondientes
	if bestMetric > 0 {
		return fmt.Errorf(
			"sequence corrupted: detected %d bit errors at positions %v",
			bestMetric, posHist[bestState],
		)
	}
	return nil
}
