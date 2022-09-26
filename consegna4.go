//ho aggiunto commenti e stampe che facciano capire cosa sta accadendo durante l'esecuzione
//ho modificato il tunnelChannel := make(chan Tunnel,1) aggiungendo il ',1' risolvendo il deadlock iniziale
//ho modificato il nome della funzione 'transumanza' in 'LanciaPalline', per maggiore chiarezza
/*ho tolto il 'time.Sleep(time.Minute)' sostituendolo con la 
funzione 'Wait(times int,waitChannel chan bool)', introducendo la variabile
'var waitChannel chan bool = make(chan bool)'. Questo affinche' la funzione 'main()' 
termini non appena terminano tutte le altre goroutines, senza aspettare per forza un 
minuto totale. La varabile 'waitChannel' e' passata come argomento alla funzione 
'LanciaPalline(..)', all'interno della quale viene modificata una volta che la funzione 
stessa ha terminato il suo lavoro. Per tenere il conto di tutte le goroutines eccetto
la funzione 'main()' ho introdotto la variabile 'timesToWait'.*/
//ho modificato il nome della funzione 'mandaPersona(..)' in 'TentaUnLancio(..)' per maggiore chiarezza
/*introdotta variabile 'lancioGruppoXInCorso' per ogni gruppo per fare in modo che
se in un certo gruppo una pallina e' in movimento, non vengano lanciate palline
dello stesso gruppo. Infatti nella consegna c'e' scritto che "Una persona NON PUO' 
lanciare una pallina finché quella che ha lanciato precedentemente non è.....".
Attraverso la variabile booleana 'lancioGruppoXInCorso' il problema del lancio
contemporaneo di piu' palline dello stesso gruppo e' risolto.*/
/*Eliminata la variabile 'c1 := make(chan int,1)' e il suo utilizzo. Il problema 
del lancio contemporaneo di piu' palline dello stesso gruppo e' stato risolto
come discusso nel commento sopra. Mancava solamente la gestione degli scontri tra
palline di due gruppi diversi. La logica applicata e' la seguente: se non si stanno
lanciando palline dello stesso gruppo e IL TUNNEL E' LIBERO, si lancia la pallina
in questione. Se in seguito una pallina dell'altro gruppo viene lanciata e si nota
che il tunnel non e' libero, significa che si e' verificato uno scontro. La seconda
pallina lanciata verra' semplicemente fatta ripartire dopo aver impostato correttamente
le variabili booleane e dopo aver liberato i canali necessari. La prima pallina che 
era stata lanciata, prima di completare il viaggio (es. decrementare g.nPalline), si
verifica che effettivamente il TUNNEL SIA ANCORA LIBERO. In caso non lo sia piu'
significa che un'altra pallina lo ha liberato perche' si era verificato uno scontro
con la stessa.*/

package main

import "fmt"
import "time"
import "math/rand"
import "strings"

type Gruppo struct {
    nome string
    nPalline int
}

type Tunnel struct {
    libero bool
}

func LanciaPalline(g Gruppo, t chan Tunnel,lancioStessoGruppoInCorso *bool,waitChannel chan bool){
    for g.nPalline > 0{
		var attesa int = rand.Intn(2)
		fmt.Println("La routine del gruppo ",g.nome," sta attendendo ",attesa," secondi prima di provare a lanciare una nuova pallina.")
        time.Sleep(time.Duration(attesa)*time.Second)
		fmt.Println("La routine del gruppo ",g.nome," sta per provare a lanciare una nuova pallina.")
        TentaUnLancio(&g,t,lancioStessoGruppoInCorso)
    }
	waitChannel <- true 
}

func TentaUnLancio(g *Gruppo, t chan Tunnel,lancioStessoGruppoInCorso *bool){
    if(!*(lancioStessoGruppoInCorso)){
		*lancioStessoGruppoInCorso = true
		tunnel := <- t
		if tunnel.libero {
			tunnel.libero = false
			t <- tunnel
			
			//la pallina viene lanciata. Tempo necessario: 1 secondo
			time.Sleep(time.Second)
			
			tunnel := <- t
			if(!tunnel.libero){  //non ci sono stati scontri,quindi il tunnel nel frattempo non e' stato liberato
				tunnel.libero = true
				t <- tunnel
				g.nPalline = g.nPalline - 1
				*lancioStessoGruppoInCorso = false //lancio della pallina concluso con successo
				fmt.Println("LANCIO PALLINA GRUPPO ",strings.ToUpper(g.nome)," AVVENUTO CORRETTAMENTE. RIMANGONO ", g.nPalline, " PALLINE NEL GRUPPO ", strings.ToUpper(g.nome),".")
			}else{  //si e' verificato uno scontro 
				t <- tunnel
				*lancioStessoGruppoInCorso = false //lancio della pallina concluso con uno scontro
	
			}
			
		} else{
			//sto provando a lanciare una pallina mentre il tunnel non e' libero, 
			//quindi una pallina dell'altro gruppo e' in movimento nel tunnel.
			//Di conseguenza c'e' uno scontro e entrambe le palline tornano indietro
			//per essere rilanciate
			tunnel.libero = true
			t <- tunnel
			*lancioStessoGruppoInCorso = false //lancio della pallina concluso con uno scontro
			fmt.Println("Scontro tra due palline dei due gruppi. Entrambe ritornano indietro.")
		}
	}
}

func main() {
	
	start := (time.Now())
	
    rand.Seed(time.Now().UnixNano())
	
	var waitChannel chan bool = make(chan bool)
	var timesToWait int = 0 //numero delle goroutines tranne la funzione main()
	
	//gruppi di palline
    gruppo1 := Gruppo{"destra",5}
    gruppo2 := Gruppo{"sinistra",5}
	
	//Una pallina di un gruppo non puo' essere lanciata se e' 
	//in corso il lancio di una pallina della stesso gruppo
	var lancioGruppo1InCorso bool = false
	var lancioGruppo2InCorso bool = false
    
	//creazione risorsa Tunnel. Inizialmente e' libero.
    tunnelChannel := make(chan Tunnel,1)  //ho modificato aggiungendo il ,1; ossia l ho fatto diventare un buffered channel
	//l'inserire e il prelevare un'istanza della risorsa Tunnel non sono per forza simultanee
    tunnel := Tunnel{true}
    tunnelChannel <- tunnel
    
	fmt.Println("Il gruppo di palline '",gruppo1.nome,"' inizia a essere lanciato. Palline totali nel gruppo: ",gruppo1.nPalline) 
	timesToWait++
    go LanciaPalline(gruppo1,tunnelChannel,&lancioGruppo1InCorso,waitChannel)
	timesToWait++
	fmt.Println("Il gruppo di palline '",gruppo2.nome,"' inizia a essere lanciato. Palline totali nel gruppo: ",gruppo2.nPalline)
    go LanciaPalline(gruppo2,tunnelChannel,&lancioGruppo2InCorso,waitChannel)
    
	//si assicura che il main() non termini troppo velocemente per permettere alle goroutines di terminare
	Wait(timesToWait,waitChannel)
	
	fmt.Println("TUTTE LE ",gruppo1.nPalline + gruppo2.nPalline," PALLINE SONO STATE LANCIATE ATTRAVERSO IL TUNNEL CORRETTAMENTE.")
	
	fmt.Println("TEMPO TOTALE IMPIEGATO NELL'ESECUZIONE: ",time.Since(start))
}

func Wait(times int,waitChannel chan bool){
	for i:=0;i<times;i++{
		<-waitChannel
	}
}
