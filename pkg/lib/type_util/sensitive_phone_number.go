package type_util

import (
	"fmt"
)

type SensitivePhoneNumber string

func (s SensitivePhoneNumber) String() string {
	maxLenToMask := 6
	lenToMask := max(len(string(s))-maxLenToMask, 0)
	return fmt.Sprintf("%s%s", string(s)[0:lenToMask], " *** ")
}
