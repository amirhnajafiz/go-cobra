package crypto

import (
	"github.com/amirhnajafiz/go-cobra/pkg/encrypt"
	"github.com/amirhnajafiz/go-cobra/pkg/logger"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "crypto",
		Short: "generate certifications",
		Long:  "command for creating certification files in certs directory",
		Run: func(_ *cobra.Command, _ []string) {
			main()
		},
	}
}

func main() {
	l := logger.New()

	enc := encrypt.Encrypt{
		Logger: l.Named("crypto"),
	}

	enc.GenerateCertificateAuthority()
	enc.GenerateCert()
}
