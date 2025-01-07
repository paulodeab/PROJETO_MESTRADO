package main

import (
    "fmt"
    "time"

    "github.com/goburrow/modbus"
)

func main() {
    handler := modbus.NewTCPClientHandler("localhost:502")
    handler.Timeout = 10 * time.Second
    handler.SlaveId = 1

    err := handler.Connect()
    if err != nil {
        fmt.Println("Erro ao conectar:", err)
        return
    }
    defer handler.Close()

    client := modbus.NewClient(handler)

    for {
        results, err := client.ReadHoldingRegisters(0, 1)
        if err != nil {
            fmt.Println("Erro ao ler registros:", err)
            return
        }
		fmt.Print(results)

        temperature := results[1]
        fmt.Printf("Temperatura: %d°C\n", temperature)
        time.Sleep(1 * time.Second)
    }
}