// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

/// @title AIModelMarket
/// @notice Simple marketplace for listing and purchasing AI models.
contract AIModelMarket {
    struct Model { address seller; uint256 price; string uri; bool sold; }
    uint256 public modelCount;
    mapping(uint256 => Model) public models;

    event ModelListed(uint256 id, address seller, uint256 price, string uri);
    event ModelPurchased(uint256 id, address buyer);

    function listModel(uint256 price, string calldata uri) external returns (uint256) {
        modelCount++;
        models[modelCount] = Model(msg.sender, price, uri, false);
        emit ModelListed(modelCount, msg.sender, price, uri);
        return modelCount;
    }

    function purchase(uint256 id) external payable {
        Model storage m = models[id];
        require(!m.sold, "sold");
        require(msg.value >= m.price, "price");
        m.sold = true;
        payable(m.seller).transfer(m.price);
        emit ModelPurchased(id, msg.sender);
    }
}

