package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/smtp"
    "os"

    "github.com/rs/cors" // Adicione esta linha
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
    // Criar um novo mux (roteador)
    mux := http.NewServeMux()
    mux.HandleFunc("/send-message", sendMessageHandler)

    // Configurar o middleware CORS
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"}, // Permite todas as origens (ajuste para maior segurança em produção)
        AllowedMethods:   []string{"POST"}, // Permite apenas POST
        AllowedHeaders:   []string{"Content-Type"}, // Permite o header Content-Type
        AllowCredentials: false,
    })

    // Aplicar o middleware CORS ao roteador
    handler := c.Handler(mux)

    // Iniciar o servidor na porta fornecida pelo Render (ou 8080 localmente)
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Porta padrão para testes locais
    }
    log.Println("Servidor rodando na porta", port)
    if err := http.ListenAndServe(":"+port, handler); err != nil {
        log.Fatalf("Erro ao iniciar servidor: %v", err)
    }
}