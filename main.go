package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// ชี้ไปที่โฟลเดอร์ static เพื่อเสิร์ฟไฟล์ HTML
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	fmt.Println("Server is running on port 8080...")
	// รันเซิร์ฟเวอร์ที่ Port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
