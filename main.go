package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

//Cores do Print
var colorReset = "\033[0m"
var colorRed = "\033[31m"
var colorGreen = "\033[32m"
var colorYellow = "\033[33m"
var colorPurple = "\033[35m"
var colorCyan = "\033[36m"
var display = true

//Filosofo
type Filosofo struct {
	id     int
	estado int
	//-2 - Terminou
	//-1 - Morreu
	//0 - Pensando
	//1 - Com Fome
	//2 - Comendo
	filEsq   *Filosofo
	filDir   *Filosofo
	hashiEsq *Hashi
	hashiDir *Hashi
}

//Print Filosofo
func (f *Filosofo) String() string {
	res := "filosofo: " + fmt.Sprint(f.id) + "\n"
	res += "estado: " + fmt.Sprint(f.estado) + "\n"
	res += "filEsq: " + fmt.Sprint(f.filEsq.id) + "\n"
	res += "fisDir: " + fmt.Sprint(f.filDir.id) + "\n"
	res += "hashiEsq: " + fmt.Sprint(f.hashiEsq.id) + "\n"
	res += "hashiDir: " + fmt.Sprint(f.hashiDir.id) + "\n"

	return res
}

//Hashi
type Hashi struct {
	id         int
	disponivel bool
	reservado  int
	// -1 esquerdo reservou
	// 0 ninguem reservou
	// 1 direito reservou
}

func (f *Filosofo) comeca(wg *sync.WaitGroup) {
	defer wg.Done()

	//Iteracoes do filosofo
	for i := 0; i < 10; i++ {
		time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)

		//Estado aleatorio
		state := rand.Intn(2)

		//stateAnt := f.estado
		f.estado = state

		switch state {
		case 0:
			//if stateAnt != 0 {
			//fmt.Printf("Filosofo %d: %sPensando%s\n", f.id, string(colorGreen), string(colorReset))
			//}

			time.Sleep(time.Duration(rand.Intn(6)+4) * time.Second)

		case 1:
			//fmt.Printf("Filosofo %d: %sCom fome%s\n", f.id, string(colorYellow), string(colorReset))

			if f.hashiEsq.disponivel && f.hashiDir.disponivel && f.hashiEsq.reservado != -1 && f.hashiDir.reservado != 1 {
				//Comendo
				f.hashiEsq.disponivel = false
				f.hashiDir.disponivel = false

				//fmt.Printf("Filosofo %d: %sComendo com hashis %d e %d%s\n", f.id, string(colorCyan), f.hashiEsq.id, f.hashiDir.id, string(colorReset))
				f.estado = 2

				time.Sleep(time.Duration(rand.Intn(8)+2) * time.Second)

				//fmt.Printf("Filosofo %d: %sTerminou de comer. Liberando hashis %d e %d%s\n", f.id, string(colorPurple), f.hashiEsq.id, f.hashiDir.id, string(colorReset))
				f.estado = 0

				f.hashiEsq.disponivel = true
				f.hashiDir.disponivel = true
			} else {

				//Esperando
				//fmt.Printf("Filosofo %d: %sEsperando para comer%s\n", f.id, string(colorYellow), string(colorReset))

				for i := 0; !f.hashiEsq.disponivel || !f.hashiDir.disponivel || f.hashiEsq.reservado == -1 || f.hashiDir.reservado == 1; i++ {

					if !f.hashiEsq.disponivel && !f.hashiDir.disponivel {
						f.hashiEsq.reservado = 1
						f.hashiDir.reservado = -1
					} else if !f.hashiEsq.disponivel {
						f.hashiEsq.reservado = 1
					} else if !f.hashiDir.disponivel {
						f.hashiDir.reservado = -1
					}

					//Morreu
					if i > int(math.Pow(10, 9)) {

						f.estado = -1

						//Desreservar hashis
						if f.hashiEsq.reservado == 1 {
							f.hashiEsq.reservado = 0
						}
						if f.hashiDir.reservado == -1 {
							f.hashiDir.reservado = 0
						}

						//fmt.Printf("%sFilosofo %d: Morreu de fome%s\n", string(colorRed), f.id, string(colorReset))
						return
					}
				}

				//Comer depois de esperar
				f.hashiEsq.disponivel = false
				f.hashiDir.disponivel = false
				f.hashiEsq.reservado = 0
				f.hashiDir.reservado = 0

				//fmt.Printf("Filosofo %d: %sComendo com hashis %d e %d%s\n", f.id, string(colorCyan), f.hashiEsq.id, f.hashiDir.id, string(colorReset))
				f.estado = 2

				time.Sleep(time.Duration(rand.Intn(8)+2) * time.Second)

				//fmt.Printf("Filosofo %d: %sTerminou de comer. Liberando hashis %d e %d%s\n", f.id, string(colorPurple), f.hashiEsq.id, f.hashiDir.id, string(colorReset))
				f.estado = 0

				f.hashiEsq.disponivel = true
				f.hashiDir.disponivel = true

			}

		default:
			fmt.Println("!Erro de estado!")
		}

	}

	f.estado = -2
}

//Inicializar Mesa de Filosofos
func initFilosofos(total int) []*Filosofo {

	filosofos := make([]*Filosofo, total, total)
	auxHashi := &Hashi{1, true, 0}

	filosofos[0] = &Filosofo{
		id:       1,
		estado:   0,
		filEsq:   nil,
		filDir:   nil,
		hashiEsq: auxHashi,
		hashiDir: nil,
	}

	for i := 1; i < total; i++ {

		auxHashi = &Hashi{i + 1, true, 0}
		filosofos[i-1].hashiDir = auxHashi

		filosofos[i] = &Filosofo{
			id:       i + 1,
			estado:   0,
			filEsq:   filosofos[i-1],
			filDir:   nil,
			hashiEsq: auxHashi,
			hashiDir: nil,
		}

		filosofos[i-1].filDir = filosofos[i]

	}

	filosofos[total-1].filDir = filosofos[0]
	filosofos[total-1].hashiDir = filosofos[0].hashiEsq
	filosofos[0].filEsq = filosofos[total-1]

	return filosofos
}

func checarFilosofos(filosofos []*Filosofo) bool {

	for _, f := range filosofos {
		if f.estado != -2 && f.estado != -1 {
			return true
		}
	}

	return false
}

func main() {
	//Total de filosofos
	var total int
	fmt.Print("Numero de filosofos: ")
	fmt.Scanln(&total)

	rand.Seed(time.Now().UTC().UnixNano())

	//Criar Filosofos e Hashis
	filosofos := initFilosofos(total)

	//Print
	//for _, f := range filosofos {
	//	fmt.Println(f.String())
	//}

	//division := "-------------------------"
	//fmt.Print("\033[H\033[2J")
	//fmt.Println(division, "Inicio", division)

	//Grupo thread
	wg := new(sync.WaitGroup)
	wg.Add(total)

	//Inicializar Filosofos
	for _, f := range filosofos {
		//fmt.Printf("Filosofo %d: %sPensando%s\n", f.id, colorGreen, colorReset)
		go f.comeca(wg)
	}

	//Display
	if display {
		for {
			res := "\033[H\033[2J"
			time.Sleep(500 * time.Millisecond)

			for _, f := range filosofos {
				//Hashi
				if !f.hashiEsq.disponivel {
					res += colorRed
				} else if f.hashiEsq.reservado != 0 {
					res += colorYellow
				} else {
					res += colorGreen
				}

				res += "Hashi " + fmt.Sprint(f.hashiEsq.id) + "\n"

				//Filosofo
				if f.estado == 0 {
					res += colorGreen
				} else if f.estado == 1 {
					res += colorCyan
				} else if f.estado == 2 {
					res += colorYellow
				} else if f.estado == 3 {
					res += colorPurple
				} else if f.estado == -1 {
					res += colorRed
				} else {
					res += colorReset
				}

				res += "Filosofo " + fmt.Sprint(f.id) + "\n"

			}

			//Hashi repetido
			if !filosofos[0].hashiEsq.disponivel {
				res += colorRed
			} else if filosofos[0].hashiEsq.reservado != 0 {
				res += colorYellow
			} else {
				res += colorGreen
			}

			res += "Hashi " + fmt.Sprint(filosofos[0].hashiEsq.id) + colorReset + "\n"

			//Print de tudo
			fmt.Println(res)

			if !checarFilosofos(filosofos) {
				break
			}
		}
	}

	//Esperar threads terminarem
	//fmt.Println(division, "Fim", division)
	wg.Wait()
}
