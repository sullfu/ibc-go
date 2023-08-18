package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	ibcmock "github.com/cosmos/ibc-go/v7/testing/mock"
)

var (
	timeoutHeight = clienttypes.NewHeight(1, 10000)
	maxSequence   = uint64(10)
)

// tests the IBC handler receiving a packet on ordered and unordered channels.
// It verifies that the storing of an acknowledgement on success occurs. It
// tests high level properties like ordering and basic sanity checks. More
// rigorous testing of 'RecvPacket' can be found in the
// 04-channel/keeper/packet_test.go.
func (suite *KeeperTestSuite) TestHandleRecvPacket() {
	var (
		packet channeltypes.Packet
		path   *ibctesting.Path
		async  bool // indicate no ack written
	)

	testCases := []struct {
		name      string
		malleate  func()
		expPass   bool
		expRevert bool
	}{
		{"success: ORDERED", func() {
			path.SetChannelOrdered()
			suite.coordinator.Setup(path)

			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
			suite.Require().NoError(err)

			packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
		}, true, false},
		{"success: UNORDERED", func() {
			suite.coordinator.Setup(path)

			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
			suite.Require().NoError(err)

			packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
		}, true, false},
		{"success: UNORDERED out of order packet", func() {
			// setup uses an UNORDERED channel
			suite.coordinator.Setup(path)

			// attempts to receive packet with sequence 10 without receiving packet with sequence 1
			for i := uint64(1); i < 10; i++ {
				sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
				suite.Require().NoError(err)

				packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			}
		}, true, false},
		{"success: OnRecvPacket callback returns revert=true", func() {
			suite.coordinator.Setup(path)

			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockFailPacketData)
			suite.Require().NoError(err)

			packet = channeltypes.NewPacket(ibctesting.MockFailPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
		}, true, true},
		{"success: ORDERED - async acknowledgement", func() {
			path.SetChannelOrdered()
			suite.coordinator.Setup(path)
			async = true

			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibcmock.MockAsyncPacketData)
			suite.Require().NoError(err)

			packet = channeltypes.NewPacket(ibcmock.MockAsyncPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
		}, true, false},
		{"success: UNORDERED - async acknowledgement", func() {
			suite.coordinator.Setup(path)
			async = true

			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibcmock.MockAsyncPacketData)
			suite.Require().NoError(err)

			packet = channeltypes.NewPacket(ibcmock.MockAsyncPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
		}, true, false},
		{"failure: ORDERED out of order packet", func() {
			path.SetChannelOrdered()
			suite.coordinator.Setup(path)

			// attempts to receive packet with sequence 10 without receiving packet with sequence 1
			for i := uint64(1); i < 10; i++ {
				sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
				suite.Require().NoError(err)

				packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			}
		}, false, false},
		{"channel does not exist", func() {
			// any non-nil value of packet is valid
			suite.Require().NotNil(packet)
		}, false, false},
		{"packet not sent", func() {
			suite.coordinator.Setup(path)
			packet = channeltypes.NewPacket(ibctesting.MockPacketData, 1, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
		}, false, false},
		{"successful no-op: ORDERED - packet already received (replay)", func() {
			// mock will panic if application callback is called twice on the same packet
			path.SetChannelOrdered()
			suite.coordinator.Setup(path)

			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
			suite.Require().NoError(err)

			packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			err = path.EndpointB.RecvPacket(packet)
			suite.Require().NoError(err)
		}, true, false},
		{"successful no-op: UNORDERED - packet already received (replay)", func() {
			// mock will panic if application callback is called twice on the same packet
			suite.coordinator.Setup(path)

			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
			suite.Require().NoError(err)

			packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			err = path.EndpointB.RecvPacket(packet)
			suite.Require().NoError(err)
		}, true, false},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			async = false     // reset
			path = ibctesting.NewPath(suite.chainA, suite.chainB)

			tc.malleate()

			var (
				proof       []byte
				proofHeight clienttypes.Height
			)

			// get proof of packet commitment from chainA
			packetKey := host.PacketCommitmentKey(packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
			if path.EndpointA.ChannelID != "" {
				proof, proofHeight = path.EndpointA.QueryProof(packetKey)
			}

			msg := channeltypes.NewMsgRecvPacket(packet, proof, proofHeight, suite.chainB.SenderAccount.GetAddress().String())

			_, err := keeper.Keeper.RecvPacket(*suite.chainB.App.GetIBCKeeper(), suite.chainB.GetContext(), msg)

			if tc.expPass {
				suite.Require().NoError(err)

				// replay should not fail since it will be treated as a no-op
				_, err := keeper.Keeper.RecvPacket(*suite.chainB.App.GetIBCKeeper(), suite.chainB.GetContext(), msg)
				suite.Require().NoError(err)

				// check that callback state was handled correctly
				_, exists := suite.chainB.GetSimApp().ScopedIBCMockKeeper.GetCapability(suite.chainB.GetContext(), ibcmock.GetMockRecvCanaryCapabilityName(packet))
				if tc.expRevert {
					suite.Require().False(exists, "capability exists in store even after callback reverted")
				} else {
					suite.Require().True(exists, "callback state not persisted when revert is false")
				}

				// verify if ack was written
				ack, found := suite.chainB.App.GetIBCKeeper().ChannelKeeper.GetPacketAcknowledgement(suite.chainB.GetContext(), packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence())

				if async {
					suite.Require().Nil(ack)
					suite.Require().False(found)

				} else {
					suite.Require().NotNil(ack)
					suite.Require().True(found)
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

// tests the IBC handler acknowledgement of a packet on ordered and unordered
// channels. It verifies that the deletion of packet commitments from state
// occurs. It test high level properties like ordering and basic sanity
// checks. More rigorous testing of 'AcknowledgePacket'
// can be found in the 04-channel/keeper/packet_test.go.
func (suite *KeeperTestSuite) TestHandleAcknowledgePacket() {
	var (
		packet channeltypes.Packet
		path   *ibctesting.Path
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{"success: ORDERED", func() {
			path.SetChannelOrdered()
			suite.coordinator.Setup(path)

			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
			suite.Require().NoError(err)

			packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			err = path.EndpointB.RecvPacket(packet)
			suite.Require().NoError(err)
		}, true},
		{"success: UNORDERED", func() {
			suite.coordinator.Setup(path)

			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
			suite.Require().NoError(err)

			packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			err = path.EndpointB.RecvPacket(packet)
			suite.Require().NoError(err)
		}, true},
		{"success: UNORDERED acknowledge out of order packet", func() {
			// setup uses an UNORDERED channel
			suite.coordinator.Setup(path)

			// attempts to acknowledge ack with sequence 10 without acknowledging ack with sequence 1 (removing packet commitment)
			for i := uint64(1); i < 10; i++ {
				sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
				suite.Require().NoError(err)

				packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
				err = path.EndpointB.RecvPacket(packet)
				suite.Require().NoError(err)
			}
		}, true},
		{"failure: ORDERED acknowledge out of order packet", func() {
			path.SetChannelOrdered()
			suite.coordinator.Setup(path)

			// attempts to acknowledge ack with sequence 10 without acknowledging ack with sequence 1 (removing packet commitment
			for i := uint64(1); i < 10; i++ {
				sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
				suite.Require().NoError(err)

				packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
				err = path.EndpointB.RecvPacket(packet)
				suite.Require().NoError(err)
			}
		}, false},
		{"channel does not exist", func() {
			// any non-nil value of packet is valid
			suite.Require().NotNil(packet)
		}, false},
		{"packet not received", func() {
			suite.coordinator.Setup(path)

			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
			suite.Require().NoError(err)

			packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
		}, false},
		{"successful no-op: ORDERED - packet already acknowledged (replay)", func() {
			suite.coordinator.Setup(path)

			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
			suite.Require().NoError(err)

			packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			err = path.EndpointB.RecvPacket(packet)
			suite.Require().NoError(err)

			err = path.EndpointA.AcknowledgePacket(packet, ibctesting.MockAcknowledgement)
			suite.Require().NoError(err)
		}, true},
		{"successful no-op: UNORDERED - packet already acknowledged (replay)", func() {
			suite.coordinator.Setup(path)

			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
			suite.Require().NoError(err)

			packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			err = path.EndpointB.RecvPacket(packet)
			suite.Require().NoError(err)

			err = path.EndpointA.AcknowledgePacket(packet, ibctesting.MockAcknowledgement)
			suite.Require().NoError(err)
		}, true},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			path = ibctesting.NewPath(suite.chainA, suite.chainB)

			tc.malleate()

			var (
				proof       []byte
				proofHeight clienttypes.Height
			)
			packetKey := host.PacketAcknowledgementKey(packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence())
			if path.EndpointB.ChannelID != "" {
				proof, proofHeight = path.EndpointB.QueryProof(packetKey)
			}

			msg := channeltypes.NewMsgAcknowledgement(packet, ibcmock.MockAcknowledgement.Acknowledgement(), proof, proofHeight, suite.chainA.SenderAccount.GetAddress().String())

			_, err := keeper.Keeper.Acknowledgement(*suite.chainA.App.GetIBCKeeper(), suite.chainA.GetContext(), msg)

			if tc.expPass {
				suite.Require().NoError(err)

				// verify packet commitment was deleted on source chain
				has := suite.chainA.App.GetIBCKeeper().ChannelKeeper.HasPacketCommitment(suite.chainA.GetContext(), packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
				suite.Require().False(has)

				// replay should not error as it is treated as a no-op
				_, err := keeper.Keeper.Acknowledgement(*suite.chainA.App.GetIBCKeeper(), suite.chainA.GetContext(), msg)
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

// tests the IBC handler timing out a packet on ordered and unordered channels.
// It verifies that the deletion of a packet commitment occurs. It tests
// high level properties like ordering and basic sanity checks. More
// rigorous testing of 'TimeoutPacket' and 'TimeoutExecuted' can be found in
// the 04-channel/keeper/timeout_test.go.
func (suite *KeeperTestSuite) TestHandleTimeoutPacket() {
	var (
		packet    channeltypes.Packet
		packetKey []byte
		path      *ibctesting.Path
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{"success: ORDERED", func() {
			path.SetChannelOrdered()
			suite.coordinator.Setup(path)

			timeoutHeight := clienttypes.GetSelfHeight(suite.chainB.GetContext())
			timeoutTimestamp := uint64(suite.chainB.GetContext().BlockTime().UnixNano())

			// create packet commitment
			sequence, err := path.EndpointA.SendPacket(timeoutHeight, timeoutTimestamp, ibctesting.MockPacketData)
			suite.Require().NoError(err)

			// need to update chainA client to prove missing ack
			err = path.EndpointA.UpdateClient()
			suite.Require().NoError(err)

			packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, timeoutTimestamp)
			packetKey = host.NextSequenceRecvKey(packet.GetDestPort(), packet.GetDestChannel())
		}, true},
		{"success: UNORDERED", func() {
			suite.coordinator.Setup(path)

			timeoutHeight := clienttypes.GetSelfHeight(suite.chainB.GetContext())
			timeoutTimestamp := uint64(suite.chainB.GetContext().BlockTime().UnixNano())

			// create packet commitment
			sequence, err := path.EndpointA.SendPacket(timeoutHeight, timeoutTimestamp, ibctesting.MockPacketData)
			suite.Require().NoError(err)

			// need to update chainA client to prove missing ack
			err = path.EndpointA.UpdateClient()
			suite.Require().NoError(err)

			packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, timeoutTimestamp)
			packetKey = host.PacketReceiptKey(packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence())
		}, true},
		{"success: UNORDERED timeout out of order packet", func() {
			// setup uses an UNORDERED channel
			suite.coordinator.Setup(path)

			// attempts to timeout the last packet sent without timing out the first packet
			// packet sequences begin at 1
			for i := uint64(1); i < maxSequence; i++ {
				timeoutHeight := clienttypes.GetSelfHeight(suite.chainB.GetContext())

				// create packet commitment
				sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
				suite.Require().NoError(err)

				packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			}

			err := path.EndpointA.UpdateClient()
			suite.Require().NoError(err)

			packetKey = host.PacketReceiptKey(packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence())
		}, true},
		{"success: ORDERED timeout out of order packet", func() {
			path.SetChannelOrdered()
			suite.coordinator.Setup(path)

			// attempts to timeout the last packet sent without timing out the first packet
			// packet sequences begin at 1
			for i := uint64(1); i < maxSequence; i++ {
				timeoutHeight := clienttypes.GetSelfHeight(suite.chainB.GetContext())

				// create packet commitment
				sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
				suite.Require().NoError(err)

				packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			}

			err := path.EndpointA.UpdateClient()
			suite.Require().NoError(err)

			packetKey = host.NextSequenceRecvKey(packet.GetDestPort(), packet.GetDestChannel())
		}, true},
		{"channel does not exist", func() {
			// any non-nil value of packet is valid
			suite.Require().NotNil(packet)

			packetKey = host.NextSequenceRecvKey(packet.GetDestPort(), packet.GetDestChannel())
		}, false},
		{"successful no-op: UNORDERED - packet not sent", func() {
			suite.coordinator.Setup(path)
			packet = channeltypes.NewPacket(ibctesting.MockPacketData, 1, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, clienttypes.NewHeight(0, 1), 0)
			packetKey = host.PacketReceiptKey(packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence())
		}, true},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			path = ibctesting.NewPath(suite.chainA, suite.chainB)

			tc.malleate()

			var (
				proof       []byte
				proofHeight clienttypes.Height
			)
			if path.EndpointB.ChannelID != "" {
				proof, proofHeight = path.EndpointB.QueryProof(packetKey)
			}

			msg := channeltypes.NewMsgTimeout(packet, 1, proof, proofHeight, suite.chainA.SenderAccount.GetAddress().String())

			_, err := keeper.Keeper.Timeout(*suite.chainA.App.GetIBCKeeper(), suite.chainA.GetContext(), msg)

			if tc.expPass {
				suite.Require().NoError(err)

				// replay should not return an error as it is treated as a no-op
				_, err := keeper.Keeper.Timeout(*suite.chainA.App.GetIBCKeeper(), suite.chainA.GetContext(), msg)
				suite.Require().NoError(err)

				// verify packet commitment was deleted on source chain
				has := suite.chainA.App.GetIBCKeeper().ChannelKeeper.HasPacketCommitment(suite.chainA.GetContext(), packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
				suite.Require().False(has)

			} else {
				suite.Require().Error(err)
			}
		})
	}
}

// tests the IBC handler timing out a packet via channel closure on ordered
// and unordered channels. It verifies that the deletion of a packet
// commitment occurs. It tests high level properties like ordering and basic
// sanity checks. More rigorous testing of 'TimeoutOnClose' and
// 'TimeoutExecuted' can be found in the 04-channel/keeper/timeout_test.go.
func (suite *KeeperTestSuite) TestHandleTimeoutOnClosePacket() {
	var (
		packet                      channeltypes.Packet
		packetKey                   []byte
		path                        *ibctesting.Path
		counterpartyUpgradeSequence uint64
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{"success: ORDERED", func() {
			path.SetChannelOrdered()
			suite.coordinator.Setup(path)

			// create packet commitment
			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
			suite.Require().NoError(err)

			// need to update chainA client to prove missing ack
			err = path.EndpointA.UpdateClient()
			suite.Require().NoError(err)

			packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			packetKey = host.NextSequenceRecvKey(packet.GetDestPort(), packet.GetDestChannel())

			// close counterparty channel
			err = path.EndpointB.SetChannelState(channeltypes.CLOSED)
			suite.Require().NoError(err)
		}, true},
		{"success: UNORDERED", func() {
			suite.coordinator.Setup(path)

			// create packet commitment
			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
			suite.Require().NoError(err)

			// need to update chainA client to prove missing ack
			err = path.EndpointA.UpdateClient()
			suite.Require().NoError(err)

			packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			packetKey = host.PacketReceiptKey(packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence())

			// close counterparty channel
			err = path.EndpointB.SetChannelState(channeltypes.CLOSED)
			suite.Require().NoError(err)
		}, true},
		{"success: UNORDERED timeout out of order packet", func() {
			// setup uses an UNORDERED channel
			suite.coordinator.Setup(path)

			// attempts to timeout the last packet sent without timing out the first packet
			// packet sequences begin at 1
			for i := uint64(1); i < maxSequence; i++ {
				// create packet commitment
				sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
				suite.Require().NoError(err)

				packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			}

			err := path.EndpointA.UpdateClient()
			suite.Require().NoError(err)

			packetKey = host.PacketReceiptKey(packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence())

			// close counterparty channel
			err = path.EndpointB.SetChannelState(channeltypes.CLOSED)
			suite.Require().NoError(err)
		}, true},
		{"success: ORDERED timeout out of order packet", func() {
			path.SetChannelOrdered()
			suite.coordinator.Setup(path)

			// attempts to timeout the last packet sent without timing out the first packet
			// packet sequences begin at 1
			for i := uint64(1); i < maxSequence; i++ {
				// create packet commitment
				sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
				suite.Require().NoError(err)

				packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			}

			err := path.EndpointA.UpdateClient()
			suite.Require().NoError(err)

			packetKey = host.NextSequenceRecvKey(packet.GetDestPort(), packet.GetDestChannel())

			// close counterparty channel
			err = path.EndpointB.SetChannelState(channeltypes.CLOSED)
			suite.Require().NoError(err)
		}, true},
		{"channel does not exist", func() {
			// any non-nil value of packet is valid
			suite.Require().NotNil(packet)

			packetKey = host.NextSequenceRecvKey(packet.GetDestPort(), packet.GetDestChannel())
		}, false},
		{"successful no-op: UNORDERED - packet not sent", func() {
			suite.coordinator.Setup(path)
			packet = channeltypes.NewPacket(ibctesting.MockPacketData, 1, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, clienttypes.NewHeight(0, 1), 0)
			packetKey = host.PacketAcknowledgementKey(packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence())

			// close counterparty channel
			err := path.EndpointB.SetChannelState(channeltypes.CLOSED)
			suite.Require().NoError(err)
		}, true},
		{"ORDERED: channel not closed", func() {
			path.SetChannelOrdered()
			suite.coordinator.Setup(path)

			// create packet commitment
			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
			suite.Require().NoError(err)

			// need to update chainA client to prove missing ack
			err = path.EndpointA.UpdateClient()
			suite.Require().NoError(err)

			packet = channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			packetKey = host.NextSequenceRecvKey(packet.GetDestPort(), packet.GetDestChannel())
		}, false},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			path = ibctesting.NewPath(suite.chainA, suite.chainB)

			tc.malleate()

			proof, proofHeight := suite.chainB.QueryProof(packetKey)

			channelKey := host.ChannelKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			proofClosed, _ := suite.chainB.QueryProof(channelKey)

			msg := channeltypes.NewMsgTimeoutOnClose(packet, 1, proof, proofClosed, proofHeight, suite.chainA.SenderAccount.GetAddress().String(), counterpartyUpgradeSequence)

			_, err := keeper.Keeper.TimeoutOnClose(*suite.chainA.App.GetIBCKeeper(), suite.chainA.GetContext(), msg)

			if tc.expPass {
				suite.Require().NoError(err)

				// replay should not return an error as it will be treated as a no-op
				_, err := keeper.Keeper.TimeoutOnClose(*suite.chainA.App.GetIBCKeeper(), suite.chainA.GetContext(), msg)
				suite.Require().NoError(err)

				// verify packet commitment was deleted on source chain
				has := suite.chainA.App.GetIBCKeeper().ChannelKeeper.HasPacketCommitment(suite.chainA.GetContext(), packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
				suite.Require().False(has)

			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestUpgradeClient() {
	var (
		path              *ibctesting.Path
		newChainID        string
		newClientHeight   clienttypes.Height
		upgradedClient    exported.ClientState
		upgradedConsState exported.ConsensusState
		lastHeight        exported.Height
		msg               *clienttypes.MsgUpgradeClient
	)
	cases := []struct {
		name    string
		setup   func()
		expPass bool
	}{
		{
			name: "successful upgrade",
			setup: func() {
				upgradedClient = ibctm.NewClientState(newChainID, ibctm.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod+ibctesting.TrustingPeriod, ibctesting.MaxClockDrift, newClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath)
				// Call ZeroCustomFields on upgraded clients to clear any client-chosen parameters in test-case upgradedClient
				upgradedClient = upgradedClient.ZeroCustomFields()

				upgradedConsState = &ibctm.ConsensusState{
					NextValidatorsHash: []byte("nextValsHash"),
				}

				// last Height is at next block
				lastHeight = clienttypes.NewHeight(0, uint64(suite.chainB.GetContext().BlockHeight()+1))

				upgradedClientBz, err := clienttypes.MarshalClientState(suite.chainA.App.AppCodec(), upgradedClient)
				suite.Require().NoError(err)
				upgradedConsStateBz, err := clienttypes.MarshalConsensusState(suite.chainA.App.AppCodec(), upgradedConsState)
				suite.Require().NoError(err)

				// zero custom fields and store in upgrade store
				suite.chainB.GetSimApp().UpgradeKeeper.SetUpgradedClient(suite.chainB.GetContext(), int64(lastHeight.GetRevisionHeight()), upgradedClientBz)            //nolint:errcheck // ignore error for testing
				suite.chainB.GetSimApp().UpgradeKeeper.SetUpgradedConsensusState(suite.chainB.GetContext(), int64(lastHeight.GetRevisionHeight()), upgradedConsStateBz) //nolint:errcheck // ignore error for testing

				// commit upgrade store changes and update clients
				suite.coordinator.CommitBlock(suite.chainB)
				err = path.EndpointA.UpdateClient()
				suite.Require().NoError(err)

				cs, found := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(suite.chainA.GetContext(), path.EndpointA.ClientID)
				suite.Require().True(found)

				proofUpgradeClient, _ := suite.chainB.QueryUpgradeProof(upgradetypes.UpgradedClientKey(int64(lastHeight.GetRevisionHeight())), cs.GetLatestHeight().GetRevisionHeight())
				proofUpgradedConsState, _ := suite.chainB.QueryUpgradeProof(upgradetypes.UpgradedConsStateKey(int64(lastHeight.GetRevisionHeight())), cs.GetLatestHeight().GetRevisionHeight())

				msg, err = clienttypes.NewMsgUpgradeClient(path.EndpointA.ClientID, upgradedClient, upgradedConsState,
					proofUpgradeClient, proofUpgradedConsState, suite.chainA.SenderAccount.GetAddress().String())
				suite.Require().NoError(err)
			},
			expPass: true,
		},
		{
			name: "VerifyUpgrade fails",
			setup: func() {
				upgradedClient = ibctm.NewClientState(newChainID, ibctm.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod+ibctesting.TrustingPeriod, ibctesting.MaxClockDrift, newClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath)
				// Call ZeroCustomFields on upgraded clients to clear any client-chosen parameters in test-case upgradedClient
				upgradedClient = upgradedClient.ZeroCustomFields()

				upgradedConsState = &ibctm.ConsensusState{
					NextValidatorsHash: []byte("nextValsHash"),
				}

				// last Height is at next block
				lastHeight = clienttypes.NewHeight(0, uint64(suite.chainB.GetContext().BlockHeight()+1))

				upgradedClientBz, err := clienttypes.MarshalClientState(suite.chainA.App.AppCodec(), upgradedClient)
				suite.Require().NoError(err)
				upgradedConsStateBz, err := clienttypes.MarshalConsensusState(suite.chainA.App.AppCodec(), upgradedConsState)
				suite.Require().NoError(err)

				// zero custom fields and store in upgrade store
				suite.chainB.GetSimApp().UpgradeKeeper.SetUpgradedClient(suite.chainB.GetContext(), int64(lastHeight.GetRevisionHeight()), upgradedClientBz)            //nolint:errcheck // ignore error for testing
				suite.chainB.GetSimApp().UpgradeKeeper.SetUpgradedConsensusState(suite.chainB.GetContext(), int64(lastHeight.GetRevisionHeight()), upgradedConsStateBz) //nolint:errcheck // ignore error for testing

				// commit upgrade store changes and update clients
				suite.coordinator.CommitBlock(suite.chainB)
				err = path.EndpointA.UpdateClient()
				suite.Require().NoError(err)

				msg, err = clienttypes.NewMsgUpgradeClient(path.EndpointA.ClientID, upgradedClient, upgradedConsState, nil, nil, suite.chainA.SenderAccount.GetAddress().String())
				suite.Require().NoError(err)
			},
			expPass: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		path = ibctesting.NewPath(suite.chainA, suite.chainB)
		suite.coordinator.SetupClients(path)

		var err error
		clientState := path.EndpointA.GetClientState().(*ibctm.ClientState)
		revisionNumber := clienttypes.ParseChainID(clientState.ChainId)

		newChainID, err = clienttypes.SetRevisionNumber(clientState.ChainId, revisionNumber+1)
		suite.Require().NoError(err)

		newClientHeight = clienttypes.NewHeight(revisionNumber+1, clientState.GetLatestHeight().GetRevisionHeight()+1)

		tc.setup()

		_, err = keeper.Keeper.UpgradeClient(*suite.chainA.App.GetIBCKeeper(), suite.chainA.GetContext(), msg)

		if tc.expPass {
			suite.Require().NoError(err, "upgrade handler failed on valid case: %s", tc.name)
			newClient, ok := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(suite.chainA.GetContext(), path.EndpointA.ClientID)
			suite.Require().True(ok)
			newChainSpecifiedClient := newClient.ZeroCustomFields()
			suite.Require().Equal(upgradedClient, newChainSpecifiedClient)
		} else {
			suite.Require().Error(err, "upgrade handler passed on invalid case: %s", tc.name)
		}
	}
}

func (suite *KeeperTestSuite) TestChannelUpgradeTry() {
	var (
		path *ibctesting.Path
		msg  *channeltypes.MsgChannelUpgradeTry
	)

	cases := []struct {
		name      string
		malleate  func()
		expResult func(res *channeltypes.MsgChannelUpgradeTryResponse, err error)
	}{
		{
			"success",
			func() {},
			func(res *channeltypes.MsgChannelUpgradeTryResponse, err error) {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(channeltypes.SUCCESS, res.Result)

				channel := path.EndpointB.GetChannel()
				suite.Require().Equal(channeltypes.STATE_FLUSHING, channel.State)
				suite.Require().Equal(uint64(1), channel.UpgradeSequence)
			},
		},
		{
			"module capability not found",
			func() {
				msg.PortId = "invalid-port"
				msg.ChannelId = "invalid-channel"
			},
			func(res *channeltypes.MsgChannelUpgradeTryResponse, err error) {
				suite.Require().Error(err)
				suite.Require().Nil(res)

				suite.Require().ErrorIs(err, capabilitytypes.ErrCapabilityNotFound)
			},
		},
		{
			"unsynchronized upgrade sequence writes upgrade error receipt",
			func() {
				channel := path.EndpointB.GetChannel()
				channel.UpgradeSequence = 99

				path.EndpointB.SetChannel(channel)
			},
			func(res *channeltypes.MsgChannelUpgradeTryResponse, err error) {
				suite.Require().NoError(err)

				suite.Require().NotNil(res)
				suite.Require().Equal(channeltypes.FAILURE, res.Result)

				errorReceipt, found := suite.chainB.GetSimApp().GetIBCKeeper().ChannelKeeper.GetUpgradeErrorReceipt(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				suite.Require().True(found)
				suite.Require().Equal(uint64(99), errorReceipt.Sequence)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			path = ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(path)

			// configure the channel upgrade version on testing endpoints
			path.EndpointA.ChannelConfig.ProposedUpgrade.Fields.Version = ibcmock.UpgradeVersion
			path.EndpointB.ChannelConfig.ProposedUpgrade.Fields.Version = ibcmock.UpgradeVersion

			err := path.EndpointA.ChanUpgradeInit()
			suite.Require().NoError(err)

			err = path.EndpointB.UpdateClient()
			suite.Require().NoError(err)

			counterpartySequence := path.EndpointA.GetChannel().UpgradeSequence
			counterpartyUpgrade, found := suite.chainA.GetSimApp().GetIBCKeeper().ChannelKeeper.GetUpgrade(suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
			suite.Require().True(found)

			proofChannel, proofUpgrade, proofHeight := path.EndpointB.QueryChannelUpgradeProof()

			msg = &channeltypes.MsgChannelUpgradeTry{
				PortId:                        path.EndpointB.ChannelConfig.PortID,
				ChannelId:                     path.EndpointB.ChannelID,
				ProposedUpgradeConnectionHops: []string{ibctesting.FirstConnectionID},
				CounterpartyUpgradeSequence:   counterpartySequence,
				CounterpartyUpgradeFields:     counterpartyUpgrade.Fields,
				ProofChannel:                  proofChannel,
				ProofUpgrade:                  proofUpgrade,
				ProofHeight:                   proofHeight,
				Signer:                        suite.chainB.SenderAccount.GetAddress().String(),
			}

			tc.malleate()

			res, err := suite.chainB.GetSimApp().GetIBCKeeper().ChannelUpgradeTry(suite.chainB.GetContext(), msg)

			tc.expResult(res, err)
		})
	}
}

func (suite *KeeperTestSuite) TestChannelUpgradeAck() {
	var (
		path *ibctesting.Path
		msg  *channeltypes.MsgChannelUpgradeAck
	)

	cases := []struct {
		name      string
		malleate  func()
		expResult func(res *channeltypes.MsgChannelUpgradeAckResponse, err error)
	}{
		{
			"success, no pending in-flight packets",
			func() {},
			func(res *channeltypes.MsgChannelUpgradeAckResponse, err error) {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(channeltypes.SUCCESS, res.Result)

				channel := path.EndpointA.GetChannel()
				suite.Require().Equal(channeltypes.STATE_FLUSHCOMPLETE, channel.State)
				suite.Require().Equal(uint64(1), channel.UpgradeSequence)
			},
		},
		{
			"success, pending in-flight packets",
			func() {
				portID := path.EndpointA.ChannelConfig.PortID
				channelID := path.EndpointA.ChannelID
				// Set a dummy packet commitment to simulate in-flight packets
				suite.chainA.GetSimApp().IBCKeeper.ChannelKeeper.SetPacketCommitment(suite.chainA.GetContext(), portID, channelID, 1, []byte("hash"))
			},
			func(res *channeltypes.MsgChannelUpgradeAckResponse, err error) {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(channeltypes.SUCCESS, res.Result)

				channel := path.EndpointA.GetChannel()
				suite.Require().Equal(channeltypes.STATE_FLUSHING, channel.State)
				suite.Require().Equal(uint64(1), channel.UpgradeSequence)
			},
		},
		{
			"module capability not found",
			func() {
				msg.PortId = ibctesting.InvalidID
				msg.ChannelId = ibctesting.InvalidID
			},
			func(res *channeltypes.MsgChannelUpgradeAckResponse, err error) {
				suite.Require().Error(err)
				suite.Require().Nil(res)

				suite.Require().ErrorIs(err, capabilitytypes.ErrCapabilityNotFound)
			},
		},
		{
			"core handler returns error and no upgrade error receipt is written",
			func() {
				// force an error by overriding the channel state to an invalid value
				channel := path.EndpointA.GetChannel()
				channel.State = channeltypes.CLOSED

				path.EndpointA.SetChannel(channel)
			},
			func(res *channeltypes.MsgChannelUpgradeAckResponse, err error) {
				suite.Require().Error(err)
				suite.Require().Nil(res)
				suite.Require().ErrorIs(err, channeltypes.ErrInvalidChannelState)

				errorReceipt, found := suite.chainA.GetSimApp().GetIBCKeeper().ChannelKeeper.GetUpgradeErrorReceipt(suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				suite.Require().Empty(errorReceipt)
				suite.Require().False(found)
			},
		},
		{
			"core handler returns error and writes upgrade error receipt",
			func() {
				// force an upgrade error by modifying the channel upgrade ordering to an incompatible value
				upgrade := path.EndpointA.GetChannelUpgrade()
				upgrade.Fields.Ordering = channeltypes.NONE

				path.EndpointA.SetChannelUpgrade(upgrade)
			},
			func(res *channeltypes.MsgChannelUpgradeAckResponse, err error) {
				suite.Require().NoError(err)

				suite.Require().NotNil(res)
				suite.Require().Equal(channeltypes.FAILURE, res.Result)

				errorReceipt, found := suite.chainA.GetSimApp().GetIBCKeeper().ChannelKeeper.GetUpgradeErrorReceipt(suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				suite.Require().True(found)
				suite.Require().Equal(uint64(1), errorReceipt.Sequence)
			},
		},
		{
			"application callback returns error and error receipt is written",
			func() {
				suite.chainA.GetSimApp().IBCMockModule.IBCApp.OnChanUpgradeAck = func(
					ctx sdk.Context, portID, channelID, counterpartyVersion string,
				) error {
					// set arbitrary value in store to mock application state changes
					store := ctx.KVStore(suite.chainA.GetSimApp().GetKey(exported.ModuleName))
					store.Set([]byte("foo"), []byte("bar"))
					return fmt.Errorf("mock app callback failed")
				}
			},
			func(res *channeltypes.MsgChannelUpgradeAckResponse, err error) {
				suite.Require().NoError(err)

				suite.Require().NotNil(res)
				suite.Require().Equal(channeltypes.FAILURE, res.Result)

				errorReceipt, found := suite.chainA.GetSimApp().GetIBCKeeper().ChannelKeeper.GetUpgradeErrorReceipt(suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				suite.Require().True(found)
				suite.Require().Equal(uint64(1), errorReceipt.Sequence)

				// assert application state changes are not committed
				store := suite.chainA.GetContext().KVStore(suite.chainA.GetSimApp().GetKey(exported.ModuleName))
				suite.Require().False(store.Has([]byte("foo")))
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			path = ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(path)

			// configure the channel upgrade version on testing endpoints
			path.EndpointA.ChannelConfig.ProposedUpgrade.Fields.Version = ibcmock.UpgradeVersion
			path.EndpointB.ChannelConfig.ProposedUpgrade.Fields.Version = ibcmock.UpgradeVersion

			err := path.EndpointA.ChanUpgradeInit()
			suite.Require().NoError(err)

			err = path.EndpointB.ChanUpgradeTry()
			suite.Require().NoError(err)

			err = path.EndpointA.UpdateClient()
			suite.Require().NoError(err)

			counterpartyUpgrade := path.EndpointB.GetChannelUpgrade()

			proofChannel, proofUpgrade, proofHeight := path.EndpointA.QueryChannelUpgradeProof()

			msg = &channeltypes.MsgChannelUpgradeAck{
				PortId:              path.EndpointA.ChannelConfig.PortID,
				ChannelId:           path.EndpointA.ChannelID,
				CounterpartyUpgrade: counterpartyUpgrade,
				ProofChannel:        proofChannel,
				ProofUpgrade:        proofUpgrade,
				ProofHeight:         proofHeight,
				Signer:              suite.chainA.SenderAccount.GetAddress().String(),
			}

			tc.malleate()

			res, err := suite.chainA.GetSimApp().GetIBCKeeper().ChannelUpgradeAck(suite.chainA.GetContext(), msg)

			tc.expResult(res, err)
		})
	}
}

func (suite *KeeperTestSuite) TestChannelUpgradeConfirm() {
	var (
		path *ibctesting.Path
		msg  *channeltypes.MsgChannelUpgradeConfirm
	)

	cases := []struct {
		name      string
		malleate  func()
		expResult func(res *channeltypes.MsgChannelUpgradeConfirmResponse, err error)
	}{
		{
			"success, no pending in-flight packets",
			func() {},
			func(res *channeltypes.MsgChannelUpgradeConfirmResponse, err error) {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(channeltypes.SUCCESS, res.Result)

				channel := path.EndpointB.GetChannel()
				suite.Require().Equal(channeltypes.OPEN, channel.State)
				suite.Require().Equal(uint64(1), channel.UpgradeSequence)
			},
		},
		{
			"success, pending in-flight packets",
			func() {
				portID := path.EndpointB.ChannelConfig.PortID
				channelID := path.EndpointB.ChannelID
				// Set a dummy packet commitment to simulate in-flight packets
				suite.chainB.GetSimApp().IBCKeeper.ChannelKeeper.SetPacketCommitment(suite.chainB.GetContext(), portID, channelID, 1, []byte("hash"))
			},
			func(res *channeltypes.MsgChannelUpgradeConfirmResponse, err error) {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(channeltypes.SUCCESS, res.Result)

				channel := path.EndpointB.GetChannel()
				suite.Require().Equal(channeltypes.STATE_FLUSHING, channel.State)
				suite.Require().Equal(uint64(1), channel.UpgradeSequence)
			},
		},
		{
			"module capability not found",
			func() {
				msg.PortId = ibctesting.InvalidID
				msg.ChannelId = ibctesting.InvalidID
			},
			func(res *channeltypes.MsgChannelUpgradeConfirmResponse, err error) {
				suite.Require().Error(err)
				suite.Require().Nil(res)

				suite.Require().ErrorIs(err, capabilitytypes.ErrCapabilityNotFound)
			},
		},
		{
			"core handler returns error and no upgrade error receipt is written",
			func() {
				// force an error by overriding the channel state to an invalid value
				channel := path.EndpointB.GetChannel()
				channel.State = channeltypes.CLOSED

				path.EndpointB.SetChannel(channel)
			},
			func(res *channeltypes.MsgChannelUpgradeConfirmResponse, err error) {
				suite.Require().Error(err)
				suite.Require().Nil(res)
				suite.Require().ErrorIs(err, channeltypes.ErrInvalidChannelState)

				errorReceipt, found := suite.chainB.GetSimApp().GetIBCKeeper().ChannelKeeper.GetUpgradeErrorReceipt(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				suite.Require().Empty(errorReceipt)
				suite.Require().False(found)
			},
		},
		{
			"core handler returns error and writes upgrade error receipt",
			func() {
				// force an upgrade error by modifying the counterparty channel upgrade timeout to be no longer valid
				upgrade := path.EndpointA.GetChannelUpgrade()
				upgrade.Timeout = channeltypes.NewTimeout(clienttypes.ZeroHeight(), 0)

				path.EndpointA.SetChannelUpgrade(upgrade)

				suite.coordinator.CommitBlock(suite.chainA)

				err := path.EndpointB.UpdateClient()
				suite.Require().NoError(err)

				proofChannel, proofUpgrade, proofHeight := path.EndpointB.QueryChannelUpgradeProof()

				msg.CounterpartyUpgrade = upgrade
				msg.ProofChannel = proofChannel
				msg.ProofUpgrade = proofUpgrade
				msg.ProofHeight = proofHeight
			},
			func(res *channeltypes.MsgChannelUpgradeConfirmResponse, err error) {
				suite.Require().NoError(err)

				suite.Require().NotNil(res)
				suite.Require().Equal(channeltypes.FAILURE, res.Result)

				errorReceipt, found := suite.chainB.GetSimApp().GetIBCKeeper().ChannelKeeper.GetUpgradeErrorReceipt(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				suite.Require().True(found)
				suite.Require().Equal(uint64(1), errorReceipt.Sequence)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			path = ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(path)

			// configure the channel upgrade version on testing endpoints
			path.EndpointA.ChannelConfig.ProposedUpgrade.Fields.Version = ibcmock.UpgradeVersion
			path.EndpointB.ChannelConfig.ProposedUpgrade.Fields.Version = ibcmock.UpgradeVersion

			err := path.EndpointA.ChanUpgradeInit()
			suite.Require().NoError(err)

			err = path.EndpointB.ChanUpgradeTry()
			suite.Require().NoError(err)

			err = path.EndpointA.ChanUpgradeAck()
			suite.Require().NoError(err)

			err = path.EndpointB.UpdateClient()
			suite.Require().NoError(err)

			counterpartyChannelState := path.EndpointA.GetChannel().State
			counterpartyUpgrade := path.EndpointA.GetChannelUpgrade()

			proofChannel, proofUpgrade, proofHeight := path.EndpointB.QueryChannelUpgradeProof()

			msg = &channeltypes.MsgChannelUpgradeConfirm{
				PortId:                   path.EndpointB.ChannelConfig.PortID,
				ChannelId:                path.EndpointB.ChannelID,
				CounterpartyChannelState: counterpartyChannelState,
				CounterpartyUpgrade:      counterpartyUpgrade,
				ProofChannel:             proofChannel,
				ProofUpgrade:             proofUpgrade,
				ProofHeight:              proofHeight,
				Signer:                   suite.chainA.SenderAccount.GetAddress().String(),
			}

			tc.malleate()

			res, err := suite.chainB.GetSimApp().GetIBCKeeper().ChannelUpgradeConfirm(suite.chainB.GetContext(), msg)

			tc.expResult(res, err)
		})
	}
}

func (suite *KeeperTestSuite) TestChannelUpgradeOpen() {
	var (
		path *ibctesting.Path
		msg  *channeltypes.MsgChannelUpgradeOpen
	)

	cases := []struct {
		name      string
		malleate  func()
		expResult func(res *channeltypes.MsgChannelUpgradeOpenResponse, err error)
	}{
		// TODO: Uncomment and address testcases when appropriate, timeout logic currently causes failures
		// {
		// 	"success",
		// 	func() {},
		// 	func(res *channeltypes.MsgChannelUpgradeOpenResponse, err error) {
		// 		suite.Require().NoError(err)
		// 		suite.Require().NotNil(res)

		// 		channel := path.EndpointA.GetChannel()
		// 		suite.Require().Equal(channeltypes.OPEN, channel.State)
		// 		suite.Require().Equal(channeltypes.NOTINFLUSH, channel.FlushStatus)
		// 	},
		// },
		// {
		// 	"module capability not found",
		// 	func() {
		// 		msg.PortId = ibctesting.InvalidID
		// 		msg.ChannelId = ibctesting.InvalidID
		// 	},
		// 	func(res *channeltypes.MsgChannelUpgradeOpenResponse, err error) {
		// 		suite.Require().Error(err)
		// 		suite.Require().Nil(res)

		// 		suite.Require().ErrorIs(err, capabilitytypes.ErrCapabilityNotFound)
		// 	},
		// },
		// {
		// 	"core handler fails",
		// 	func() {
		// 		portID := path.EndpointA.ChannelConfig.PortID
		// 		channelID := path.EndpointA.ChannelID
		// 		// Set a dummy packet commitment to simulate in-flight packets
		// 		suite.chainA.GetSimApp().IBCKeeper.ChannelKeeper.SetPacketCommitment(suite.chainA.GetContext(), portID, channelID, 1, []byte("hash"))
		// 	},
		// 	func(res *channeltypes.MsgChannelUpgradeOpenResponse, err error) {
		// 		suite.Require().Error(err)
		// 		suite.Require().Nil(res)

		// 		suite.Require().ErrorIs(err, channeltypes.ErrPendingInflightPackets)
		// 	},
		// },
	}

	for _, tc := range cases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			path = ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(path)

			// configure the channel upgrade version on testing endpoints
			path.EndpointA.ChannelConfig.ProposedUpgrade.Fields.Version = ibcmock.UpgradeVersion
			path.EndpointB.ChannelConfig.ProposedUpgrade.Fields.Version = ibcmock.UpgradeVersion

			err := path.EndpointB.ChanUpgradeInit()
			suite.Require().NoError(err)

			err = path.EndpointA.ChanUpgradeTry()
			suite.Require().NoError(err)

			err = path.EndpointB.ChanUpgradeAck()
			suite.Require().NoError(err)

			err = path.EndpointA.UpdateClient()
			suite.Require().NoError(err)

			counterpartyChannel := path.EndpointB.GetChannel()
			proofChannel, _, proofHeight := path.EndpointA.QueryChannelUpgradeProof()

			msg = &channeltypes.MsgChannelUpgradeOpen{
				PortId:                   path.EndpointA.ChannelConfig.PortID,
				ChannelId:                path.EndpointA.ChannelID,
				CounterpartyChannelState: counterpartyChannel.State,
				ProofChannel:             proofChannel,
				ProofHeight:              proofHeight,
				Signer:                   suite.chainA.SenderAccount.GetAddress().String(),
			}

			tc.malleate()
			res, err := suite.chainA.GetSimApp().GetIBCKeeper().ChannelUpgradeOpen(suite.chainA.GetContext(), msg)

			tc.expResult(res, err)
		})
	}
}

func (suite *KeeperTestSuite) TestChannelUpgradeCancel() {
	var (
		path *ibctesting.Path
		msg  *channeltypes.MsgChannelUpgradeCancel
	)

	cases := []struct {
		name     string
		malleate func()
		expErr   error
	}{
		// {
		// 	name:     "success",
		// 	malleate: func() {},
		// 	expErr:   nil,
		// },
		// {
		// 	name: "invalid proof",
		// 	malleate: func() {
		// 		msg.ProofErrorReceipt = []byte("invalid proof")
		// 	},
		// 	expErr: commitmenttypes.ErrInvalidProof,
		// },
		// {
		// 	name: "invalid error receipt sequence",
		// 	malleate: func() {
		// 		const invalidSequence = 0

		// 		errorReceipt, ok := suite.chainB.GetSimApp().IBCKeeper.ChannelKeeper.GetUpgradeErrorReceipt(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
		// 		suite.Require().True(ok)

		// 		errorReceipt.Sequence = invalidSequence

		// 		// overwrite the error receipt with an invalid sequence.
		// 		suite.chainB.GetSimApp().IBCKeeper.ChannelKeeper.SetUpgradeErrorReceipt(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, errorReceipt)

		// 		// ensure that the error receipt is committed to state.
		// 		suite.coordinator.CommitBlock(suite.chainB)
		// 		suite.Require().NoError(path.EndpointA.UpdateClient())

		// 		// retrieve the error receipt proof and proof height.
		// 		errorReceiptProof, proofHeight := path.EndpointB.QueryProof(host.ChannelUpgradeErrorKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID))

		// 		// provide a valid proof of the error receipt with an invalid sequence.
		// 		msg.ErrorReceipt.Sequence = invalidSequence
		// 		msg.ProofErrorReceipt = errorReceiptProof
		// 		msg.ProofHeight = proofHeight
		// 	},
		// 	expErr: channeltypes.ErrInvalidUpgradeSequence,
		// },
		// {
		//	name: "capability not found",
		//	malleate: func() {
		//		msg.ChannelId = ibctesting.InvalidID
		//	},
		//	expErr: capabilitytypes.ErrCapabilityNotFound,
		// },
	}

	for _, tc := range cases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			path = ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(path)

			// configure the channel upgrade version on testing endpoints
			path.EndpointA.ChannelConfig.ProposedUpgrade.Fields.Version = ibcmock.UpgradeVersion
			path.EndpointB.ChannelConfig.ProposedUpgrade.Fields.Version = ibcmock.UpgradeVersion

			suite.Require().NoError(path.EndpointA.ChanUpgradeInit())

			// fetch the previous channel when it is in the INITUPGRADE state.
			prevChannel := path.EndpointA.GetChannel()

			// cause the upgrade to fail on chain b so an error receipt is written.
			suite.chainB.GetSimApp().IBCMockModule.IBCApp.OnChanUpgradeTry = func(
				ctx sdk.Context, portID, channelID string, order channeltypes.Order, connectionHops []string, counterpartyVersion string,
			) (string, error) {
				return "", fmt.Errorf("mock app callback failed")
			}

			suite.Require().NoError(path.EndpointB.ChanUpgradeTry())

			suite.Require().NoError(path.EndpointA.UpdateClient())

			upgradeErrorReceiptKey := host.ChannelUpgradeErrorKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			errorReceiptProof, proofHeight := path.EndpointB.QueryProof(upgradeErrorReceiptKey)

			errorReceipt, ok := suite.chainB.GetSimApp().IBCKeeper.ChannelKeeper.GetUpgradeErrorReceipt(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			suite.Require().True(ok)

			msg = &channeltypes.MsgChannelUpgradeCancel{
				PortId:            path.EndpointA.ChannelConfig.PortID,
				ChannelId:         path.EndpointA.ChannelID,
				ErrorReceipt:      errorReceipt,
				ProofErrorReceipt: errorReceiptProof,
				ProofHeight:       proofHeight,
				Signer:            suite.chainA.SenderAccount.GetAddress().String(),
			}

			tc.malleate()

			res, err := suite.chainA.GetSimApp().GetIBCKeeper().ChannelUpgradeCancel(suite.chainA.GetContext(), msg)

			expPass := tc.expErr == nil
			if expPass {
				suite.Require().NoError(err)
				channel := path.EndpointA.GetChannel()
				suite.Require().Equal(prevChannel.Version, channel.Version, "channel version should be reverted")
				suite.Require().Equalf(channeltypes.OPEN, channel.State, "channel state should be %s", channeltypes.OPEN.String())
				suite.Require().Equalf(channeltypes.NOTINFLUSH, channel.FlushStatus, "channel flush status should be %s", channeltypes.NOTINFLUSH.String())
				suite.Require().Equal(errorReceipt.Sequence, channel.UpgradeSequence, "channel upgrade sequence should be set to error receipt sequence")
			} else {
				suite.Require().Nil(res)
				suite.Require().ErrorIs(err, tc.expErr)

				channel := path.EndpointA.GetChannel()

				suite.Require().Equal(prevChannel.Version, channel.Version, "channel version should not be changed")
				suite.Require().Equalf(prevChannel.State, channel.State, "channel state should be %s", prevChannel.State.String())
				suite.Require().Equalf(prevChannel.FlushStatus, channel.FlushStatus, "channel flush status should be %s", prevChannel.FlushStatus.String())
				suite.Require().Equal(prevChannel.UpgradeSequence, channel.UpgradeSequence, "channel upgrade sequence should not incremented")
			}
		})
	}
}

func (suite *KeeperTestSuite) TestChannelUpgradeTimeout() {
	var (
		path *ibctesting.Path
		msg  *channeltypes.MsgChannelUpgradeTimeout
	)

	cases := []struct {
		name     string
		malleate func()
		expErr   error
	}{
		// {
		// 	name:     "success",
		// 	malleate: func() {},
		// 	expErr:   nil,
		// },
		// {
		// 	name: "success with error receipt",
		// 	malleate: func() {
		// 		// this error receipt is for a past upgrade so is stale
		// 		errReceipt := channeltypes.ErrorReceipt{
		// 			Sequence: 0,
		// 			Message:  "error",
		// 		}
		// 		path.EndpointB.Chain.App.GetIBCKeeper().ChannelKeeper.SetUpgradeErrorReceipt(path.EndpointB.Chain.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, errReceipt)

		// 		suite.Require().NoError(path.EndpointB.UpdateClient())
		// 		suite.Require().NoError(path.EndpointA.UpdateClient())

		// 		upgradeErrorReceiptKey := host.ChannelUpgradeErrorKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
		// 		errorReceiptProof, proofHeight := path.EndpointB.QueryProof(upgradeErrorReceiptKey)
		// 		counterpartyChannel := path.EndpointB.GetChannel()
		// 		// the channel proof is being regenerated as the proof height will change
		// 		channelProof, _, _ := path.EndpointA.QueryChannelUpgradeProof()

		// 		msg.PreviousErrorReceipt = &errReceipt
		// 		msg.ProofErrorReceipt = errorReceiptProof
		// 		msg.CounterpartyChannel = counterpartyChannel
		// 		msg.ProofChannel = channelProof
		// 		msg.ProofHeight = proofHeight
		// 	},
		// 	expErr: nil,
		// },
		// {
		// 	name: "invalid proof",
		// 	malleate: func() {
		// 		msg.ProofErrorReceipt = []byte("invalid proof")
		// 	},
		// 	expErr: commitmenttypes.ErrInvalidProof,
		// },
		// {
		// 	name: "invalid error receipt sequence, this error receipt is for this same upgrade so UpgradeCancel should be used instead",
		// 	malleate: func() {
		// 		errReceipt := channeltypes.ErrorReceipt{
		// 			Sequence: 1,
		// 			Message:  "error",
		// 		}
		// 		path.EndpointB.Chain.App.GetIBCKeeper().ChannelKeeper.SetUpgradeErrorReceipt(path.EndpointB.Chain.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, errReceipt)
		// 		suite.Require().NoError(path.EndpointB.UpdateClient())
		// 		suite.Require().NoError(path.EndpointA.UpdateClient())

		// 		upgradeErrorReceiptKey := host.ChannelUpgradeErrorKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
		// 		errorReceiptProof, proofHeight := path.EndpointB.QueryProof(upgradeErrorReceiptKey)
		// 		counterpartyChannel := path.EndpointB.GetChannel()
		// 		channelProof, _, _ := path.EndpointA.QueryChannelUpgradeProof()

		// 		msg.PreviousErrorReceipt = &errReceipt
		// 		msg.ProofErrorReceipt = errorReceiptProof
		// 		msg.CounterpartyChannel = counterpartyChannel
		// 		msg.ProofChannel = channelProof
		// 		msg.ProofHeight = proofHeight
		// 	},
		// 	expErr: channeltypes.ErrInvalidUpgradeSequence,
		// },
	}

	for _, tc := range cases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			path = ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(path)

			// configure the channel upgrade version on testing endpoints
			path.EndpointA.ChannelConfig.ProposedUpgrade.Fields.Version = ibcmock.UpgradeVersion
			path.EndpointA.ChannelConfig.ProposedUpgrade.Timeout.Height = clienttypes.NewHeight(1, 1)

			suite.Require().NoError(path.EndpointA.ChanUpgradeInit())

			// fetch the previous channel when it is in the INITUPGRADE state.
			prevChannel := path.EndpointA.GetChannel()

			suite.Require().NoError(path.EndpointB.UpdateClient())

			upgradeErrorReceiptKey := host.ChannelUpgradeErrorKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			errorReceiptProof, proofHeight := path.EndpointB.QueryProof(upgradeErrorReceiptKey)

			channelProof, _, _ := path.EndpointA.QueryChannelUpgradeProof()

			msg = &channeltypes.MsgChannelUpgradeTimeout{
				PortId:               path.EndpointA.ChannelConfig.PortID,
				ChannelId:            path.EndpointA.ChannelID,
				CounterpartyChannel:  path.EndpointB.GetChannel(),
				PreviousErrorReceipt: nil,
				ProofChannel:         channelProof,
				ProofErrorReceipt:    errorReceiptProof,
				ProofHeight:          proofHeight,
				Signer:               suite.chainA.SenderAccount.GetAddress().String(),
			}

			tc.malleate()

			ctx := suite.chainA.GetContext()
			res, err := suite.chainA.GetSimApp().GetIBCKeeper().ChannelUpgradeTimeout(ctx, msg)

			events := ctx.EventManager().Events().ToABCIEvents()

			expPass := tc.expErr == nil
			if expPass {
				suite.Require().NoError(err)

				channel := path.EndpointA.GetChannel()
				upgrade := path.EndpointA.GetProposedUpgrade()

				expEvents := ibctesting.EventsMap{
					channeltypes.EventTypeChannelUpgradeTimeout: {
						channeltypes.AttributeKeyPortID:                path.EndpointA.ChannelConfig.PortID,
						channeltypes.AttributeKeyChannelID:             path.EndpointA.ChannelID,
						channeltypes.AttributeCounterpartyPortID:       path.EndpointB.ChannelConfig.PortID,
						channeltypes.AttributeCounterpartyChannelID:    path.EndpointB.ChannelID,
						channeltypes.AttributeKeyUpgradeConnectionHops: upgrade.Fields.ConnectionHops[0],
						channeltypes.AttributeKeyUpgradeVersion:        upgrade.Fields.Version,
						channeltypes.AttributeKeyUpgradeOrdering:       upgrade.Fields.Ordering.String(),
						channeltypes.AttributeKeyUpgradeTimeout:        upgrade.Timeout.String(),
						channeltypes.AttributeKeyUpgradeSequence:       fmt.Sprintf("%d", channel.UpgradeSequence),
					},
					sdk.EventTypeMessage: {
						sdk.AttributeKeyModule: channeltypes.AttributeValueCategory,
					},
				}

				suite.Require().Equalf(channeltypes.OPEN, channel.State, "channel state should be %s", channeltypes.OPEN.String())
				suite.Require().Equalf(channeltypes.NOTINFLUSH, channel.FlushStatus, "channel flush status should be %s", channeltypes.NOTINFLUSH.String())

				_, found := path.EndpointA.Chain.GetSimApp().IBCKeeper.ChannelKeeper.GetUpgrade(suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				suite.Require().False(found, "channel upgrade should be nil")

				suite.Require().NotNil(res)
				suite.Require().Equal(channeltypes.SUCCESS, res.Result)
				ibctesting.AssertEvents(&suite.Suite, expEvents, events)
			} else {
				suite.Require().Nil(res)
				suite.Require().ErrorIs(err, tc.expErr)

				channel := path.EndpointA.GetChannel()
				suite.Require().Equalf(prevChannel.State, channel.State, "channel state should be %s", prevChannel.State.String())
				suite.Require().Equalf(prevChannel.FlushStatus, channel.FlushStatus, "channel flush status should be %s", prevChannel.FlushStatus.String())
				suite.Require().Equal(prevChannel.UpgradeSequence, channel.UpgradeSequence, "channel upgrade sequence should not incremented")

				_, found := path.EndpointA.Chain.GetSimApp().IBCKeeper.ChannelKeeper.GetUpgrade(suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				suite.Require().True(found, "channel upgrade should not be nil")
				ibctesting.AssertEvents(&suite.Suite, ibctesting.EventsMap{}, events)
			}
		})
	}
}

// TestUpdateClientParams tests the UpdateClientParams rpc handler
func (suite *KeeperTestSuite) TestUpdateClientParams() {
	validAuthority := suite.chainA.App.GetIBCKeeper().GetAuthority()
	testCases := []struct {
		name    string
		msg     *clienttypes.MsgUpdateParams
		expPass bool
	}{
		{
			"success: valid authority and default params",
			clienttypes.NewMsgUpdateParams(validAuthority, clienttypes.DefaultParams()),
			true,
		},
		{
			"failure: malformed authority address",
			clienttypes.NewMsgUpdateParams(ibctesting.InvalidID, clienttypes.DefaultParams()),
			false,
		},
		{
			"failure: empty authority address",
			clienttypes.NewMsgUpdateParams("", clienttypes.DefaultParams()),
			false,
		},
		{
			"failure: whitespace authority address",
			clienttypes.NewMsgUpdateParams("    ", clienttypes.DefaultParams()),
			false,
		},
		{
			"failure: unauthorized authority address",
			clienttypes.NewMsgUpdateParams(ibctesting.TestAccAddress, clienttypes.DefaultParams()),
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			_, err := keeper.Keeper.UpdateClientParams(*suite.chainA.App.GetIBCKeeper(), suite.chainA.GetContext(), tc.msg)
			if tc.expPass {
				suite.Require().NoError(err)
				p := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetParams(suite.chainA.GetContext())
				suite.Require().Equal(tc.msg.Params, p)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

// TestUpdateConnectionParams tests the UpdateConnectionParams rpc handler
func (suite *KeeperTestSuite) TestUpdateConnectionParams() {
	validAuthority := suite.chainA.App.GetIBCKeeper().GetAuthority()
	testCases := []struct {
		name    string
		msg     *connectiontypes.MsgUpdateParams
		expPass bool
	}{
		{
			"success: valid authority and default params",
			connectiontypes.NewMsgUpdateParams(validAuthority, connectiontypes.DefaultParams()),
			true,
		},
		{
			"failure: malformed authority address",
			connectiontypes.NewMsgUpdateParams(ibctesting.InvalidID, connectiontypes.DefaultParams()),
			false,
		},
		{
			"failure: empty authority address",
			connectiontypes.NewMsgUpdateParams("", connectiontypes.DefaultParams()),
			false,
		},
		{
			"failure: whitespace authority address",
			connectiontypes.NewMsgUpdateParams("    ", connectiontypes.DefaultParams()),
			false,
		},
		{
			"failure: unauthorized authority address",
			connectiontypes.NewMsgUpdateParams(ibctesting.TestAccAddress, connectiontypes.DefaultParams()),
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			_, err := keeper.Keeper.UpdateConnectionParams(*suite.chainA.App.GetIBCKeeper(), suite.chainA.GetContext(), tc.msg)
			if tc.expPass {
				suite.Require().NoError(err)
				p := suite.chainA.App.GetIBCKeeper().ConnectionKeeper.GetParams(suite.chainA.GetContext())
				suite.Require().Equal(tc.msg.Params, p)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
