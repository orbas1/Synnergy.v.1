// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

/// @title MultisigWallet
/// @notice Basic N-of-M multisignature wallet.
contract MultisigWallet {
    address[] public owners;
    uint256 public required;

    struct Transaction {
        address to;
        uint256 value;
        bytes data;
        bool executed;
        uint256 confirmations;
    }

    Transaction[] public transactions;
    mapping(uint256 => mapping(address => bool)) public approved;

    modifier onlyOwner() {
        bool isOwner = false;
        for (uint256 i = 0; i < owners.length; i++) {
            if (owners[i] == msg.sender) isOwner = true;
        }
        require(isOwner, "not owner");
        _;
    }

    constructor(address[] memory _owners, uint256 _required) {
        owners = _owners;
        required = _required;
    }

    receive() external payable {}

    function submit(address to, uint256 value, bytes calldata data) external onlyOwner returns (uint256) {
        transactions.push(Transaction(to, value, data, false, 0));
        return transactions.length - 1;
    }

    function approve(uint256 id) external onlyOwner {
        Transaction storage txn = transactions[id];
        require(!approved[id][msg.sender], "approved");
        approved[id][msg.sender] = true;
        txn.confirmations += 1;
        if (txn.confirmations >= required) {
            execute(id);
        }
    }

    function execute(uint256 id) public onlyOwner {
        Transaction storage txn = transactions[id];
        require(!txn.executed, "executed");
        require(txn.confirmations >= required, "insufficient approvals");
        txn.executed = true;
        (bool ok, ) = txn.to.call{value: txn.value}(txn.data);
        require(ok, "call failed");
    }
}

