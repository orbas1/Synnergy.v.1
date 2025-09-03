// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

/// @title RegulatoryCompliance
/// @notice Maintains a registry of compliant addresses.
contract RegulatoryCompliance {
    address public regulator;
    mapping(address => bool) public approved;

    constructor() { regulator = msg.sender; }

    function setApproval(address user, bool status) external {
        require(msg.sender == regulator, "only regulator");
        approved[user] = status;
    }

    function isApproved(address user) external view returns (bool) {
        return approved[user];
    }
}

