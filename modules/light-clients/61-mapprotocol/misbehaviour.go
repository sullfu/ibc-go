package mapprotocol

import (
	"time"

	"github.com/sullfu/ibc-go/v7/modules/core/exported"
)

var _ exported.ClientMessage = &Misbehaviour{}

// ClientType is mapprotocol light client
func (misbehaviour Misbehaviour) ClientType() string {
	return ModuleName
}

// GetTime returns the timestamp at which misbehaviour occurred. It uses the
// maximum value from both headers to prevent producing an invalid header outside
// of the misbehaviour age range.
func (misbehaviour Misbehaviour) GetTime() time.Time {
	t1, t2 := misbehaviour.Header1.GetTime(), misbehaviour.Header2.GetTime()
	if t1.After(t2) {
		return t1
	}
	return t2
}

// ValidateBasic implements Misbehaviour interface
func (misbehaviour Misbehaviour) ValidateBasic() error {
	return nil
}
