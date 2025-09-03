// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

/// @title DaoGovernance
/// @notice Minimalistic DAO governance contract with proposals and voting.
contract DaoGovernance {
    struct Proposal {
        string description;
        uint256 votesFor;
        uint256 votesAgainst;
        bool executed;
        uint256 deadline;
    }

    uint256 public proposalCount;
    mapping(uint256 => Proposal) public proposals;
    mapping(uint256 => mapping(address => bool)) public voted;

    function createProposal(string calldata description, uint256 duration) external returns (uint256) {
        proposalCount++;
        proposals[proposalCount] = Proposal(description, 0, 0, false, block.timestamp + duration);
        return proposalCount;
    }

    function vote(uint256 id, bool support) external {
        Proposal storage p = proposals[id];
        require(block.timestamp < p.deadline, "ended");
        require(!voted[id][msg.sender], "voted");
        voted[id][msg.sender] = true;
        if (support) p.votesFor++;
        else p.votesAgainst++;
    }

    function execute(uint256 id) external returns (bool accepted) {
        Proposal storage p = proposals[id];
        require(block.timestamp >= p.deadline, "active");
        require(!p.executed, "done");
        p.executed = true;
        return p.votesFor > p.votesAgainst;
    }
}

