package internal

import (
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
)

type Stats struct {
	Mutex        sync.Mutex
	NbDonations  int64
	NbSuccessful int64
	TotalAmount  *big.Int
	TotalFaulty  *big.Int
}

func (s Stats) String() string {
	currency := strings.ToUpper(os.Getenv("CURRENCY"))
	total := big.NewInt(0).Add(s.TotalAmount, s.TotalFaulty)
	average := big.NewInt(0)
	if s.NbSuccessful > 0 {
		// avoid division by 0 and panic
		average.Div(s.TotalAmount, big.NewInt(s.NbSuccessful))
	}
	line1 := fmt.Sprintf("total received: %s\t%s\n", currency, total.String())
	line2 := fmt.Sprintf("successfully donated: %s\t%s\n", currency, s.TotalAmount.String())
	line3 := fmt.Sprintf("faulty donation: %s\t%s", currency, s.TotalFaulty.String())
	line4 := fmt.Sprintf("average per person: %s\t%s", currency, average.String())
	return fmt.Sprintf("%s%s%s\n%s\n", line1, line2, line3, line4)
}
