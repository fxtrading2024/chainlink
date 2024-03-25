// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {CrossDomainForwarder} from "../CrossDomainForwarder.sol";

import {iOVM_CrossDomainMessenger} from "../../../vendor/@eth-optimism/contracts/v0.4.7/contracts/optimistic-ethereum/iOVM/bridge/messaging/iOVM_CrossDomainMessenger.sol";

/// @title OptimismCrossDomainForwarder - L1 xDomain account representation
/// @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
/// @dev Any other L2 contract which uses this contract's address as a privileged position,
///   can be considered to be owned by the `l1Owner`
contract OptimismCrossDomainForwarder is CrossDomainForwarder {
  string public constant override typeAndVersion = "OptimismCrossDomainForwarder 1.0.0";

  /// OVM_L2CrossDomainMessenger is a precompile usually deployed to 0x4200000000000000000000000000000000000007
  address internal immutable i_crossDomainMessengerAddr;

  /// @notice creates a new Optimism xDomain Forwarder contract
  /// @param crossDomainMessengerAddr the xDomain bridge messenger (Optimism bridge L2) contract address
  /// @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
  constructor(
    iOVM_CrossDomainMessenger crossDomainMessengerAddr,
    address l1OwnerAddr
  ) CrossDomainForwarder(l1OwnerAddr) {
    i_crossDomainMessengerAddr = address(crossDomainMessengerAddr);

    // solhint-disable-next-line gas-custom-errors
    require(i_crossDomainMessengerAddr != address(0), "Invalid xDomain Messenger address");
  }

  /// @notice This is always the address of the OVM_L2CrossDomainMessenger contract
  function crossDomainMessenger() public view override returns (address) {
    return i_crossDomainMessengerAddr;
  }

  /// @notice The call MUST come from the L1 owner (via cross-chain message.) Reverts otherwise.
  modifier onlyL1Owner() override {
    // solhint-disable-next-line gas-custom-errors
    require(msg.sender == i_crossDomainMessengerAddr, "Sender is not the L2 messenger");
    // solhint-disable-next-line gas-custom-errors
    require(
      iOVM_CrossDomainMessenger(i_crossDomainMessengerAddr).xDomainMessageSender() == l1Owner(),
      "xDomain sender is not the L1 owner"
    );
    _;
  }

  /// @notice The call MUST come from the proposed L1 owner (via cross-chain message.) Reverts otherwise.
  modifier onlyProposedL1Owner() override {
    // solhint-disable-next-line gas-custom-errors
    require(msg.sender == i_crossDomainMessengerAddr, "Sender is not the L2 messenger");
    // solhint-disable-next-line gas-custom-errors
    require(
      iOVM_CrossDomainMessenger(i_crossDomainMessengerAddr).xDomainMessageSender() == s_l1PendingOwner,
      "Must be proposed L1 owner"
    );
    _;
  }
}
