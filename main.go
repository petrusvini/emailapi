package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/smtp"
)

// Estrutura para a mensagem recebida do site
type Message struct {
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

// Configurações do servidor SMTP (Gmail)
const (
    smtpHost     = "smtp.gmail.com"
    smtpPort     = "587"
    senderEmail  = "petrusvini.ar@gmail.com" // Seu e-mail do Gmail
    senderPass   = "qgml ysbs wncn kgqg"         // Sua Senha de App
    recipient    = "petrusvini.ar@gmail.com"  // E-mail que receberá a mensagem
)

// Função para enviar e-mail
func sendEmail(subject, body string) error {
    msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", recipient, subject, body)

    auth := smtp.PlainAuth("", senderEmail, senderPass, smtpHost)
    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{recipient}, []byte(msg))
    if err != nil {
        return fmt.Errorf("falha ao enviar e-mail: %v", err)
    }
    return nil
}

// Handler para a rota /send-message
func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

    var msg Message
    if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
        http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
        return
    }

    // Enviar e-mail
    if err := sendEmail(msg.Subject, msg.Body); err != nil {
        http.Error(w, fmt.Sprintf("Erro ao enviar e-mail: %v", err), http.StatusInternalServerError)
        return
    }

    // Resposta de sucesso
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Mensagem enviada com sucesso!")
}

func main() {
    // Definir a rota
    http.HandleFunc("/send-message", sendMessageHandler)

    // Iniciar o servidor na porta 8080
    log.Println("Servidor rodando em http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Erro ao iniciar servidor: %v", err)
    }
}