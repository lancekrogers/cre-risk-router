// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "forge-std/Script.sol";
import "../src/RiskDecisionReceipt.sol";

contract DeployScript is Script {
    function run() external {
        // Sepolia KeystoneForwarder
        address forwarder = 0x15fC6ae953E024d975e77382eEeC56A9101f9F88;

        vm.startBroadcast();
        RiskDecisionReceipt receipt = new RiskDecisionReceipt(forwarder);
        vm.stopBroadcast();

        console.log("Deployed RiskDecisionReceipt at:", address(receipt));
    }
}
