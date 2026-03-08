// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IERC165} from "./IERC165.sol";

interface IReceiver is IERC165 {
    function onReport(bytes calldata metadata, bytes calldata report) external;
}
