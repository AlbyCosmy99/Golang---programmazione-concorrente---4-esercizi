package main

import(
	"fmt"
	"time"
	"math/rand"
)

type Operaio struct{
	nome string
	haLavoratoConMartello bool
	haLavoratoConTrapano bool
	haLavoratoConCacciavite bool
}

type Martello struct{
	numero int
}

type Trapano struct{
	numero int
}

type Cacciavite struct{
	numero int
}

func Lavora(operaioRichiesto Operaio,martelli chan Martello,trapani chan Trapano,cacciaviti chan Cacciavite,waitChannel chan bool){
	var operaio Operaio = operaioRichiesto
	
	for{
		if(!operaio.haLavoratoConMartello){
			select{
				case martello:= <-martelli:
					
					var tempoNecessario int = rand.Intn(3) + 1 //il tempo varia casualmente tra 1 e 3 secondi
					fmt.Println("L'operaio ",operaio.nome," ha ottenuto il martello. Inizio dei lavori. Tempo necessario: ",tempoNecessario," secondi.")
					//lavoro in corso...
					time.Sleep(time.Duration(tempoNecessario)*time.Second)
					fmt.Println("L'operaio ",operaio.nome," ha finito di lavorare con il martello.")
					
					martelli <- martello
					
					operaio.haLavoratoConMartello = true
				default:
					//serve affinche' il programma vada avanti e l'operaio chieda un altro attrezzo
			}
		}
		
		if(!operaio.haLavoratoConTrapano){
			select{
				case trapano:= <- trapani:
				
					var tempoNecessario int = rand.Intn(3) + 1 //il tempo varia casualmente tra 1 e 3 secondi
					fmt.Println("L'operaio ",operaio.nome," ha ottenuto il trapano numero ",trapano.numero, ". Inizio dei lavori. Tempo necessario: ",tempoNecessario," secondi.")
					//lavoro in corso...
					time.Sleep(time.Duration(tempoNecessario)*time.Second)
					fmt.Println("L'operaio ",operaio.nome," ha finito di lavorare con il trapano.")
					
					trapani <- trapano
					
					operaio.haLavoratoConTrapano = true
				default:
					//serve affinche' il programma vada avanti e l'operaio chieda un altro attrezzo
			}
		}
		
		if(operaio.haLavoratoConTrapano && !operaio.haLavoratoConCacciavite){
			select{
				case cacciavite:= <-cacciaviti:
				
					var tempoNecessario int = rand.Intn(3) + 1 //il tempo varia casualmente tra 1 e 3 secondi
					fmt.Println("L'operaio ",operaio.nome," ha ottenuto il cacciavite. Inizio dei lavori. Tempo necessario: ",tempoNecessario," secondi.")
					//lavoro in corso...
					time.Sleep(time.Duration(tempoNecessario)*time.Second)
					fmt.Println("L'operaio ",operaio.nome," ha finito di lavorare con il cacciavite.")
					
					cacciaviti <- cacciavite
					
					operaio.haLavoratoConCacciavite = true
				default: 
					//serve affinche' il programma vada avanti e l'operaio chieda un altro attrezzo
			}
		}
		
		if(operaio.haLavoratoConMartello && operaio.haLavoratoConTrapano && operaio.haLavoratoConCacciavite){
			//ha finito di lavorare con tutti gli attrezzi
			break
		}
	}
	waitChannel <- true
	//qui arriva quando l'operaio ha finito tutto il suo lavoro
}

func main(){

	start := (time.Now())
	
	var waitChannel chan bool = make(chan bool) //canale utilizzato per garantire che
	 //il main() aspetti la fine delle altre goroutines
	var timesToWait int = 0
	
	var operai []Operaio = []Operaio{Operaio{nome: "Qui",haLavoratoConMartello: false ,haLavoratoConTrapano: false ,haLavoratoConCacciavite: false},Operaio{nome: "Quo",haLavoratoConMartello: false ,haLavoratoConTrapano: false ,haLavoratoConCacciavite: false},Operaio{nome: "Qua",haLavoratoConMartello: false ,haLavoratoConTrapano: false ,haLavoratoConCacciavite: false}}
	var numeroOperai int = len(operai)
	
	//creo la risorsa Martello con un'istanza
	var martelli chan Martello = make(chan Martello,1)
	martelli <- Martello{numero: 1}
	
	//creo la risorsa Trapano con due istanze
	var trapani chan Trapano = make(chan Trapano,2)
	trapani <- Trapano{numero: 1}
	trapani <- Trapano{numero: 2}
	
	//creo la risorsa Cacciavite con un'istanza
	var cacciaviti chan Cacciavite = make(chan Cacciavite,1)
	cacciaviti <- Cacciavite{numero: 1}
	
	for i:=0;i<numeroOperai;i++{
		timesToWait++
		go Lavora(operai[i],martelli,trapani,cacciaviti,waitChannel)
	}
	
	Wait(timesToWait,waitChannel)
	
	//fine
	fmt.Println("TUTTI E ",numeroOperai," GLI OPERAI HANNO FINITO DI LAVORARE CON I TRE ATTREZZI.")
	
	fmt.Println("TEMPO TOTALE IMPIEGATO NELL'ESECUZIONE: ",time.Since(start))
}

func Wait(times int,waitChannel chan bool){ //funzione che termina quando il canale 'waitChannel' rilascia elementi 'times' volte
	for i:=0;i<times;i++{
		<-waitChannel
	}
}
