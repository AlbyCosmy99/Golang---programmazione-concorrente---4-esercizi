/*A schermo vengono stampate le varie fasi durante l'esecuzione per
comprendere cosa sta accadendo. Per ogni piatto le fasi descritte a schermo sono:
-Piatto ordinato: X.
-Il piatto X e' stato assegnato al fornello numero Y.
-Cucinando il piatto: X. Tempo necessario: Z secondi.
-Piatto cucinato: X
-Piatto X assegnato al cameriere: K
-Facendo uscire dalla cucina il piatto: X. Tempo necessario: N secondi.
-Il piatto X e' stato ordinato, cucinato e portato con successo al cliente: W
*/

package main

import(
	"fmt"
	"math/rand"
	"time"
)

type Piatto struct{
	nomeCibo string
}

type Fornello struct{
	numero int
}

type Cameriere struct{
	nome string
	piatto Piatto
}

type Cliente struct{
	nome string
	piatto Piatto
}

type Pair struct{  //coppia piatto-booleana 
	piatto Piatto
	booleana bool
}

func Ordina(ordinazioneEffettuata chan Pair,waitChannel chan bool){
	piattoOrdinato := (<-ordinazioneEffettuata).piatto
	fmt.Println("Piatto ordinato: ",piattoOrdinato.nomeCibo)
	ordinazioneEffettuata <-Pair{piatto: piattoOrdinato,booleana: true} //piatto ordinato? true
	waitChannel<-true
}

func Cucina(ordinazioneEffettuata chan Pair,piattoCucinato chan Pair,fornelli chan Fornello,waitChannel chan bool){
	for{
		pair := <-ordinazioneEffettuata
		if(pair.booleana==true){
			fornelloUtilizzato := <-fornelli //richiede un fornello per proseguire
			
			piattoInCucina := (<-piattoCucinato).piatto
			
			fmt.Println("Piatto ",piattoInCucina.nomeCibo," assegnato al fornello numero: ",fornelloUtilizzato.numero)
			
			var tempoImpiegato int = rand.Intn(3) + 4 //il tempo impiegato per cucinare il piatto varia da 4 a 6 sec
												  //lo genero casualmente
												  
			fmt.Println("Cucinando il piatto: ",piattoInCucina.nomeCibo,". Tempo necessario: ",tempoImpiegato," secondi.")
	
			time.Sleep(time.Duration(tempoImpiegato) * time.Second)
			
			fornelli<-fornelloUtilizzato  //rilascio il fornello
			
			fmt.Println("Piatto cucinato: ",piattoInCucina.nomeCibo)
			
			piattoCucinato <- Pair{piatto: piattoInCucina, booleana: true}  //il piatto e'stato cucinato
	
			waitChannel<-true 
			
			break
			
		}else{  //ordinazione non ancora effettuata
			ordinazioneEffettuata <- pair
		}
	}	
}

func UscitaPiatto(piattoCucinato chan Pair,camerieri chan Cameriere,cliente Cliente,waitChannel chan bool){
	for{
		pair := <-piattoCucinato
		if(pair.booleana==true){
			cameriereDesignato := <-camerieri //richiede un cameriere per proseguire
			
			piattoDaTrasportare := pair.piatto
			
			cameriereDesignato.piatto = piattoDaTrasportare
			
			fmt.Println("Piatto ",cameriereDesignato.piatto.nomeCibo," assegnato al cameriere: ",cameriereDesignato.nome)
			
			var tempoImpiegato int = 3
												  
			fmt.Println("Facendo uscire dalla cucina il piatto: ",cameriereDesignato.piatto.nomeCibo,". Tempo necessario: ",tempoImpiegato," secondi.")
	
			time.Sleep(time.Duration(tempoImpiegato) * time.Second)
			
			fmt.Println("Il piatto ",cameriereDesignato.piatto.nomeCibo, " e' stato ordinato, cucinato e portato con successo al cliente: ",cliente.nome)
			
			cameriereDesignato.piatto = Piatto{}
			
			camerieri<-cameriereDesignato  //rilascio il cameriere
	
			waitChannel<-true 
			
			break
			
		}else{  //piatto non ancora cucinato
			piattoCucinato <- pair
		}
	}
}

func main(){
	
	start := (time.Now())
	
	clienti:= []Cliente{Cliente{nome : "Luca",piatto: Piatto{nomeCibo: "Funghi trifolati"}},Cliente{nome: "Elena",piatto: Piatto{nomeCibo: "Salsicce e patate al forno"}},Cliente{nome: "Salif",piatto: Piatto{nomeCibo: "Lasagna"}},Cliente{nome: "Tommaso",piatto: Piatto{nomeCibo: "Pollo e uova"}},Cliente{nome: "Genoveffa",piatto: Piatto{nomeCibo: "Moussaka'"}},Cliente{nome: "Alex",piatto: Piatto{nomeCibo: "Pasta e fagioli"}},Cliente{nome: "Marius",piatto: Piatto{nomeCibo: "Pizza Margherita"}},Cliente{nome: "Davide",piatto: Piatto{nomeCibo: "Carne e riso"}},Cliente{nome: "Lucia",piatto: Piatto{nomeCibo: "Pizza Diavola"}},Cliente{nome: "Marco",piatto: Piatto{nomeCibo: "Crepes"}}}
	var lenClienti int = len(clienti)
	
	var waitChannel chan bool = make(chan bool) //canale che garantisce che il main() non termini prima della terminazione delle altre goroutines
	var timesToWait int = 0
	
	var ordinazioniEffettuate []chan Pair = make([]chan Pair,lenClienti)
	
	var piattiCucinati []chan Pair = make([]chan Pair,lenClienti)
	
	//creo un canale per i fornelli e li inizializzo con tre fornelli
	var fornelli chan Fornello = make(chan Fornello,3)
	fornelli<-Fornello{numero: 1}
	fornelli<-Fornello{numero: 2}
	fornelli<-Fornello{numero: 3}
	
	var camerieri chan Cameriere = make(chan Cameriere,2)
	camerieri<-Cameriere{nome: "Andrei"}
	camerieri<-Cameriere{nome: "Daniel"}
	
	for i:=0;i<lenClienti;i++{
		timesToWait++
		ordinazioniEffettuate[i] = make(chan Pair,1)
		ordinazioniEffettuate[i] <- Pair{piatto: clienti[i].piatto,booleana: false} //ordinazioniEffettuata? false
		go Ordina(ordinazioniEffettuate[i],waitChannel) 
	}
	
	for i:=0;i<lenClienti;i++{
		timesToWait++
		piattiCucinati[i] = make(chan Pair,1)
		piattiCucinati[i] <- Pair{piatto: clienti[i].piatto,booleana: false} //piatto cucinato? false
		go Cucina(ordinazioniEffettuate[i],piattiCucinati[i],fornelli,waitChannel)
	}
	
	for i:=0;i<lenClienti;i++{
		timesToWait++
		go UscitaPiatto(piattiCucinati[i],camerieri,clienti[i],waitChannel) //se il piatto e' stato cucinato, esce dalla cucina con il cameriere
	}
	
	Wait(timesToWait,waitChannel)  //funziona che controlla che il main non termini prima di tutte le altre goroutines
	
	fmt.Println("TUTTE LE RICHIESTE DEI ",lenClienti," CLIENTI SONO STATE COMPLETATE CON SUCCESSO.")
	
	fmt.Println("TEMPO TOTALE IMPIEGATO NELL'ESECUZIONE: ",time.Since(start))
}

func Wait(times int,waitChannel chan bool){  //funzione che termina quando il canale 'waitChannel' rilascia elementi 'times' volte
	for i:=0;i<times;i++{
		<-waitChannel
	}
}
