// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

/// @title NFTMinting
/// @notice Bare-bones ERC721-style minting contract.
contract NFTMinting {
    string public constant name = "ExampleNFT";
    string public constant symbol = "ENFT";
    uint256 public tokenCount;
    mapping(uint256 => address) public ownerOf;
    mapping(address => uint256) public balanceOf;

    event Transfer(address indexed from, address indexed to, uint256 indexed tokenId);

    function mint() external {
        tokenCount++;
        ownerOf[tokenCount] = msg.sender;
        balanceOf[msg.sender] += 1;
        emit Transfer(address(0), msg.sender, tokenCount);
    }
}

