package keeper_test

import "github.com/sullfu/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"

func (suite *KeeperTestSuite) TestParams() {
	expParams := types.DefaultParams()

	params := suite.chainA.GetSimApp().ICAHostKeeper.GetParams(suite.chainA.GetContext())
	suite.Require().Equal(expParams, params)

	expParams.HostEnabled = false
	expParams.AllowMessages = []string{"/cosmos.staking.v1beta1.MsgDelegate"}
	suite.chainA.GetSimApp().ICAHostKeeper.SetParams(suite.chainA.GetContext(), expParams)
	params = suite.chainA.GetSimApp().ICAHostKeeper.GetParams(suite.chainA.GetContext())
	suite.Require().Equal(expParams, params)
}
