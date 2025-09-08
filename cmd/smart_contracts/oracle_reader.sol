// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./common.sol";

/// @title OracleReader
/// @notice Queries a Synnergy oracle via the QueryOracle opcode
contract OracleReader is Ownable, ReentrancyGuard {
    event OracleResult(bytes32 indexed key, uint256 value);

    /// @notice Fetches a value from the Oracle module
    /// @param key Identifier for the data to query
    /// @return value Returned oracle value
    /// @custom:opcode QueryOracle 0x0A0008
    /// @custom:gas QueryOracle
    function query(bytes32 key) external onlyOwner nonReentrant returns (uint256 value) {
        require(key != bytes32(0), "empty key");
        bytes32[1] memory input;
        input[0] = key;
        bytes32 out;
        assembly {
            // QueryOracle opcode = 0x0A0008
            let success := call(gas(), 0x0A0008, 0, input, 32, out, 32)
            if iszero(success) { revert(0, 0) }
        }
        value = uint256(out);
        emit OracleResult(key, value);
    }
}

