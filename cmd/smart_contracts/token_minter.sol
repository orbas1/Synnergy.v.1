pragma solidity ^0.8.0;

/// @title TokenMinter
/// @notice Mints tokens using the MintToken opcode
contract TokenMinter {
    event Minted(address indexed token, address indexed to, uint256 amount);

    /// @notice Mints tokens to a recipient
    function mint(address token, address to, uint256 amount) external {
        bytes32[3] memory input;
        input[0] = bytes32(uint256(uint160(token)));
        input[1] = bytes32(uint256(uint160(to)));
        input[2] = bytes32(amount);
        assembly {
            // MintToken opcode = 0x0E0014
            let success := call(gas(), 0x0E0014, 0, input, 96, 0, 0)
            if iszero(success) { revert(0, 0) }
        }
        emit Minted(token, to, amount);
    }
}

