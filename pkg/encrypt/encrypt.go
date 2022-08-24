package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"

	"go.uber.org/zap"
)

type Encrypt struct {
	Logger *zap.Logger
}

func (e Encrypt) GenerateCertificateAuthority() {
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

	caB, err := x509.CreateCertificate(rand.Reader, ca, ca, &privacy.PublicKey, privacy)
	if err != nil {
		e.Logger.Error("create ca failed : " + err.Error())

		return
	}

	_ = os.Mkdir("certs", 0600)

	// Public key
	certOut, _ := os.Create("certs/ca.crt")

	_ = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: caB})
	_ = certOut.Close()

	e.Logger.Info("written certs/cat.crt\n")

	// Private key
	keyOut, _ := os.OpenFile("certs/ca.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)

	_ = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privacy)})
	_ = keyOut.Close()

	e.Logger.Info("written certs/ca.key\n")
}

func (e Encrypt) GenerateCert() {
	// Load CA
	cults, err := tls.LoadX509KeyPair("certs/ca.crt", "certs/ca.key")
	if err != nil {
		e.Logger.Panic(err.Error())
	}

	ca, err := x509.ParseCertificate(cults.Certificate[0])
	if err != nil {
		e.Logger.Panic(err.Error())
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

	privy, _ := rsa.GenerateKey(rand.Reader, 2048)

	// Sign the certificate
	certB, err := x509.CreateCertificate(rand.Reader, cert, ca, &privy.PublicKey, cults.PrivateKey)

	// Public key
	certOut, err := os.Create("certs/cert.pem")

	_ = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certB})
	_ = certOut.Close()

	e.Logger.Info("written certs/cert.pem\n")

	// Private key
	keyOut, err := os.OpenFile("certs/key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)

	_ = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privy)})
	_ = keyOut.Close()

	e.Logger.Info("written certs/key.pem\n")
}
