// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

interface IERC20 {
    function transfer(address to, uint256 amount) external returns (bool);
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
}

/// @title TokenFaucet
/// @notice Simple faucet dispensing a fixed amount of ERC20 tokens per request.
contract TokenFaucet {
    IERC20 public immutable token;
    address public owner;
    uint256 public dripAmount;
    uint256 public waitTime = 1 days;
    mapping(address => uint256) public lastRequest;

    constructor(address _token, uint256 _dripAmount) {
        token = IERC20(_token);
        owner = msg.sender;
        dripAmount = _dripAmount;
    }

    /// @notice Request tokens from the faucet. Can be called once per wait period.
    function requestTokens() external {
        require(block.timestamp - lastRequest[msg.sender] >= waitTime, "wait");
        lastRequest[msg.sender] = block.timestamp;
        require(token.transfer(msg.sender, dripAmount), "transfer failed");
    }

    /// @notice Owner can deposit additional tokens into the faucet.
    function deposit(uint256 amount) external {
        require(token.transferFrom(msg.sender, address(this), amount), "deposit failed");
    }

    /// @notice Owner may change the drip amount.
    function setDripAmount(uint256 amount) external {
        require(msg.sender == owner, "not owner");
        dripAmount = amount;
    }

    /// @notice Owner may withdraw tokens from the faucet.
    function withdraw(uint256 amount) external {
        require(msg.sender == owner, "not owner");
        require(token.transfer(owner, amount), "withdraw failed");
    }
}

