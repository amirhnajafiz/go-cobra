package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"
)

func GenerateCertificateAuthority() {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Organization:  []string{"CaptainCore"},
			Country:       []string{"USA"},
			Province:      []string{"PA"},
			Locality:      []string{"Lancaster"},
			StreetAddress: []string{"342 N Queen St"},
			PostalCode:    []string{"17603"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	privacy, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &privacy.PublicKey
	caB, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, privacy)
	if err != nil {
		log.Println("create ca failed", err)
		return
	}

	// Public key
	certOut, err := os.Create("certs/ca.crt")
	_ = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: caB})
	_ = certOut.Close()
	log.Print("written certs/cat.crt\n")

	// Private key
	keyOut, err := os.OpenFile("certs/ca.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	_ = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privacy)})
	keyOut.Close()
	log.Print("written certs/ca.key\n")
}

func GenerateCert() {
	// Load CA
	calts, err := tls.LoadX509KeyPair("certs/ca.crt", "certs/ca.key")
	if err != nil {
		panic(err)
	}
	ca, err := x509.ParseCertificate(calts.Certificate[0])
	if err != nil {
		panic(err)
	}

	// Prepare certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"CaptainCore"},
			Country:       []string{"USA"},
			Province:      []string{"PA"},
			Locality:      []string{"Lancaster"},
			StreetAddress: []string{"342 N Queen St"},
			PostalCode:    []string{"17603"},
			CommonName:    "CaptainCore",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey

	// Sign the certificate
	certB, err := x509.CreateCertificate(rand.Reader, cert, ca, pub, calts.PrivateKey)

	// Public key
	certOut, err := os.Create("certs/cert.pem")
	_ = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certB})
	certOut.Close()
	log.Print("written certs/cert.pem\n")

	// Private key
	keyOut, err := os.OpenFile("certs/key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	_ = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	log.Print("written certs/key.pem\n")
}
