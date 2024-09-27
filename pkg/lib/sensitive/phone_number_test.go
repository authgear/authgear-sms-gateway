package sensitive

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPhoneNumber(t *testing.T) {
	Convey("PhoneNumber", t, func() {
		test := func(s string, expected string) {
			phoneNumber := PhoneNumber(s)
			actual := phoneNumber.String()
			So(actual, ShouldEqual, expected)
		}

		test("", "")
		test("+85298765432", "+852*****432")
		test("+852987", "+852987")
		test("+8521823", "+852*823")
		test("+852999", "+852999")
	})
}
