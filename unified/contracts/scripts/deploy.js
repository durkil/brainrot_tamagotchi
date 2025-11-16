const hre = require("hardhat");

async function main() {
  console.log("🚀 Deploying Brainrot Tamagotchi contracts to", hre.network.name);
  
  const [deployer] = await hre.ethers.getSigners();
  console.log("📝 Deploying with account:", deployer.address);
  
  const balance = await hre.ethers.provider.getBalance(deployer.address);
  console.log("💰 Account balance:", hre.ethers.formatEther(balance), "ETH\n");
  
  // 1. Deploy BrainrotNFT
  console.log("1️⃣  Deploying BrainrotNFT...");
  const BrainrotNFT = await hre.ethers.getContractFactory("BrainrotNFT");
  const nft = await BrainrotNFT.deploy();
  await nft.waitForDeployment();
  const nftAddress = await nft.getAddress();
  console.log("✅ BrainrotNFT deployed to:", nftAddress);
  
  // 2. Deploy CaseOpening
  console.log("\n2️⃣  Deploying CaseOpening...");
  const CaseOpening = await hre.ethers.getContractFactory("CaseOpening");
  const caseOpening = await CaseOpening.deploy(nftAddress);
  await caseOpening.waitForDeployment();
  const caseOpeningAddress = await caseOpening.getAddress();
  console.log("✅ CaseOpening deployed to:", caseOpeningAddress);
  
  // 3. Deploy Marketplace
  console.log("\n3️⃣  Deploying Marketplace...");
  const Marketplace = await hre.ethers.getContractFactory("Marketplace");
  const marketplace = await Marketplace.deploy(nftAddress);
  await marketplace.waitForDeployment();
  const marketplaceAddress = await marketplace.getAddress();
  console.log("✅ Marketplace deployed to:", marketplaceAddress);
  
  // 4. Deploy BurnUpgrade
  console.log("\n4️⃣  Deploying BurnUpgrade...");
  const BurnUpgrade = await hre.ethers.getContractFactory("BurnUpgrade");
  const burnUpgrade = await BurnUpgrade.deploy(nftAddress);
  await burnUpgrade.waitForDeployment();
  const burnUpgradeAddress = await burnUpgrade.getAddress();
  console.log("✅ BurnUpgrade deployed to:", burnUpgradeAddress);
  
  // 5. Authorize contracts to mint/burn NFTs
  console.log("\n5️⃣  Setting up permissions...");
  
  console.log("   Authorizing CaseOpening to mint...");
  await nft.setAuthorizedMinter(caseOpeningAddress, true);
  
  console.log("   Authorizing BurnUpgrade to mint/burn...");
  await nft.setAuthorizedMinter(burnUpgradeAddress, true);
  
  console.log("✅ Permissions set up successfully");
  
  // Summary
  console.log("\n" + "=".repeat(60));
  console.log("🎉 DEPLOYMENT COMPLETE!");
  console.log("=".repeat(60));
  console.log("\n📋 Contract Addresses:");
  console.log("   BrainrotNFT:", nftAddress);
  console.log("   CaseOpening:", caseOpeningAddress);
  console.log("   Marketplace:", marketplaceAddress);
  console.log("   BurnUpgrade:", burnUpgradeAddress);
  
  console.log("\n💾 Save these addresses to your .env file:");
  console.log(`CONTRACT_NFT_ADDRESS=${nftAddress}`);
  console.log(`CONTRACT_CASE_ADDRESS=${caseOpeningAddress}`);
  console.log(`CONTRACT_MARKETPLACE_ADDRESS=${marketplaceAddress}`);
  console.log(`CONTRACT_BURN_ADDRESS=${burnUpgradeAddress}`);
  
  // Wait for block confirmations before verification
  if (hre.network.name !== "hardhat" && hre.network.name !== "localhost") {
    console.log("\n⏳ Waiting for block confirmations before verification...");
    await nft.deploymentTransaction().wait(5);
    
    console.log("\n🔍 Verifying contracts on Basescan...");
    console.log("Run these commands:");
    console.log(`npx hardhat verify --network ${hre.network.name} ${nftAddress}`);
    console.log(`npx hardhat verify --network ${hre.network.name} ${caseOpeningAddress} ${nftAddress}`);
    console.log(`npx hardhat verify --network ${hre.network.name} ${marketplaceAddress} ${nftAddress}`);
    console.log(`npx hardhat verify --network ${hre.network.name} ${burnUpgradeAddress} ${nftAddress}`);
  }
  
  console.log("\n✨ Ready to build the most brainrot tamagotchi! 🧠🎮\n");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });

