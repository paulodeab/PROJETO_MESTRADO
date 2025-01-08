package main

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "fmt"
    "math/big"
    "net"
    "os"
    "time"
)

func main() {
    // Gerar chave privada da CA
    caPriv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        fmt.Println("Erro ao gerar chave privada da CA:", err)
        return
    }

    // Configurar o certificado da CA
    caTemplate := &x509.Certificate{
        SerialNumber: big.NewInt(1),
        Subject: pkix.Name{
            Organization: []string{"Minha CA"},
        },
        NotBefore:   time.Now(),
        NotAfter:    time.Now().AddDate(1, 0, 0),
        KeyUsage:    x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
        IsCA:        true,
        BasicConstraintsValid: true,
    }

    // Gerar o certificado da CA
    caCertDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caPriv.PublicKey, caPriv)
    if err != nil {
        fmt.Println("Erro ao criar certificado da CA:", err)
        return
    }

    saveCertAndKey("ca.crt", "ca.key", caCertDER, caPriv)
    fmt.Println("CA gerada com sucesso!")

    // Gerar certificado do servidor
    serverPriv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    serverTemplate := &x509.Certificate{
        SerialNumber: big.NewInt(2),
        Subject: pkix.Name{
            Organization: []string{"Meu Servidor"},
        },
        NotBefore:  time.Now(),
        NotAfter:   time.Now().AddDate(1, 0, 0),
        KeyUsage:   x509.KeyUsageDigitalSignature,
        ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
        IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
        DNSNames:    []string{"localhost"},
    }

    serverCertDER, _ := x509.CreateCertificate(rand.Reader, serverTemplate, caTemplate, &serverPriv.PublicKey, caPriv)
    saveCertAndKey("server.crt", "server.key", serverCertDER, serverPriv)
    fmt.Println("Certificado do servidor gerado com sucesso!")

    // Gerar certificado do cliente
    clientPriv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    clientTemplate := &x509.Certificate{
        SerialNumber: big.NewInt(3),
        Subject: pkix.Name{
            Organization: []string{"Meu Cliente"},
        },
        NotBefore:  time.Now(),
        NotAfter:   time.Now().AddDate(1, 0, 0),
        KeyUsage:   x509.KeyUsageDigitalSignature,
        ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
    }

    clientCertDER, _ := x509.CreateCertificate(rand.Reader, clientTemplate, caTemplate, &clientPriv.PublicKey, caPriv)
    saveCertAndKey("client.crt", "client.key", clientCertDER, clientPriv)
    fmt.Println("Certificado do cliente gerado com sucesso!")
}

func saveCertAndKey(certFile, keyFile string, certDER []byte, privKey *ecdsa.PrivateKey) {
    certOut, _ := os.Create(certFile)
    pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
    certOut.Close()

    keyBytes, _ := x509.MarshalECPrivateKey(privKey)
    keyOut, _ := os.Create(keyFile)
    pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
    keyOut.Close()
}
