// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./common.sol";

/// @title MultiSigWallet
/// @notice Basic multi-signature wallet requiring multiple confirmations
contract MultiSigWallet is ReentrancyGuard {
    event Deposit(address indexed from, uint256 amount);
    event Submit(uint indexed txId);
    event Confirm(address indexed owner, uint indexed txId);
    event Revoke(address indexed owner, uint indexed txId);
    event Execute(uint indexed txId);

    address[] public owners;
    mapping(address => bool) public isOwner;
    uint public required;

    struct Transaction {
        address to;
        uint value;
        bytes data;
        bool executed;
        uint confirmCount;
    }

    Transaction[] public transactions;
    mapping(uint => mapping(address => bool)) public confirmed;

    modifier onlyOwner() {
        require(isOwner[msg.sender], "not owner");
        _;
    }

    constructor(address[] memory _owners, uint _required) {
        require(_owners.length > 0, "owners required");
        require(_required > 0 && _required <= _owners.length, "invalid required");
        for (uint i; i < _owners.length; i++) {
            address owner = _owners[i];
            require(owner != address(0), "zero address");
            require(!isOwner[owner], "owner not unique");
            isOwner[owner] = true;
            owners.push(owner);
        }
        required = _required;
    }

    receive() external payable {
        emit Deposit(msg.sender, msg.value);
    }

    /// @custom:opcode MultiSig_Submit 0x200001
    /// @custom:gas MultiSig_Submit
    function submit(address to, uint value, bytes calldata data) external onlyOwner {
        transactions.push(Transaction(to, value, data, false, 0));
        bytes32[3] memory input;
        input[0] = bytes32(uint256(uint160(to)));
        input[1] = bytes32(value);
        input[2] = keccak256(data);
        assembly {
            // MultiSig_Submit opcode = 0x200001
            let success := call(gas(), 0x200001, 0, input, 96, 0, 0)
            if iszero(success) { revert(0, 0) }
        }
        emit Submit(transactions.length - 1);
    }

    /// @custom:opcode MultiSig_Confirm 0x200002
    /// @custom:gas MultiSig_Confirm
    function confirm(uint txId) external onlyOwner {
        require(txId < transactions.length, "tx not exist");
        Transaction storage txn = transactions[txId];
        require(!confirmed[txId][msg.sender], "already confirmed");
        confirmed[txId][msg.sender] = true;
        txn.confirmCount += 1;
        bytes32[2] memory input;
        input[0] = bytes32(txId);
        input[1] = bytes32(uint256(uint160(msg.sender)));
        assembly {
            // MultiSig_Confirm opcode = 0x200002
            let success := call(gas(), 0x200002, 0, input, 64, 0, 0)
            if iszero(success) { revert(0, 0) }
        }
        emit Confirm(msg.sender, txId);
    }

    /// @custom:opcode MultiSig_Revoke 0x200003
    /// @custom:gas MultiSig_Revoke
    function revoke(uint txId) external onlyOwner {
        require(txId < transactions.length, "tx not exist");
        Transaction storage txn = transactions[txId];
        require(confirmed[txId][msg.sender], "not confirmed");
        confirmed[txId][msg.sender] = false;
        txn.confirmCount -= 1;
        bytes32[2] memory input;
        input[0] = bytes32(txId);
        input[1] = bytes32(uint256(uint160(msg.sender)));
        assembly {
            // MultiSig_Revoke opcode = 0x200003
            let success := call(gas(), 0x200003, 0, input, 64, 0, 0)
            if iszero(success) { revert(0, 0) }
        }
        emit Revoke(msg.sender, txId);
    }

    /// @custom:opcode MultiSig_Execute 0x200004
    /// @custom:gas MultiSig_Execute
    function execute(uint txId) external onlyOwner nonReentrant {
        require(txId < transactions.length, "tx not exist");
        Transaction storage txn = transactions[txId];
        require(!txn.executed, "already executed");
        require(txn.confirmCount >= required, "not enough confirmations");
        txn.executed = true;
        bytes32[3] memory input;
        input[0] = bytes32(txId);
        input[1] = bytes32(uint256(uint160(txn.to)));
        input[2] = bytes32(txn.value);
        assembly {
            // MultiSig_Execute opcode = 0x200004
            let success := call(gas(), 0x200004, 0, input, 96, 0, 0)
            if iszero(success) { revert(0, 0) }
        }
        (bool ok, ) = txn.to.call{value: txn.value}(txn.data);
        require(ok, "tx failed");
        emit Execute(txId);
    }
}

