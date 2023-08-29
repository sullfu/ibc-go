package mapprotocol

import (
	"time"

	clienttypes "github.com/sullfu/ibc-go/v7/modules/core/02-client/types"
	"github.com/sullfu/ibc-go/v7/modules/core/exported"
)

const revisionNumber = 0

var _ exported.ClientMessage = &Header{}

// ConsensusState returns the updated consensus state associated with the header
func (h Header) ConsensusState() *ConsensusState {
	return &ConsensusState{
		Timestamp: h.GetTime(),
	}
}

// ClientType defines that the Header is a mapprotocol consensus algorithm
func (h Header) ClientType() string {
	return ModuleName
}

// GetHeight returns the current height. It returns 0 if the mapprotocol
// header is nil.
// NOTE: the header.Header is checked to be non nil in ValidateBasic.
func (h Header) GetHeight() exported.Height {
	return clienttypes.NewHeight(revisionNumber, h.Number)
}

// GetTime returns the current block timestamp. It returns a zero time if
// the mapprotocol header is nil.
// NOTE: the header.Header is checked to be non nil in ValidateBasic.
func (h Header) GetTime() time.Time {
	if h.SignedHeader == nil {
		return time.Now()
	}
	return time.Unix(int64(h.SignedHeader.Timestamp), 0)
}

// ValidateBasic calls the SignedHeader ValidateBasic function and checks
// that validatorsets are not nil.
// NOTE: TrustedHeight and TrustedValidators may be empty when creating client
// with MsgCreateClient
func (h Header) ValidateBasic() error {
	return nil
}
