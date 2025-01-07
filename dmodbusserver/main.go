package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/tbrandon/mbserver"
)

func main() {
    // Cria um novo servidor Modbus
    server := mbserver.NewServer()

    // Inicia o servidor Modbus TCP
    err := server.ListenTCP("0.0.0.0:502")
    if err != nil {
        log.Fatalf("Erro ao iniciar o servidor: %v", err)
    }
    defer server.Close()

    // Simula dados de temperatura
    go func() {
        for {
            temperature := uint16(rand.Intn(100)) // Temperatura simulada entre 0 e 99
			fmt.Printf("Temp: %dC°\n", temperature)
            server.HoldingRegisters[0] = temperature
            time.Sleep(1 * time.Second)
        }
    }()

    // Mantém o servidor em execução
    select {}
}