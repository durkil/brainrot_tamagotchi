// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./BrainrotNFT.sol";

/**
 * @title CaseOpening
 * @dev Контракт для відкриття кейсів з NFT
 */
contract CaseOpening is Ownable, ReentrancyGuard {
    BrainrotNFT public nftContract;
    
    // Case types
    enum CaseType {
        Bronze,   // 0 - $0.50
        Silver,   // 1 - $2.00
        Gold      // 2 - $10.00
    }
    
    // Case configuration
    struct CaseConfig {
        uint256 price;
        bool active;
    }
    
    mapping(CaseType => CaseConfig) public caseConfigs;
    
    // IPFS base URI для metadata
    string public baseMetadataURI;
    
    // Events
    event CasePurchased(
        address indexed buyer,
        CaseType caseType,
        uint256 caseId
    );
    
    event CaseOpened(
        address indexed opener,
        uint256 caseId,
        uint256 tokenId,
        BrainrotNFT.Rarity rarity,
        BrainrotNFT.MemeType memeType
    );
    
    constructor(address _nftContract) Ownable(msg.sender) {
        nftContract = BrainrotNFT(_nftContract);
        
        // Default prices (approximate in ETH, adjust for Base)
        caseConfigs[CaseType.Bronze] = CaseConfig({
            price: 0.0005 ether,  // ~$0.50
            active: true
        });
        
        caseConfigs[CaseType.Silver] = CaseConfig({
            price: 0.002 ether,   // ~$2
            active: true
        });
        
        caseConfigs[CaseType.Gold] = CaseConfig({
            price: 0.01 ether,    // ~$10
            active: true
        });
        
        baseMetadataURI = "ipfs://";
    }
    
    /**
     * @dev Купити і відразу відкрити кейс
     * @param caseType Тип кейса
     */
    function buyAndOpenCase(CaseType caseType) external payable nonReentrant returns (uint256) {
        CaseConfig memory config = caseConfigs[caseType];
        require(config.active, "Case type not active");
        require(msg.value >= config.price, "Insufficient payment");
        
        // Випадкова генерація параметрів NFT
        (
            BrainrotNFT.MemeType memeType,
            BrainrotNFT.Rarity rarity,
            uint8 colorVariant
        ) = _generateRandomNFT(caseType);
        
        // Створити metadata URI
        string memory tokenURI = _buildTokenURI(memeType, rarity, colorVariant);
        
        // Mint NFT
        uint256 tokenId = nftContract.mint(
            msg.sender,
            memeType,
            rarity,
            colorVariant,
            tokenURI
        );
        
        emit CaseOpened(msg.sender, 0, tokenId, rarity, memeType);
        
        // Повернути здачу
        if (msg.value > config.price) {
            payable(msg.sender).transfer(msg.value - config.price);
        }
        
        return tokenId;
    }
    
    /**
     * @dev Генерувати випадкові параметри NFT
     */
    function _generateRandomNFT(CaseType caseType) 
        internal 
        view 
        returns (
            BrainrotNFT.MemeType memeType,
            BrainrotNFT.Rarity rarity,
            uint8 colorVariant
        ) 
    {
        // Pseudo-random (для production використати Chainlink VRF)
        uint256 randomSeed = uint256(
            keccak256(abi.encodePacked(block.timestamp, msg.sender, block.prevrandao))
        );
        
        // Визначити рідкість на основі типу кейса
        rarity = _determineRarity(caseType, randomSeed % 100);
        
        // Випадковий тип мема (0-7)
        memeType = BrainrotNFT.MemeType((randomSeed / 100) % 8);
        
        // Випадковий варіант кольору (0-4)
        colorVariant = uint8((randomSeed / 1000) % 5);
    }
    
    /**
     * @dev Визначити рідкість на основі типу кейса
     */
    function _determineRarity(CaseType caseType, uint256 roll) 
        internal 
        pure 
        returns (BrainrotNFT.Rarity) 
    {
        if (caseType == CaseType.Bronze) {
            // Bronze: 80% Common, 20% Rare
            if (roll < 80) return BrainrotNFT.Rarity.Common;
            return BrainrotNFT.Rarity.Rare;
            
        } else if (caseType == CaseType.Silver) {
            // Silver: 70% Rare, 25% Epic, 5% Legendary
            if (roll < 70) return BrainrotNFT.Rarity.Rare;
            if (roll < 95) return BrainrotNFT.Rarity.Epic;
            return BrainrotNFT.Rarity.Legendary;
            
        } else {
            // Gold: 60% Epic, 40% Legendary
            if (roll < 60) return BrainrotNFT.Rarity.Epic;
            return BrainrotNFT.Rarity.Legendary;
        }
    }
    
    /**
     * @dev Побудувати token URI
     */
    function _buildTokenURI(
        BrainrotNFT.MemeType memeType,
        BrainrotNFT.Rarity rarity,
        uint8 colorVariant
    ) internal view returns (string memory) {
        // Format: ipfs://{hash}/{memeType}_{rarity}_{colorVariant}.json
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
     * @dev Update case price
     */
    function updateCasePrice(CaseType caseType, uint256 newPrice) external onlyOwner {
        caseConfigs[caseType].price = newPrice;
    }
    
    /**
     * @dev Toggle case availability
     */
    function toggleCaseActive(CaseType caseType) external onlyOwner {
        caseConfigs[caseType].active = !caseConfigs[caseType].active;
    }
    
    /**
     * @dev Update base metadata URI
     */
    function setBaseMetadataURI(string memory newURI) external onlyOwner {
        baseMetadataURI = newURI;
    }
    
    /**
     * @dev Withdraw accumulated funds
     */
    function withdraw() external onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");
        payable(owner()).transfer(balance);
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

