var ConvertLib = artifacts.require("./ConvertLib.sol");
var MetaCoin = artifacts.require("./MetaCoin.sol");
var HelloWorld = artifacts.require("./HelloWorld.sol");
var Dividends = artifacts.require("./Dividends.sol");

module.exports = function(deployer) {
  // deployer.deploy(ConvertLib);
  // deployer.link(ConvertLib, MetaCoin);
  // deployer.deploy(MetaCoin);
  // deployer.deploy(HelloWorld);
  deployer.deploy(Dividends,0xffcf8fdee72ac11b5c542428b35eef5769c409f0);
};
