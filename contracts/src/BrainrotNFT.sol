// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title BrainrotNFT
 * @dev ERC-721 NFT для Brainrot Tamagotchi
 * Кожен NFT - це мемний персонаж з рівнем та рідкістю
 */
contract BrainrotNFT is ERC721, ERC721URIStorage, Ownable {
    uint256 private _tokenIdCounter;
    
    // Rarity levels
    enum Rarity {
        Common,    // 0
        Rare,      // 1
        Epic,      // 2
        Legendary  // 3
    }
    
    // Meme types
    enum MemeType {
        Pepe,           // 0
        Doge,           // 1
        Gigachad,       // 2
        Wojak,          // 3
        Cheems,         // 4
        Drake,          // 5
        VibingCat,      // 6
        Pikachu         // 7
    }
    
    // NFT metadata
    struct TokenMetadata {
        MemeType memeType;
        Rarity rarity;
        uint8 level;           // 1-30
        uint256 mintedAt;
        uint8 colorVariant;    // 0-4 (different backgrounds)
    }
    
    // Mapping: tokenId => metadata
    mapping(uint256 => TokenMetadata) public tokenMetadata;
    
    // Authorized minters (CaseOpening, BurnUpgrade contracts)
    mapping(address => bool) public authorizedMinters;
    
    // Events
    event NFTMinted(
        uint256 indexed tokenId,
        address indexed owner,
        MemeType memeType,
        Rarity rarity,
        uint8 level
    );
    
    event LevelUpgraded(
        uint256 indexed tokenId,
        uint8 oldLevel,
        uint8 newLevel
    );
    
    event NFTBurned(uint256 indexed tokenId);
    
    constructor() ERC721("BrainrotTamagotchi", "BRAINROT") Ownable(msg.sender) {}
    
    /**
     * @dev Mint новий NFT
     * @param to Адреса отримувача
     * @param memeType Тип мема
     * @param rarity Рідкість
     * @param colorVariant Варіант кольору (0-4)
     * @param tokenURI IPFS URI для metadata
     */
    function mint(
        address to,
        MemeType memeType,
        Rarity rarity,
        uint8 colorVariant,
        string memory tokenURI
    ) public returns (uint256) {
        require(
            authorizedMinters[msg.sender] || msg.sender == owner(),
            "Not authorized to mint"
        );
        require(colorVariant < 5, "Invalid color variant");
        
        _tokenIdCounter++;
        uint256 tokenId = _tokenIdCounter;
        
        _safeMint(to, tokenId);
        _setTokenURI(tokenId, tokenURI);
        
        tokenMetadata[tokenId] = TokenMetadata({
            memeType: memeType,
            rarity: rarity,
            level: 1,
            mintedAt: block.timestamp,
            colorVariant: colorVariant
        });
        
        emit NFTMinted(tokenId, to, memeType, rarity, 1);
        
        return tokenId;
    }
    
    /**
     * @dev Upgrade рівень NFT (платний)
     * @param tokenId ID токена
     * @param newLevel Новий рівень (1-30)
     */
    function upgradeLevel(uint256 tokenId, uint8 newLevel) external payable {
        require(_ownerOf(tokenId) == msg.sender, "Not token owner");
        require(newLevel > tokenMetadata[tokenId].level, "Level must increase");
        require(newLevel <= 30, "Max level is 30");
        
        uint8 oldLevel = tokenMetadata[tokenId].level;
        
        // Ціни за upgrade
        uint256 requiredPayment = _getUpgradePrice(oldLevel, newLevel);
        require(msg.value >= requiredPayment, "Insufficient payment");
        
        tokenMetadata[tokenId].level = newLevel;
        
        emit LevelUpgraded(tokenId, oldLevel, newLevel);
        
        // Повернути здачу
        if (msg.value > requiredPayment) {
            payable(msg.sender).transfer(msg.value - requiredPayment);
        }
    }
    
    /**
     * @dev Розрахувати ціну upgrade
     */
    function _getUpgradePrice(uint8 fromLevel, uint8 toLevel) internal pure returns (uint256) {
        // Приблизні ціни (в wei, треба буде конвертувати з USD)
        if (toLevel <= 5) return 0.001 ether;      // ~$1
        if (toLevel <= 10) return 0.002 ether;     // ~$2
        if (toLevel <= 15) return 0.003 ether;     // ~$3
        if (toLevel <= 20) return 0.005 ether;     // ~$5
        if (toLevel <= 25) return 0.008 ether;     // ~$8
        return 0.015 ether;                        // ~$15
    }
    
    /**
     * @dev Спалити NFT
     */
    function burn(uint256 tokenId) external {
        require(
            authorizedMinters[msg.sender] || 
            _ownerOf(tokenId) == msg.sender ||
            msg.sender == owner(),
            "Not authorized to burn"
        );
        
        delete tokenMetadata[tokenId];
        _burn(tokenId);
        
        emit NFTBurned(tokenId);
    }
    
    /**
     * @dev Отримати метадату токена
     */
    function getTokenMetadata(uint256 tokenId) external view returns (TokenMetadata memory) {
        require(_ownerOf(tokenId) != address(0), "Token does not exist");
        return tokenMetadata[tokenId];
    }
    
    /**
     * @dev Отримати всі токени користувача
     */
    function tokensOfOwner(address owner) external view returns (uint256[] memory) {
        uint256 tokenCount = balanceOf(owner);
        uint256[] memory tokens = new uint256[](tokenCount);
        uint256 index = 0;
        
        for (uint256 i = 1; i <= _tokenIdCounter; i++) {
            if (_ownerOf(i) == owner) {
                tokens[index] = i;
                index++;
            }
        }
        
        return tokens;
    }
    
    /**
     * @dev Authorize contract to mint/burn
     */
    function setAuthorizedMinter(address minter, bool authorized) external onlyOwner {
        authorizedMinters[minter] = authorized;
    }
    
    /**
     * @dev Withdraw зібрані кошти
     */
    function withdraw() external onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");
        payable(owner()).transfer(balance);
    }
    
    /**
     * @dev Get total supply
     */
    function totalSupply() external view returns (uint256) {
        return _tokenIdCounter;
    }
    
    // Override functions
    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }
    
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}

