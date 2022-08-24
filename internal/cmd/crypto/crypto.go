package crypto

import (
	"github.com/amirhnajafiz/go-cobra/pkg/encrypt"
	"github.com/amirhnajafiz/go-cobra/pkg/logger"
)

func main() {
	l := logger.New()

	enc := encrypt.Encrypt{
		Logger: l.Named("crypto"),
	}

	enc.GenerateCertificateAuthority()
	enc.GenerateCert()
}
