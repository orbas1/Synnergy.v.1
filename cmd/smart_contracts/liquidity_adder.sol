// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./common.sol";

/// @title LiquidityAdder
/// @notice Adds liquidity to pools via Synnergy-specific opcode
contract LiquidityAdder is Ownable, ReentrancyGuard {
    event Added(address tokenA, address tokenB, uint256 amtA, uint256 amtB);

    /// @notice Adds liquidity to a pool via Synnergy-specific opcode
    /// @custom:opcode Liquidity_AddLiquidity 0x0F0004
    /// @custom:gas AddLiquidity
    function add(address tokenA, address tokenB, uint256 amtA, uint256 amtB)
        external
        onlyOwner
        nonReentrant
    {
        require(tokenA != address(0) && tokenB != address(0), "zero address");
        require(amtA > 0 && amtB > 0, "zero amount");
        bytes32[4] memory input;
        input[0] = bytes32(uint256(uint160(tokenA)));
        input[1] = bytes32(uint256(uint160(tokenB)));
        input[2] = bytes32(amtA);
        input[3] = bytes32(amtB);
        assembly {
            // Liquidity_AddLiquidity opcode = 0x0F0004
            let success := call(gas(), 0x0F0004, 0, input, 128, 0, 0)
            if iszero(success) { revert(0, 0) }
        }
        emit Added(tokenA, tokenB, amtA, amtB);
    }
}

