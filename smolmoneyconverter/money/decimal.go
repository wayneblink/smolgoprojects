package money

import (
	"fmt"
	"strconv"
	"strings"
)

// Decimal can represent a floating-point number with a fixed precision.
// example:
// 1.52 = 152 * 10^(-2) will be stored as {152, 2}
type Decimal struct {
	// multiply it by the precision to get the real value
	subunits int64
	// number of "subunits" in a unit, expressed as a power of 10.
	precision byte
}

func (d *Decimal) String() string {
	if d.precision == 0 {
		return fmt.Sprintf("%d", d.subunits)
	}

	centsPerUnit := pow10(d.precision)
	frac := d.subunits % centsPerUnit
	integer := d.subunits / centsPerUnit

	decimalFormat := "%d.%0" + strconv.Itoa(int(d.precision)) + "d"

	return fmt.Sprintf(decimalFormat, integer, frac)
}

const maxDecimal = 1e12

func ParseDecimal(value string) (Decimal, error) {
	intPart, fractPart, _ := strings.Cut(value, ".")

	subunits, err := strconv.ParseInt(intPart+fractPart, 10, 64)
	if err != nil {
		return Decimal{}, fmt.Errorf("%w: %s", ErrInvalidDecimal, err.Error())
	}

	if subunits > maxDecimal {
		return Decimal{}, ErrTooLarge
	}

	precision := byte(len(fractPart))
	return Decimal{subunits: subunits, precision: precision}, nil
}

func (d *Decimal) simplify() {
	// Using %10 returns the last digit in base 10 of a number.
	// If the precision is positive, that digit belongs
	// to the right side of the decimal separator.
	for d.subunits%10 == 0 && d.precision > 0 {
		d.precision--
		d.subunits /= 10
	}
}
