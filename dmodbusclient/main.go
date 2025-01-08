package main

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io/ioutil"
    "time"

    "github.com/simonvetter/modbus"
)

func main() {
    // Carregar certificado da CA
    caCert, err := ioutil.ReadFile("ca.crt")
    if err != nil {
        fmt.Println("Erro ao ler certificado da CA:", err)
        return
    }

    caCertPool := x509.NewCertPool()
    if !caCertPool.AppendCertsFromPEM(caCert) {
        fmt.Println("Erro ao adicionar certificado da CA ao pool")
        return
    }

    // Carregar certificado do cliente
    clientCert, err := tls.LoadX509KeyPair("client.crt", "client.key")
    if err != nil {
        fmt.Println("Erro ao carregar certificado do cliente:", err)
        return
    }

    // Configuração do cliente Modbus TCP com TLS
    client, err := modbus.NewClient(&modbus.ClientConfiguration{
        URL:           "tcp+tls://127.0.0.1:502",
        Timeout:       10 * time.Second,
        TLSClientCert: &clientCert,
        TLSRootCAs:    caCertPool,
    })
    if err != nil {
        fmt.Println("Erro ao criar cliente Modbus:", err)
        return
    }

    // Conectar ao servidor
    err = client.Open()
    if err != nil {
        fmt.Println("Erro ao conectar:", err)
        if _, ok := err.(*tls.CertificateVerificationError); ok {
            fmt.Println("Conexão rejeitada devido a certificado inválido.")
        }
        return
    }
    defer client.Close()

    fmt.Println("Conexão estabelecida com sucesso!")

    // Loop para leitura de registros
    for {
        // Leitura de registros de temperatura
        results, err := client.ReadRegisters(0, 1, modbus.HOLDING_REGISTER)
        if err != nil {
            fmt.Println("Erro ao ler registros:", err)
            return
        }

        temperature := results[0]
        fmt.Printf("Temperatura: %d°C\n", temperature)
        time.Sleep(1 * time.Second)
    }
}
