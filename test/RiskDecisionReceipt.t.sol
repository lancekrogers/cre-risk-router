// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "forge-std/Test.sol";
import "../contracts/evm/src/RiskDecisionReceipt.sol";

contract RiskDecisionReceiptTest is Test {
    RiskDecisionReceipt public receipt;

    bytes32 constant RUN_ID = keccak256("test-run-1");
    bytes32 constant DECISION_HASH = keccak256("decision-data");

    function setUp() public {
        receipt = new RiskDecisionReceipt(address(this));
    }

    function test_RecordApprovedDecision() public {
        receipt.recordDecision(
            RUN_ID,
            DECISION_HASH,
            true,       // approved
            1000e6,     // maxPositionUsd (6 decimals)
            50,         // maxSlippageBps
            300,        // ttlSeconds
            2000e8      // chainlinkPrice (8 decimals)
        );

        assertEq(receipt.totalDecisions(), 1);
        assertEq(receipt.totalApproved(), 1);
        assertEq(receipt.totalDenied(), 0);
        assertTrue(receipt.recorded(RUN_ID));
        assertTrue(receipt.isDecisionValid(RUN_ID));
        assertEq(receipt.getRunCount(), 1);
    }

    function test_RecordDeniedDecision() public {
        receipt.recordDecision(
            RUN_ID,
            DECISION_HASH,
            false,      // denied
            0,
            0,
            300,
            2000e8
        );

        assertEq(receipt.totalDecisions(), 1);
        assertEq(receipt.totalApproved(), 0);
        assertEq(receipt.totalDenied(), 1);
        assertTrue(receipt.recorded(RUN_ID));
    }

    function test_RejectDuplicateRunId() public {
        receipt.recordDecision(RUN_ID, DECISION_HASH, true, 1000e6, 50, 300, 2000e8);

        vm.expectRevert("Decision already recorded");
        receipt.recordDecision(RUN_ID, DECISION_HASH, true, 1000e6, 50, 300, 2000e8);
    }

    function test_DecisionExpiry() public {
        receipt.recordDecision(RUN_ID, DECISION_HASH, true, 1000e6, 50, 300, 2000e8);

        assertTrue(receipt.isDecisionValid(RUN_ID));

        vm.warp(block.timestamp + 301);

        assertFalse(receipt.isDecisionValid(RUN_ID));
    }
}
