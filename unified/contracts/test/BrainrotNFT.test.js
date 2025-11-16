const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("BrainrotNFT", function () {
  let nft;
  let owner;
  let user1;
  let user2;

  beforeEach(async function () {
    [owner, user1, user2] = await ethers.getSigners();
    
    const BrainrotNFT = await ethers.getContractFactory("BrainrotNFT");
    nft = await BrainrotNFT.deploy();
    await nft.waitForDeployment();
  });

  describe("Minting", function () {
    it("Should mint NFT with correct metadata", async function () {
      await nft.setAuthorizedMinter(owner.address, true);
      
      const tx = await nft.mint(
        user1.address,
        0, // Pepe
        0, // Common
        0, // Color variant
        "ipfs://test"
      );
      
      await tx.wait();
      
      const metadata = await nft.getTokenMetadata(1);
      expect(metadata.memeType).to.equal(0); // Pepe
      expect(metadata.rarity).to.equal(0); // Common
      expect(metadata.level).to.equal(1);
      expect(metadata.colorVariant).to.equal(0);
    });

    it("Should only allow authorized minters", async function () {
      await expect(
        nft.connect(user1).mint(user2.address, 0, 0, 0, "ipfs://test")
      ).to.be.revertedWith("Not authorized to mint");
    });
  });

  describe("Level Upgrades", function () {
    beforeEach(async function () {
      await nft.setAuthorizedMinter(owner.address, true);
      await nft.mint(user1.address, 0, 0, 0, "ipfs://test");
    });

    it("Should upgrade level with payment", async function () {
      const upgradeFee = ethers.parseEther("0.001");
      
      await nft.connect(user1).upgradeLevel(1, 5, { value: upgradeFee });
      
      const metadata = await nft.getTokenMetadata(1);
      expect(metadata.level).to.equal(5);
    });

    it("Should reject upgrade without sufficient payment", async function () {
      await expect(
        nft.connect(user1).upgradeLevel(1, 5, { value: 0 })
      ).to.be.revertedWith("Insufficient payment");
    });

    it("Should only allow owner to upgrade", async function () {
      await expect(
        nft.connect(user2).upgradeLevel(1, 5, { value: ethers.parseEther("0.001") })
      ).to.be.revertedWith("Not token owner");
    });
  });

  describe("Burning", function () {
    beforeEach(async function () {
      await nft.setAuthorizedMinter(owner.address, true);
      await nft.mint(user1.address, 0, 0, 0, "ipfs://test");
    });

    it("Should allow owner to burn", async function () {
      await nft.connect(user1).burn(1);
      
      await expect(nft.ownerOf(1)).to.be.reverted;
    });

    it("Should allow authorized burner", async function () {
      await nft.setAuthorizedMinter(user2.address, true);
      await nft.connect(user2).burn(1);
      
      await expect(nft.ownerOf(1)).to.be.reverted;
    });
  });

  describe("Token Enumeration", function () {
    it("Should return all tokens of owner", async function () {
      await nft.setAuthorizedMinter(owner.address, true);
      
      await nft.mint(user1.address, 0, 0, 0, "ipfs://test1");
      await nft.mint(user1.address, 1, 1, 1, "ipfs://test2");
      await nft.mint(user2.address, 2, 2, 2, "ipfs://test3");
      
      const tokens = await nft.tokensOfOwner(user1.address);
      expect(tokens.length).to.equal(2);
      expect(tokens[0]).to.equal(1n);
      expect(tokens[1]).to.equal(2n);
    });
  });
});

