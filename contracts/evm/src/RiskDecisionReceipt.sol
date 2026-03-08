// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IReceiver} from "./interfaces/IReceiver.sol";
import {IERC165} from "./interfaces/IERC165.sol";

contract RiskDecisionReceipt is IReceiver {
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

    address public forwarder;

    event DecisionRecorded(
        bytes32 indexed runId,
        bytes32 indexed decisionHash,
        address indexed recorder,
        bool approved
    );

    constructor(address _forwarder) {
        forwarder = _forwarder;
    }

    /// @notice IReceiver entry point called by the CRE KeystoneForwarder.
    function onReport(bytes calldata, bytes calldata report) external override {
        (
            bytes32 runId,
            bytes32 decisionHash,
            bool approved,
            uint256 maxPositionUsd,
            uint256 maxSlippageBps,
            uint256 ttlSeconds,
            uint256 chainlinkPrice
        ) = abi.decode(report, (bytes32, bytes32, bool, uint256, uint256, uint256, uint256));

        _recordDecision(runId, decisionHash, approved, maxPositionUsd, maxSlippageBps, ttlSeconds, chainlinkPrice);
    }

    /// @notice Direct entry point for testing without CRE forwarder.
    function recordDecision(
        bytes32 runId,
        bytes32 decisionHash,
        bool approved,
        uint256 maxPositionUsd,
        uint256 maxSlippageBps,
        uint256 ttlSeconds,
        uint256 chainlinkPrice
    ) external {
        _recordDecision(runId, decisionHash, approved, maxPositionUsd, maxSlippageBps, ttlSeconds, chainlinkPrice);
    }

    function _recordDecision(
        bytes32 runId,
        bytes32 decisionHash,
        bool approved,
        uint256 maxPositionUsd,
        uint256 maxSlippageBps,
        uint256 ttlSeconds,
        uint256 chainlinkPrice
    ) internal {
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

    /// @notice ERC165 interface detection.
    function supportsInterface(bytes4 interfaceId) external pure override returns (bool) {
        return interfaceId == type(IReceiver).interfaceId || interfaceId == type(IERC165).interfaceId;
    }
}
