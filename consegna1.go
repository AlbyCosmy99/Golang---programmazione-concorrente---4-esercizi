package main

import(
	"fmt"
	"math/rand"
)

type Viaggio struct{
	meta string
	personeAttesa int	 
	clientiMinimiNecessari int 
}

type Cliente struct{
	nome string
	destinazione string
}

func Prenota(cliente chan Cliente,meta chan []Viaggio,waitChannel chan bool){
	
	var clientePrenotato Cliente = <-cliente  //Cliente che sta prenotando
	var viaggi []Viaggio = <-meta   //array con i viaggi disponibili
	
	var indiceMeta int = rand.Intn(len(viaggi)) //avendo due mete, genero un numero intero casuale tra 0 e 1
	
	var viaggioScelto Viaggio = viaggi[indiceMeta] //dopo aver scelto casualmente una meta, la salvo nella variabile viaggioScelto
	
	clientePrenotato.destinazione = viaggioScelto.meta //imposto la meta al Cliente  
	
	viaggioScelto.personeAttesa++ 
	
	cliente <-clientePrenotato  //aggiorno il canale 'cliente'
	
	viaggi[indiceMeta] = viaggioScelto   //aggiorno l'array viaggi 
	meta <-viaggi   //aggiorno il canale 'meta' con il viaggio scelto aggiornato
	
	waitChannel <- true //la goroutine e' conclusa
}

func StampaPartecipanti(clienti []chan Cliente,viaggi chan []Viaggio,waitChannel chan bool){
	
	var lenClienti int = len(clienti) 
	var arrayClienti []Cliente = make([]Cliente,lenClienti) //slice con tutti i Clienti
	for k:=0;k<lenClienti;k++{
		arrayClienti[k] = <-clienti[k]   //riempio lo slice di 'Cliente' con i tipi Cliente del canale 'clienti'
	}
	
	var arrayViaggi []Viaggio = <-viaggi  //array con tutti i tipi Viaggio
	
	for i:=0;i<len(arrayViaggi);i++{  //iterazione per ogni viaggio
		clientiPrenotati := ""   
		
		for j:=0;j<len(arrayClienti);j++{   //iterazione per ogni cliente
			if arrayClienti[j].destinazione == arrayViaggi[i].meta{
				clientiPrenotati = clientiPrenotati + arrayClienti[j].nome + "; "
			}
		}
		
		fmt.Println(); //linea vuota
		
		if clientiPrenotati==""{
			fmt.Println("Per il viaggio in ",arrayViaggi[i].meta," non ci sono state prenotazioni.")
		}else{
			fmt.Print("Per il viaggio in ",arrayViaggi[i].meta," hanno prenotato i seguenti clienti: ")
			fmt.Println(clientiPrenotati)
		}
		
		fmt.Println("Il viaggio in ",arrayViaggi[i].meta, " ha bisogno di almeno ",arrayViaggi[i].clientiMinimiNecessari," clienti per essere confermato.")
		
		if arrayViaggi[i].personeAttesa >= arrayViaggi[i].clientiMinimiNecessari{
			fmt.Println("Ci sono ",arrayViaggi[i].personeAttesa, " persone in attesa, quindi il viaggio e' stato confermato.")
		}else{
			fmt.Println("Ci sono ",arrayViaggi[i].personeAttesa, " persone in attesa, quindi il viaggio NON e' stato confermato.")
		}
	}	
	
	waitChannel <- true  //la goroutine e' conclusa
}

func main(){

	var waitChannel chan bool = make(chan bool) //canale che garantisce che il main() non termini prima della terminazione delle altre goroutines
	var timesToWait int = 0 //numero delle goroutines totali tranne il main()
	
	const clientiTotali int = 7
	const viaggiTotali int = 2
	
	clienti := []Cliente{Cliente{nome: "Brontolo"},Cliente{nome: "Cucciolo"},Cliente{nome: "Dotto"},Cliente{nome: "Eolo"},Cliente{nome: "Gongolo"},Cliente{nome: "Mammolo"},Cliente{nome: "Pisolo"}}
	viaggi := []Viaggio{Viaggio{meta: "Spagna",personeAttesa: 0,clientiMinimiNecessari: 4},Viaggio{meta: "Francia",personeAttesa: 0,clientiMinimiNecessari:2}}
	
	var canaliClienti []chan Cliente = []chan Cliente{make(chan Cliente,1),make(chan Cliente,1),make(chan Cliente,1),make(chan Cliente,1),make(chan Cliente,1),make(chan Cliente,1),make(chan Cliente,1)}  //array di canali contenenti ciascuno un Cliente
	
	for j:=0;j<len(canaliClienti);j++{
		canaliClienti[j] <- clienti[j]  //aggiorno i canali dell'array canaliClienti
	}
	
	var chanViaggi chan []Viaggio = make(chan []Viaggio,1) //canale contenente un array di tipi Viaggio
	
	chanViaggi <-viaggi  //aggiorno il canale chanViaggi con i viaggi disponibili           
	
	for i:=0;i<len(clienti);i++{  //per ciascun cliente
		timesToWait++  //quante goroutines ci sono in tutto e' salvato in 'timesToWait'
		go Prenota(canaliClienti[i],chanViaggi,waitChannel)
	}
	
	Wait(timesToWait,waitChannel)  //si prosegue se il canale waitChannel rilascia 'timesToWait' elementi, ossia se sono
		//terminati tutte le goroutines previste
	
	go StampaPartecipanti(canaliClienti,chanViaggi,waitChannel)
	
	Wait(1,waitChannel)  //si assicura che la goroutine 'StampaPartecipanti(....)' termini prima del main()
}

func Wait(times int,waitChannel chan bool){  //funzione che termina quando il canale 'waitChannel' rilascia elementi 'times' volte
	for i:=0;i<times;i++{
		<-waitChannel
	}
}
