pragma solidity ^0.8.0;

/// @title CrossChainETHBridge
/// @notice Simple Ether bridge emitting events for off-chain relayers
contract CrossChainETHBridge {
    /// @dev Emitted when Ether is locked on this chain to be released elsewhere
    event Deposit(address indexed from, uint256 amount, string target);

    /// @dev Emitted when Ether is withdrawn back to this chain
    event Withdraw(address indexed to, uint256 amount, string source);

    /// @notice Lock Ether and specify the target chain or address
    function deposit(string calldata target) external payable {
        require(msg.value > 0, "no value");
        emit Deposit(msg.sender, msg.value, target);
    }

    /// @notice Release Ether back to a user after off-chain validation
    /// @dev In a production bridge this would require cryptographic proofs
    function withdraw(address payable to, uint256 amount, string calldata source) external {
        require(address(this).balance >= amount, "insufficient");
        emit Withdraw(to, amount, source);
        to.transfer(amount);
    }
}

