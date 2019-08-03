pragma solidity ^0.4.24;

contract Multisend {
    function multiSend(address[] recipients, uint[] amounts,string payload) public payable{
        for(uint i = 0; i < recipients.length; i++) {
            recipients[i].transfer(amounts[i]);
        }
    }
}
