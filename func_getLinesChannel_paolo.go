package main

import (
	"io"
	"log"
	"strings"
)

// getLinesChannel legge i dati dal file f in una goroutine,
// accumula i byte per ricostruire le linee complete e le invia
// a un canale di stringhe.
//
// Quando la lettura è completata, chiude il file e il canale.
func getLinesChannel_paolo(f io.ReadCloser) <-chan string {
	// 2. Crea un canale di stringhe
	linesChannel := make(chan string)

	// Avvia una goroutine per la lettura asincrona del file
	go func() {
		// 7. Chiude il file quando la goroutine termina (defer)
		defer f.Close()
		// 6. Chiude il canale quando la goroutine termina (defer)
		defer close(linesChannel)

		// Buffer di 8 byte per la lettura
		buffer := make([]byte, 8)
		persistentString := "" // Accumula i byte per tenere traccia della linea corrente

		for {
			nread, err := f.Read(buffer)

			// 8. Gestione degli errori (incluso EOF)
			if err != nil && err != io.EOF {
				// In caso di errore non-EOF (es. I/O error), lo logghiamo e usciamo dalla goroutine.
				log.Printf("Errore di lettura del file: %v", err)
				return
			}

			if nread > 0 {
				// Aggiunge i byte letti alla stringa persistente
				persistentString += string(buffer[:nread])

				// Logica di suddivisione (split) delle linee
				sections := strings.Split(persistentString, "\n")

				// Itera su tutte le sezioni tranne l'ultima (che è la parte incompleta della linea)
				for index, section := range sections {
					if index < len(sections)-1 {
						// 5. Invia la linea completa al canale
						linesChannel <- section
					} else {
						// L'ultima sezione (la parte incompleta) diventa la nuova stringa persistente
						persistentString = section
					}
				}
			}

			// Se abbiamo raggiunto la fine del file (EOF), usciamo dal ciclo
			if err == io.EOF {
				// Dopo che il ciclo è terminato, gestiamo l'ultima linea
				// rimasta in persistentString che non termina con un newline.
				if persistentString != "" {
					linesChannel <- persistentString
				}
				break
			}
		}
	}()

	// 1. Ritorna il canale in sola lettura immediatamente
	return linesChannel
}


