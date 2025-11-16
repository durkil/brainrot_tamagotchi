// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "./BrainrotNFT.sol";

/**
 * @title BurnUpgrade
 * @dev Контракт для обміну (burn) 3 NFT однієї рідкості на шанс отримати вищу
 */
contract BurnUpgrade is Ownable, ReentrancyGuard {
    BrainrotNFT public nftContract;
    
    // Success chances (out of 100)
    mapping(BrainrotNFT.Rarity => uint8) public upgradeChances;
    
    // IPFS base URI
    string public baseMetadataURI;
    
    // Events
    event BurnAttempt(
        address indexed user,
        uint256[] burnedTokenIds,
        BrainrotNFT.Rarity fromRarity,
        bool success
    );
    
    event UpgradeSuccess(
        address indexed user,
        uint256 newTokenId,
        BrainrotNFT.Rarity newRarity
    );
    
    constructor(address _nftContract) Ownable(msg.sender) {
        nftContract = BrainrotNFT(_nftContract);
        
        // Default upgrade chances
        upgradeChances[BrainrotNFT.Rarity.Common] = 40;    // 40% шанс Common -> Rare
        upgradeChances[BrainrotNFT.Rarity.Rare] = 25;      // 25% шанс Rare -> Epic
        upgradeChances[BrainrotNFT.Rarity.Epic] = 15;      // 15% шанс Epic -> Legendary
        
        baseMetadataURI = "ipfs://";
    }
    
    /**
     * @dev Спалити 3 NFT за шанс отримати вищу рідкість
     * @param tokenIds Масив з 3 token IDs
     */
    function burnForUpgrade(uint256[] calldata tokenIds) external nonReentrant {
        require(tokenIds.length == 3, "Must burn exactly 3 NFTs");
        
        BrainrotNFT.Rarity rarity;
        BrainrotNFT.MemeType memeType;
        
        // Verify ownership and same rarity
        for (uint256 i = 0; i < 3; i++) {
            require(nftContract.ownerOf(tokenIds[i]) == msg.sender, "Not owner of all NFTs");
            
            BrainrotNFT.TokenMetadata memory metadata = nftContract.getTokenMetadata(tokenIds[i]);
            
            if (i == 0) {
                rarity = metadata.rarity;
                memeType = metadata.memeType; // Use first meme type
            } else {
                require(metadata.rarity == rarity, "All NFTs must have same rarity");
            }
        }
        
        // Cannot upgrade Legendary
        require(rarity != BrainrotNFT.Rarity.Legendary, "Cannot upgrade Legendary");
        
        // Burn all 3 NFTs
        for (uint256 i = 0; i < 3; i++) {
            nftContract.burn(tokenIds[i]);
        }
        
        // Determine if upgrade succeeds
        bool success = _rollUpgrade(rarity);
        
        emit BurnAttempt(msg.sender, tokenIds, rarity, success);
        
        if (success) {
            // Mint new NFT with higher rarity
            BrainrotNFT.Rarity newRarity = BrainrotNFT.Rarity(uint8(rarity) + 1);
            
            // Random color variant
            uint8 colorVariant = uint8(
                uint256(keccak256(abi.encodePacked(block.timestamp, msg.sender))) % 5
            );
            
            // Build token URI
            string memory tokenURI = _buildTokenURI(memeType, newRarity, colorVariant);
            
            // Mint
            uint256 newTokenId = nftContract.mint(
                msg.sender,
                memeType,
                newRarity,
                colorVariant,
                tokenURI
            );
            
            emit UpgradeSuccess(msg.sender, newTokenId, newRarity);
        }
        // If failed, user just loses the 3 NFTs (burned)
    }
    
    /**
     * @dev Перевірити успіх upgrade
     */
    function _rollUpgrade(BrainrotNFT.Rarity rarity) internal view returns (bool) {
        uint8 chance = upgradeChances[rarity];
        
        // Pseudo-random (для production використати Chainlink VRF)
        uint256 roll = uint256(
            keccak256(abi.encodePacked(block.timestamp, msg.sender, block.prevrandao))
        ) % 100;
        
        return roll < chance;
    }
    
    /**
     * @dev Побудувати token URI
     */
    function _buildTokenURI(
        BrainrotNFT.MemeType memeType,
        BrainrotNFT.Rarity rarity,
        uint8 colorVariant
    ) internal view returns (string memory) {
        return string(
            abi.encodePacked(
                baseMetadataURI,
                _uint2str(uint256(memeType)),
                "_",
                _uint2str(uint256(rarity)),
                "_",
                _uint2str(uint256(colorVariant)),
                ".json"
            )
        );
    }
    
    /**
     * @dev Guaranteed upgrade: burn 5 NFTs for 100% success
     */
    function burnForGuaranteedUpgrade(uint256[] calldata tokenIds) external nonReentrant {
        require(tokenIds.length == 5, "Must burn exactly 5 NFTs");
        
        BrainrotNFT.Rarity rarity;
        BrainrotNFT.MemeType memeType;
        
        // Verify ownership and same rarity
        for (uint256 i = 0; i < 5; i++) {
            require(nftContract.ownerOf(tokenIds[i]) == msg.sender, "Not owner of all NFTs");
            
            BrainrotNFT.TokenMetadata memory metadata = nftContract.getTokenMetadata(tokenIds[i]);
            
            if (i == 0) {
                rarity = metadata.rarity;
                memeType = metadata.memeType;
            } else {
                require(metadata.rarity == rarity, "All NFTs must have same rarity");
            }
        }
        
        require(rarity != BrainrotNFT.Rarity.Legendary, "Cannot upgrade Legendary");
        
        // Burn all 5 NFTs
        for (uint256 i = 0; i < 5; i++) {
            nftContract.burn(tokenIds[i]);
        }
        
        // Guaranteed upgrade
        BrainrotNFT.Rarity newRarity = BrainrotNFT.Rarity(uint8(rarity) + 1);
        
        uint8 colorVariant = uint8(
            uint256(keccak256(abi.encodePacked(block.timestamp, msg.sender))) % 5
        );
        
        string memory tokenURI = _buildTokenURI(memeType, newRarity, colorVariant);
        
        uint256 newTokenId = nftContract.mint(
            msg.sender,
            memeType,
            newRarity,
            colorVariant,
            tokenURI
        );
        
        emit BurnAttempt(msg.sender, tokenIds, rarity, true);
        emit UpgradeSuccess(msg.sender, newTokenId, newRarity);
    }
    
    /**
     * @dev Get upgrade chance for rarity
     */
    function getUpgradeChance(BrainrotNFT.Rarity rarity) external view returns (uint8) {
        return upgradeChances[rarity];
    }
    
    /**
     * @dev Update upgrade chances (owner only)
     */
    function setUpgradeChance(BrainrotNFT.Rarity rarity, uint8 chance) external onlyOwner {
        require(chance <= 100, "Chance must be <= 100");
        upgradeChances[rarity] = chance;
    }
    
    /**
     * @dev Update base metadata URI
     */
    function setBaseMetadataURI(string memory newURI) external onlyOwner {
        baseMetadataURI = newURI;
    }
    
    /**
     * @dev Helper: uint to string
     */
    function _uint2str(uint256 _i) internal pure returns (string memory) {
        if (_i == 0) {
            return "0";
        }
        uint256 j = _i;
        uint256 len;
        while (j != 0) {
            len++;
            j /= 10;
        }
        bytes memory bstr = new bytes(len);
        uint256 k = len;
        while (_i != 0) {
            k = k - 1;
            uint8 temp = (48 + uint8(_i - _i / 10 * 10));
            bytes1 b1 = bytes1(temp);
            bstr[k] = b1;
            _i /= 10;
        }
        return string(bstr);
    }
}

