pragma solidity ^0.8.0;

/// @title LiquidityAdder
/// @notice Demonstrates calling the custom Liquidity_AddLiquidity opcode
contract LiquidityAdder {
    event Added(address tokenA, address tokenB, uint256 amtA, uint256 amtB);

    /// @notice Adds liquidity to a pool via Synnergy-specific opcode
    function add(address tokenA, address tokenB, uint256 amtA, uint256 amtB) external {
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

