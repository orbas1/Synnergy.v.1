// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./common.sol";

/// @title CrossChainETHBridge
/// @notice Locks and releases Ether across chains using signed messages
contract CrossChainETHBridge is Ownable, ReentrancyGuard {
    /// @dev address allowed to sign off-chain withdrawal messages
    address public signer;
    /// @dev track processed nonces to prevent replay
    mapping(uint256 => bool) public processed;

    event Deposit(address indexed from, uint256 amount, string target);
    event Withdraw(address indexed to, uint256 amount, string source, uint256 nonce);
    event SignerUpdated(address indexed newSigner);

    constructor(address _signer) {
        require(_signer != address(0), "signer zero");
        signer = _signer;
    }

    /// @notice update signer address
    function setSigner(address newSigner) external onlyOwner {
        require(newSigner != address(0), "signer zero");
        signer = newSigner;
        emit SignerUpdated(newSigner);
    }

    /// @notice Lock Ether and specify the target chain or address
    /// @custom:opcode LockAndMint 0x090004
    /// @custom:gas LockAndMint
    function deposit(string calldata target) external payable nonReentrant {
        require(msg.value > 0, "no value");
        emit Deposit(msg.sender, msg.value, target);
        // relay deposit to Synnergy VM for cross-chain tracking
        bytes32[3] memory input;
        input[0] = bytes32(uint256(uint160(msg.sender)));
        input[1] = bytes32(msg.value);
        input[2] = keccak256(abi.encodePacked(target));
        assembly {
            // LockAndMint opcode = 0x090004
            let success := call(gas(), 0x090004, 0, input, 96, 0, 0)
            if iszero(success) { revert(0, 0) }
        }
    }

    /// @notice Release Ether back to a user after off-chain validation
    /// @param to recipient on this chain
    /// @param amount amount of Ether to release
    /// @param source description of the source chain
    /// @param nonce unique nonce of the withdrawal message
    /// @param sig signer signature authorizing the withdrawal
    /// @custom:opcode BurnAndRelease 0x090005
    /// @custom:gas BurnAndRelease
    function withdraw(
        address payable to,
        uint256 amount,
        string calldata source,
        uint256 nonce,
        bytes calldata sig
    ) external nonReentrant {
        require(!processed[nonce], "nonce used");
        require(address(this).balance >= amount, "insufficient");
        bytes32 msgHash = keccak256(abi.encodePacked(to, amount, source, nonce, address(this)));
        require(recoverSigner(msgHash, sig) == signer, "bad sig");
        processed[nonce] = true;
        emit Withdraw(to, amount, source, nonce);
        bytes32[4] memory input;
        input[0] = bytes32(uint256(uint160(address(to))));
        input[1] = bytes32(amount);
        input[2] = keccak256(abi.encodePacked(source));
        input[3] = bytes32(nonce);
        assembly {
            // BurnAndRelease opcode = 0x090005
            let success := call(gas(), 0x090005, 0, input, 128, 0, 0)
            if iszero(success) { revert(0, 0) }
        }
        to.transfer(amount);
    }

    function recoverSigner(bytes32 msgHash, bytes calldata sig) internal pure returns (address) {
        require(sig.length == 65, "bad sig length");
        bytes32 r;
        bytes32 s;
        uint8 v;
        assembly {
            r := calldataload(sig.offset)
            s := calldataload(add(sig.offset, 32))
            v := byte(0, calldataload(add(sig.offset, 64)))
        }
        if (v < 27) v += 27;
        require(v == 27 || v == 28, "bad v");
        bytes32 ethHash = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", msgHash));
        return ecrecover(ethHash, v, r, s);
    }
}

