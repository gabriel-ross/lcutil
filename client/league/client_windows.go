package league

import (
	"github.com/abatewongc/bartender-bastion/client"
)

func NewFromExisting() (client.Client, error) {
	return CreateFromWindows()
}
