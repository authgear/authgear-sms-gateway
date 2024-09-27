package sensitive

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

type PhoneNumber string

func (s PhoneNumber) String() string {
	num, err := phonenumbers.Parse(string(s), "")
	if err != nil {
		// Not a phone number, so do not mask it.
		return string(s)
	}

	countryCallingCodeWithoutPlusSign := strconv.Itoa(int(num.GetCountryCode()))
	nationalNumber := phonenumbers.GetNationalSignificantNumber(num)

	startIndex := 0
	endIndex := len(nationalNumber) - 3
	// Then do not mask
	if endIndex < 0 {
		endIndex = 0
	}

	masked := nationalNumber[startIndex:endIndex]
	unmasked := nationalNumber[endIndex:]

	masked = strings.Repeat("*", len(masked))

	return fmt.Sprintf("+%v%v%v", countryCallingCodeWithoutPlusSign, masked, unmasked)
}
