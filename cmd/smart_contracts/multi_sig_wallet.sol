pragma solidity ^0.8.0;

/// @title MultiSigWallet
/// @notice Basic multi-signature wallet requiring multiple confirmations
contract MultiSigWallet {
    event Deposit(address indexed from, uint256 amount);
    event Submit(uint indexed txId);
    event Confirm(address indexed owner, uint indexed txId);
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

    function submit(address to, uint value, bytes calldata data) external onlyOwner {
        transactions.push(Transaction(to, value, data, false, 0));
        emit Submit(transactions.length - 1);
    }

    function confirm(uint txId) external onlyOwner {
        Transaction storage txn = transactions[txId];
        require(!confirmed[txId][msg.sender], "already confirmed");
        confirmed[txId][msg.sender] = true;
        txn.confirmCount += 1;
        emit Confirm(msg.sender, txId);
    }

    function execute(uint txId) external onlyOwner {
        Transaction storage txn = transactions[txId];
        require(!txn.executed, "already executed");
        require(txn.confirmCount >= required, "not enough confirmations");
        txn.executed = true;
        (bool ok, ) = txn.to.call{value: txn.value}(txn.data);
        require(ok, "tx failed");
        emit Execute(txId);
    }
}

