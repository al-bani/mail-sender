package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type EmailRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func main() {
	// Load konfigurasi dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Tidak bisa load .env, lanjut pakai environment variable sistem")
	}
	log.Println("ngetes")

	r := gin.Default()
	// aktifkan CORS default (Allow-All). Untuk produksi, set konfigurasi khusus.
	r.Use(cors.Default())

	r.POST("/send-email", func(c *gin.Context) {
		var req EmailRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

        from := os.Getenv("SMTP_USER")
        pass := os.Getenv("SMTP_PASS")
        // tentukan penerima: SMTP_TO > req.Email > fallback ke from
        to := os.Getenv("SMTP_TO")
        if to == "" {
            if req.Email != "" {
                to = req.Email
            } else {
                to = from
            }
        }
        subject := req.Subject
        body := fmt.Sprintf(
            "From: %s\nEmail: %s\n\n%s",
            req.Name, req.Email, req.Message,
        )

        // header email terurut agar deterministik
        var sb strings.Builder
        sb.WriteString("From: ")
        sb.WriteString(from)
        sb.WriteString("\r\n")

        sb.WriteString("To: ")
        sb.WriteString(to)
        sb.WriteString("\r\n")

        sb.WriteString("Subject: ")
        sb.WriteString(subject)
        sb.WriteString("\r\n")

        sb.WriteString("MIME-Version: 1.0\r\n")
        sb.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
        sb.WriteString("\r\n")
        sb.WriteString(body)

        msg := []byte(sb.String())

        auth := smtp.PlainAuth("", from, pass, os.Getenv("SMTP_HOST"))
        addr := os.Getenv("SMTP_HOST") + ":" + os.Getenv("SMTP_PORT")

        if err := smtp.SendMail(addr, auth, from, []string{to}, msg); err != nil {
			log.Println("Gagal kirim email:", err)
			c.JSON(500, gin.H{"error": "Gagal mengirim email"})
			return
		}

		c.JSON(200, gin.H{"message": "Email berhasil dikirim!"})
	})

	r.Run(":8080")
}
