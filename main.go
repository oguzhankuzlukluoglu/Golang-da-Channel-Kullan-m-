package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// İki kanal oluştur
	messageChannel := make(chan string)
	doneChannel := make(chan bool)

	// Wait group oluştur
	var wg sync.WaitGroup

	// Gönderici fonksiyonu
	sender := func(messages []string) {
		// defer keyi önüne geldiği işlemin en son gerçekleştirilmesini sağlar
		defer close(messageChannel)
		defer wg.Done()

		for _, msg := range messages {
			// Mesajı kanala gönder
			messageChannel <- msg
			time.Sleep(time.Millisecond * 500) // Zaman gecikmesi eklendi
		}
	}

	// Alıcı fonksiyonu
	receiver := func() {
		defer wg.Done()

		for {
			select {
			case msg, ok := <-messageChannel:
				if !ok {
					// Kanal kapandıysa işlemi sonlandır
					fmt.Println("Alıcı: Kanal kapandı.")
					doneChannel <- true
					return
				}
				// Mesajı yazdır
				fmt.Println("Alındı:", msg)
			}
		}
	}

	// WaitGroup'a gönderici ve alıcı işlemleri ekleniyor
	wg.Add(2)
	messages := []string{"Merhaba", "Go", "Dilinde", "Kanal", "Kullanımı"}
	go sender(messages)
	go receiver()

	// Gönderici işlem tamamlandığında, kanalı kapat ve alıcıyı beklet
	go func() {
		wg.Wait()
		close(doneChannel)
	}()

	// Alıcı işlem tamamlandığında programı sonlandır
	<-doneChannel
}
