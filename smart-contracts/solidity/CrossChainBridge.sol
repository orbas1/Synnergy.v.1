// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

interface IERC20 {
    function transfer(address to, uint256 amount) external returns (bool);
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
}

/// @title CrossChainBridge
/// @notice Locks tokens on this chain and releases them under admin control to mimic a bridge.
contract CrossChainBridge {
    IERC20 public immutable token;
    address public admin;

    event Locked(address indexed from, uint256 amount, string targetChainAddress);
    event Released(address indexed to, uint256 amount);

    constructor(address _token) {
        token = IERC20(_token);
        admin = msg.sender;
    }

    function lock(uint256 amount, string calldata target) external {
        require(token.transferFrom(msg.sender, address(this), amount), "lock failed");
        emit Locked(msg.sender, amount, target);
    }

    function release(address to, uint256 amount) external {
        require(msg.sender == admin, "admin");
        require(token.transfer(to, amount), "release failed");
        emit Released(to, amount);
    }
}

