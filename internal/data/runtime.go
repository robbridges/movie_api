package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFOrmat = errors.New("invalid runtime format")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	quotedJsonValue := strconv.Quote(jsonValue)

	return []byte(quotedJsonValue), nil
}

// UnmarshalJSON I'm not thrilled with the idea of the receiever being a pointer in one and not in another, but in this case it needs
// to be a pointer so that it can be modified when we try to unmarshal it
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unquotedJsonValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFOrmat
	}
	parts := strings.Split(unquotedJsonValue, " ")
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFOrmat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFOrmat
	}
	// convert the int32 to a runtime and assign it the reciever we have to deference it to assign the underlying value
	*r = Runtime(i)

	return nil
}
