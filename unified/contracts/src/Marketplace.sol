// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "./BrainrotNFT.sol";

/**
 * @title Marketplace
 * @dev P2P marketplace для торгівлі Brainrot NFT
 */
contract Marketplace is Ownable, ReentrancyGuard {
    BrainrotNFT public nftContract;
    
    // Platform fee (5% = 500 basis points)
    uint256 public platformFeeBps = 500;
    uint256 public constant MAX_FEE_BPS = 1000; // Max 10%
    
    // Listing structure
    struct Listing {
        address seller;
        uint256 price;
        bool active;
        uint256 listedAt;
    }
    
    // Mapping: tokenId => Listing
    mapping(uint256 => Listing) public listings;
    
    // Events
    event NFTListed(
        uint256 indexed tokenId,
        address indexed seller,
        uint256 price
    );
    
    event NFTSold(
        uint256 indexed tokenId,
        address indexed seller,
        address indexed buyer,
        uint256 price,
        uint256 platformFee
    );
    
    event ListingCancelled(
        uint256 indexed tokenId,
        address indexed seller
    );
    
    event PriceUpdated(
        uint256 indexed tokenId,
        uint256 oldPrice,
        uint256 newPrice
    );
    
    constructor(address _nftContract) Ownable(msg.sender) {
        nftContract = BrainrotNFT(_nftContract);
    }
    
    /**
     * @dev Виставити NFT на продаж
     * @param tokenId ID токена
     * @param price Ціна в wei
     */
    function listNFT(uint256 tokenId, uint256 price) external {
        require(nftContract.ownerOf(tokenId) == msg.sender, "Not token owner");
        require(price > 0, "Price must be greater than 0");
        require(!listings[tokenId].active, "Already listed");
        
        // Перевірити, що контракт має approval
        require(
            nftContract.getApproved(tokenId) == address(this) ||
            nftContract.isApprovedForAll(msg.sender, address(this)),
            "Marketplace not approved"
        );
        
        listings[tokenId] = Listing({
            seller: msg.sender,
            price: price,
            active: true,
            listedAt: block.timestamp
        });
        
        emit NFTListed(tokenId, msg.sender, price);
    }
    
    /**
     * @dev Купити NFT з marketplace
     * @param tokenId ID токена
     */
    function buyNFT(uint256 tokenId) external payable nonReentrant {
        Listing memory listing = listings[tokenId];
        
        require(listing.active, "NFT not listed");
        require(msg.value >= listing.price, "Insufficient payment");
        require(msg.sender != listing.seller, "Cannot buy own NFT");
        
        // Verify ownership
        require(nftContract.ownerOf(tokenId) == listing.seller, "Seller no longer owns NFT");
        
        // Calculate fees
        uint256 platformFee = (listing.price * platformFeeBps) / 10000;
        uint256 sellerProceeds = listing.price - platformFee;
        
        // Mark as inactive
        listings[tokenId].active = false;
        
        // Transfer NFT
        nftContract.safeTransferFrom(listing.seller, msg.sender, tokenId);
        
        // Transfer payments
        payable(listing.seller).transfer(sellerProceeds);
        // Platform fee stays in contract
        
        emit NFTSold(tokenId, listing.seller, msg.sender, listing.price, platformFee);
        
        // Return excess payment
        if (msg.value > listing.price) {
            payable(msg.sender).transfer(msg.value - listing.price);
        }
    }
    
    /**
     * @dev Скасувати listing
     * @param tokenId ID токена
     */
    function cancelListing(uint256 tokenId) external {
        Listing memory listing = listings[tokenId];
        
        require(listing.active, "NFT not listed");
        require(listing.seller == msg.sender, "Not the seller");
        
        listings[tokenId].active = false;
        
        emit ListingCancelled(tokenId, msg.sender);
    }
    
    /**
     * @dev Оновити ціну listing
     * @param tokenId ID токена
     * @param newPrice Нова ціна
     */
    function updatePrice(uint256 tokenId, uint256 newPrice) external {
        Listing storage listing = listings[tokenId];
        
        require(listing.active, "NFT not listed");
        require(listing.seller == msg.sender, "Not the seller");
        require(newPrice > 0, "Price must be greater than 0");
        
        uint256 oldPrice = listing.price;
        listing.price = newPrice;
        
        emit PriceUpdated(tokenId, oldPrice, newPrice);
    }
    
    /**
     * @dev Отримати інформацію про listing
     */
    function getListing(uint256 tokenId) external view returns (Listing memory) {
        return listings[tokenId];
    }
    
    /**
     * @dev Перевірити, чи виставлений NFT
     */
    function isListed(uint256 tokenId) external view returns (bool) {
        return listings[tokenId].active;
    }
    
    /**
     * @dev Отримати всі активні listings (off-chain helper)
     * Note: Для production використати The Graph або backend indexer
     */
    function getActiveListings(uint256 maxSupply) external view returns (uint256[] memory) {
        uint256 activeCount = 0;
        
        // Count active listings
        for (uint256 i = 1; i <= maxSupply; i++) {
            if (listings[i].active) {
                activeCount++;
            }
        }
        
        // Populate array
        uint256[] memory activeTokenIds = new uint256[](activeCount);
        uint256 index = 0;
        
        for (uint256 i = 1; i <= maxSupply; i++) {
            if (listings[i].active) {
                activeTokenIds[index] = i;
                index++;
            }
        }
        
        return activeTokenIds;
    }
    
    /**
     * @dev Update platform fee (owner only)
     * @param newFeeBps New fee in basis points
     */
    function setPlatformFee(uint256 newFeeBps) external onlyOwner {
        require(newFeeBps <= MAX_FEE_BPS, "Fee too high");
        platformFeeBps = newFeeBps;
    }
    
    /**
     * @dev Withdraw accumulated platform fees
     */
    function withdrawFees() external onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No fees to withdraw");
        payable(owner()).transfer(balance);
    }
    
    /**
     * @dev Emergency: cancel any listing (owner only)
     */
    function emergencyCancelListing(uint256 tokenId) external onlyOwner {
        listings[tokenId].active = false;
        emit ListingCancelled(tokenId, listings[tokenId].seller);
    }
}

