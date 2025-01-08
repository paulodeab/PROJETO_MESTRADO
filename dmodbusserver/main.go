package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/tbrandon/mbserver"
)

func main() {
    // Carregar certificados TLS
    cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
    if err != nil {
        log.Fatalf("Erro ao carregar certificados: %v", err)
    }

    // Configuração TLS
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
    }

    // Cria um novo servidor Modbus
    server := mbserver.NewServer()

    // Inicia o servidor Modbus TCP com TLS
    err = server.ListenTLS("0.0.0.0:502", tlsConfig)
    if err != nil {
        log.Fatalf("Erro ao iniciar o servidor: %v", err)
    }
    defer server.Close()

    // Simula dados de temperatura
    go func() {
        for {
            temperature := uint16(rand.Intn(100)) // Temperatura simulada entre 0 e 99
            server.HoldingRegisters[0] = temperature
            time.Sleep(1 * time.Second)
            fmt.Println(temperature)
        }
    }()

    // Mantém o servidor em execução
    select {}
}