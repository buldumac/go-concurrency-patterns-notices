package t

import (
	"context"
)

// Linux users, admins, heroes might have a Déjà vu

// echo HelloGo | tee 1.txt 2.txt 3.txt
// or
// echo HelloGo | tee -a 1.txt 2.txt 3.txt

// There is a pattern in Go which name is Tee Channel that is doing the same
// This pattern have a lot in common with Fan-Out pattern but the difference is that Tee Channel is not splitting,
// it's sending the same data

// TeePlease - same as in Linux: echo Hello | tee 1 2 3
func TeePlease[T interface{}](ctx context.Context, sourceCh <-chan T) (_, _, _ <-chan T) {
	outCh1 := make(chan T)
	outCh2 := make(chan T)
	outCh3 := make(chan T)

	go func() {
		defer func() {
			close(outCh1)
			close(outCh2)
			close(outCh3)
		}()

		for sourceValue := range sourceCh {
			outCh1, outCh2, outCh3 := outCh1, outCh2, outCh3
			for i := 0; i < 3; i++ {
				select {
				case <-ctx.Done():
					return
				case outCh1 <- sourceValue:
					outCh1 = nil
				case outCh2 <- sourceValue:
					outCh2 = nil
				case outCh3 <- sourceValue:
					outCh3 = nil
				}
			}
		}
	}()

	return outCh1, outCh2, outCh3
}
