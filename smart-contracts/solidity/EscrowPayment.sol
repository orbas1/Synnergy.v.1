// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

/// @title EscrowPayment
/// @notice Holds funds until an arbiter releases them to the payee or refunds the payer.
contract EscrowPayment {
    address public payer;
    address public payee;
    address public arbiter;

    constructor(address _payee, address _arbiter) payable {
        payer = msg.sender;
        payee = _payee;
        arbiter = _arbiter;
    }

    function release() external {
        require(msg.sender == arbiter, "arbiter");
        payable(payee).transfer(address(this).balance);
    }

    function refund() external {
        require(msg.sender == arbiter, "arbiter");
        payable(payer).transfer(address(this).balance);
    }
}

