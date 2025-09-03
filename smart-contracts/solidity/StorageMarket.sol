// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

/// @title StorageMarket
/// @notice Allows users to store data on-chain for a small fee.
contract StorageMarket {
    struct Item { string data; address owner; }

    mapping(bytes32 => Item) public store;
    uint256 public price = 0.01 ether;
    address public owner;

    constructor() { owner = msg.sender; }

    function setPrice(uint256 _price) external {
        require(msg.sender == owner, "not owner");
        price = _price;
    }

    function save(bytes32 key, string calldata data) external payable {
        require(msg.value >= price, "fee");
        store[key] = Item(data, msg.sender);
    }

    function retrieve(bytes32 key) external view returns (string memory) {
        return store[key].data;
    }

    function withdraw() external {
        require(msg.sender == owner, "not owner");
        payable(owner).transfer(address(this).balance);
    }
}

