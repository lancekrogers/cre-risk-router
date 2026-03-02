// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract RiskDecisionReceipt {
    struct Decision {
        bytes32 runId;
        bytes32 decisionHash;
        bool approved;
        uint256 maxPositionUsd;   // 6-decimal precision
        uint256 maxSlippageBps;
        uint256 ttlSeconds;
        uint256 chainlinkPrice;   // 8-decimal precision
        uint256 timestamp;
        address recorder;
    }

    mapping(bytes32 => Decision) public decisions;
    mapping(bytes32 => bool) public recorded;
    bytes32[] public runIds;

    uint256 public totalDecisions;
    uint256 public totalApproved;
    uint256 public totalDenied;

    event DecisionRecorded(
        bytes32 indexed runId,
        bytes32 indexed decisionHash,
        address indexed recorder,
        bool approved
    );

    function recordDecision(
        bytes32 runId,
        bytes32 decisionHash,
        bool approved,
        uint256 maxPositionUsd,
        uint256 maxSlippageBps,
        uint256 ttlSeconds,
        uint256 chainlinkPrice
    ) external {
        require(!recorded[runId], "Decision already recorded");

        decisions[runId] = Decision({
            runId: runId,
            decisionHash: decisionHash,
            approved: approved,
            maxPositionUsd: maxPositionUsd,
            maxSlippageBps: maxSlippageBps,
            ttlSeconds: ttlSeconds,
            chainlinkPrice: chainlinkPrice,
            timestamp: block.timestamp,
            recorder: msg.sender
        });

        recorded[runId] = true;
        runIds.push(runId);
        totalDecisions++;

        if (approved) {
            totalApproved++;
        } else {
            totalDenied++;
        }

        emit DecisionRecorded(runId, decisionHash, msg.sender, approved);
    }

    function isDecisionValid(bytes32 runId) external view returns (bool) {
        if (!recorded[runId]) {
            return false;
        }
        Decision storage d = decisions[runId];
        return block.timestamp <= d.timestamp + d.ttlSeconds;
    }

    function getRunCount() external view returns (uint256) {
        return runIds.length;
    }
}
