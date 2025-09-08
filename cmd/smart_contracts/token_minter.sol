// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./common.sol";

/// @title TokenMinter
/// @notice Mints tokens using the MintToken opcode
contract TokenMinter is Ownable, ReentrancyGuard {
    event Minted(address indexed token, address indexed to, uint256 amount);

    /// @notice Mints tokens to a recipient
    /// @custom:opcode MintToken 0x0E0014
    /// @custom:gas MintToken
    function mint(address token, address to, uint256 amount)
        external
        onlyOwner
        nonReentrant
    {
        require(token != address(0) && to != address(0), "zero address");
        require(amount > 0, "zero amount");
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

