# 🎮 Brainrot Tamagotchi - MVP Development Plan

## 📅 Timeline: 3-4 тижні

---

## 🎯 MVP Goal

**Brainrot Tamagotchi** - Base mini app де користувачі:
- Доглядають за мемним тамагочі (NFT)
- Відкривають кейси для нових мемів
- Прокачують персонажів (Level 1-30)
- Торгують на маркетплейсі
- Обмінюють (burn) 3 NFT на шанс отримати кращий

---

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────┐
│                  Frontend                        │
│         React/Next.js + Base SDK                 │
│            RainbowKit + Wagmi                    │
└───────────────┬─────────────────────────────────┘
                │
                │ REST API + WebSocket
                │
┌───────────────▼─────────────────────────────────┐
│               Backend (Golang)                   │
│        Gin + GORM + Redis + go-ethereum         │
│                                                  │
│  ┌────────────┐  ┌──────────────┐              │
│  │ Tamagotchi │  │  Marketplace │              │
│  │  Service   │  │   Service    │              │
│  └────────────┘  └──────────────┘              │
│                                                  │
│  ┌────────────┐  ┌──────────────┐              │
│  │   Case     │  │  Blockchain  │              │
│  │  Service   │  │   Service    │              │
│  └────────────┘  └──────────────┘              │
└──────────┬────────────────┬─────────────────────┘
           │                │
           │                │ Web3 Calls
           │                │
┌──────────▼────────┐  ┌────▼────────────────────┐
│   PostgreSQL      │  │   Base Chain            │
│   (User data,     │  │   (Smart Contracts)     │
│    NFT states)    │  │                         │
└───────────────────┘  │  - BrainrotNFT.sol     │
                       │  - CaseOpening.sol      │
┌───────────────────┐  │  - Marketplace.sol     │
│      Redis        │  │  - BurnUpgrade.sol     │
│  (Hunger/Mood/    │  │                         │
│   Energy cache)   │  └─────────────────────────┘
└───────────────────┘
```

---

## 📦 Tech Stack

### **Smart Contracts**
- Solidity ^0.8.20
- Hardhat (development)
- OpenZeppelin (ERC-721, Ownable, ReentrancyGuard)
- Base Sepolia (testnet) → Base Mainnet

### **Backend**
- Go 1.21+
- Gin (web framework)
- GORM + PostgreSQL
- go-redis
- go-ethereum (blockchain interaction)
- WebSocket (gorilla)

### **Frontend**
- Next.js 14 (App Router)
- React 18
- TypeScript
- Wagmi + Viem (Base integration)
- RainbowKit (wallet connect)
- TailwindCSS + Framer Motion
- Base SDK

### **Infrastructure**
- Docker + Docker Compose
- IPFS (NFT metadata)
- Base RPC

---

## 🗂️ Project Structure

```
/unified/
├── /contracts/                    # Smart contracts
│   ├── /src/
│   │   ├── BrainrotNFT.sol
│   │   ├── CaseOpening.sol
│   │   ├── Marketplace.sol
│   │   └── BurnUpgrade.sol
│   ├── /test/
│   ├── hardhat.config.js
│   └── package.json
│
├── /backend/                      # Golang API
│   ├── /cmd/
│   │   └── main.go
│   ├── /internal/
│   │   ├── /api/                  # HTTP handlers
│   │   ├── /blockchain/           # Web3 integration
│   │   ├── /models/               # Data models
│   │   ├── /services/             # Business logic
│   │   └── /repository/           # Database layer
│   ├── /pkg/
│   ├── /contracts/                # Generated Go bindings
│   ├── go.mod
│   └── .env
│
├── /frontend/                     # Next.js app
│   ├── /app/
│   │   ├── page.tsx               # Home (Tamagotchi)
│   │   ├── /cases/
│   │   ├── /marketplace/
│   │   └── /inventory/
│   ├── /components/
│   ├── /lib/
│   ├── /public/
│   │   └── /memes/                # Meme images
│   ├── package.json
│   └── next.config.js
│
├── /docs/
│   ├── API.md
│   ├── CONTRACTS.md
│   └── DEPLOYMENT.md
│
├── docker-compose.yml
└── README.md
```

---

## 📝 Development Phases

### **Phase 1: Smart Contracts (Week 1)**

#### 1.1 BrainrotNFT.sol
```solidity
- ERC-721 implementation
- Mint function with rarity
- Burn function
- Level tracking (mapping tokenId => level)
- upgradeLevel() function (paid)
- Metadata URI (IPFS)
```

#### 1.2 CaseOpening.sol
```solidity
- Buy case (Bronze/Silver/Gold)
- Open case → random rarity
- Mint NFT through BrainrotNFT
- Payment handling (Base ETH)
```

#### 1.3 Marketplace.sol
```solidity
- List NFT (tokenId, price)
- Buy NFT (with 5% fee)
- Cancel listing
- Escrow mechanism
```

#### 1.4 BurnUpgrade.sol
```solidity
- Burn 3 NFTs of same rarity
- Random chance for higher rarity
- Mint new NFT if success
```

**Deliverables:**
- ✅ All contracts written
- ✅ Unit tests (>80% coverage)
- ✅ Deploy to Base Sepolia
- ✅ Verify on Basescan

---

### **Phase 2: Backend - Core (Week 1-2)**

#### 2.1 Setup
```go
- Initialize Go project
- Setup Gin server
- Connect PostgreSQL (GORM)
- Connect Redis
- Environment config
```

#### 2.2 Database Models
```go
type User struct {
    WalletAddress string
    CreatedAt     time.Time
}

type NFT struct {
    TokenID       uint
    Owner         string
    MemeType      string  // "pepe", "doge", etc.
    Rarity        string  // "common", "rare", etc.
    Level         int
    Hunger        int
    Mood          int
    Energy        int
    LastFed       time.Time
    LastPlayed    time.Time
}

type MarketListing struct {
    TokenID       uint
    Seller        string
    Price         float64
    IsActive      bool
}

type CaseOpening struct {
    UserAddress   string
    CaseType      string
    TokenID       uint
    Rarity        string
    TxHash        string
}
```

#### 2.3 Blockchain Integration
```go
- Connect to Base RPC
- Load contract ABIs
- Listen to NFT mint events
- Listen to marketplace events
- Transaction signing
```

#### 2.4 Services

**TamagotchiService:**
```go
- GetPetState(tokenID) - get hunger/mood/energy
- FeedPet(tokenID, isPaid) - feed (free once/day or paid)
- PlayWithPet(tokenID) - improve mood
- StartHungerDecay() - background goroutine (every 6h)
- CheckDeath() - if hunger/mood/energy = 0
```

**CaseService:**
```go
- BuyCase(userAddress, caseType) - create payment
- OpenCase(caseID) - call smart contract
- GetRandomRarity(caseType) - determine rarity
- MintNFT(userAddress, rarity) - call contract
```

**MarketplaceService:**
```go
- ListNFT(tokenID, price)
- BuyNFT(tokenID, buyerAddress)
- GetListings(filters) - browse marketplace
- CancelListing(tokenID)
```

**Deliverables:**
- ✅ REST API running on :8080
- ✅ PostgreSQL schema created
- ✅ Redis caching working
- ✅ Blockchain events listener

---

### **Phase 3: Frontend - UI (Week 2-3)**

#### 3.1 Setup
```bash
- Initialize Next.js 14
- Setup TailwindCSS
- Install Wagmi + RainbowKit
- Configure Base chain
```

#### 3.2 Core Pages

**Home (`/`):**
```tsx
- Display current pet (NFT)
- Show Hunger/Mood/Energy bars
- [Feed] [Play] [Shop] buttons
- Real-time updates (WebSocket)
```

**Cases (`/cases`):**
```tsx
- 3 case cards (Bronze, Silver, Gold)
- Price + rarity info
- Buy button → Base Pay
- Opening animation
- Reveal NFT
```

**Marketplace (`/marketplace`):**
```tsx
- Grid of listed NFTs
- Filters (rarity, level, price)
- Buy button
- "My Listings" tab
```

**Inventory (`/inventory`):**
```tsx
- Grid of user's NFTs
- Select to use as pet
- [Upgrade Level] button
- [List on Market] button
- [Burn 3 for Upgrade] section
```

#### 3.3 Components
```tsx
- PetDisplay.tsx - 3D/animated pet
- StatBar.tsx - hunger/mood/energy
- NFTCard.tsx - NFT preview
- CaseCard.tsx - case display
- WalletButton.tsx - connect wallet
```

#### 3.4 Base Integration
```tsx
- Connect wallet (RainbowKit)
- Read contracts (wagmi)
- Write contracts (transactions)
- Base Pay for purchases
```

**Deliverables:**
- ✅ All pages functional
- ✅ Wallet connection
- ✅ Contract interactions
- ✅ Responsive design

---

### **Phase 4: Integration & Polish (Week 3-4)**

#### 4.1 End-to-End Testing
- User flow: Connect wallet → Buy case → Open → Get NFT
- User flow: Feed pet → Play → Level up
- User flow: List NFT → Another user buys
- User flow: Burn 3 NFTs → Get better rarity

#### 4.2 NFT Metadata
```json
{
  "name": "Pepe #123",
  "description": "Brainrot Tamagotchi - Rare Pepe",
  "image": "ipfs://QmXxx...",
  "attributes": [
    {"trait_type": "Meme", "value": "Pepe"},
    {"trait_type": "Rarity", "value": "Rare"},
    {"trait_type": "Level", "value": 15},
    {"trait_type": "Background", "value": "Purple"}
  ]
}
```

#### 4.3 Deployment
- Backend → VPS / Cloud (Docker)
- Frontend → Vercel
- Contracts → Base Mainnet
- Database → Managed PostgreSQL
- Redis → Managed Redis

**Deliverables:**
- ✅ MVP fully functional
- ✅ Deployed to production
- ✅ Base App integration

---

## 🎨 NFT Design Plan

### **8 Base Memes × 4-5 Variants Each**

| Meme | Common (Gray) | Rare (Blue) | Epic (Purple) | Legendary (Gold) |
|------|---------------|-------------|---------------|------------------|
| Pepe | Gray bg | Blue bg | Purple gradient | Gold holographic |
| Doge | Gray bg | Cyan bg | Purple gradient | Diamond effect |
| Gigachad | Gray bg | Blue bg | Purple gradient | Gold glow |
| Wojak | Gray bg | Blue bg | Purple gradient | Gold shimmer |
| Cheems | Gray bg | Blue bg | Purple gradient | Rainbow |
| Drake | Gray bg | Blue bg | Purple gradient | Gold crown |
| Vibing Cat | Gray bg | Cyan bg | Purple gradient | Neon rainbow |
| Pikachu | Gray bg | Yellow bg | Purple gradient | Electric gold |

**Total: ~35 unique NFT variants**

---

## 💰 Pricing Structure

### **Cases:**
- Bronze: $0.50 (80% Common, 20% Rare)
- Silver: $2.00 (70% Rare, 25% Epic, 5% Legendary)
- Gold: $10.00 (60% Epic, 40% Legendary)

### **Pet Care:**
- Feed: Free (1x/day) or $0.10
- Play: Free (unlimited)
- Energy boost: $0.10

### **Level Upgrades:**
- Lv1→5: $1
- Lv5→10: $2
- Lv10→15: $3
- Lv15→20: $5
- Lv20→25: $8
- Lv25→30: $15

### **Marketplace:**
- 5% platform fee on all sales
- Prices set by users (market driven)

---

## 📊 Success Metrics

### **Week 1:**
- ✅ Smart contracts deployed & tested
- ✅ Backend API functional
- ✅ Database schema implemented

### **Week 2:**
- ✅ Frontend basic UI
- ✅ Wallet connection
- ✅ Case opening works

### **Week 3:**
- ✅ Full integration
- ✅ Marketplace functional
- ✅ Tamagotchi mechanics working

### **Week 4:**
- ✅ Polish & bug fixes
- ✅ Deploy to production
- ✅ Base App submission

---

## 🚀 Post-MVP Features (v2)

- Evolution stones
- Leaderboards (top pets, top traders)
- Social features (show off pets)
- Multiple pets per user
- Mini-games for mood/energy
- Achievements & rewards
- Referral system
- Staking mechanisms

---

## 📞 Development Workflow

### **Daily:**
1. Pull latest code
2. Check TODO list
3. Develop feature
4. Test locally
5. Commit & push
6. Update progress

### **Testing:**
- Smart contracts: Hardhat tests
- Backend: Go tests + Postman
- Frontend: Manual testing + Playwright
- Integration: Full user flows

### **Git Workflow:**
```
main
├── develop
│   ├── feature/smart-contracts
│   ├── feature/backend-api
│   ├── feature/frontend-ui
│   └── feature/integration
```

---

## ✅ MVP Checklist

### **Smart Contracts:**
- [ ] BrainrotNFT.sol (ERC-721)
- [ ] CaseOpening.sol
- [ ] Marketplace.sol
- [ ] BurnUpgrade.sol
- [ ] Tests written
- [ ] Deployed to Base Sepolia
- [ ] Deployed to Base Mainnet

### **Backend:**
- [ ] Gin server setup
- [ ] Database models (GORM)
- [ ] Blockchain integration (go-ethereum)
- [ ] TamagotchiService
- [ ] CaseService
- [ ] MarketplaceService
- [ ] REST API endpoints
- [ ] WebSocket (real-time)
- [ ] Background jobs (hunger decay)

### **Frontend:**
- [ ] Next.js setup
- [ ] Wallet connection (RainbowKit)
- [ ] Home page (pet display)
- [ ] Cases page
- [ ] Marketplace page
- [ ] Inventory page
- [ ] Contract interactions
- [ ] Animations (case opening)

### **Integration:**
- [ ] Frontend ↔ Backend API
- [ ] Backend ↔ Smart Contracts
- [ ] IPFS metadata upload
- [ ] Base Pay integration
- [ ] WebSocket real-time updates

### **Deployment:**
- [ ] Contracts on Base Mainnet
- [ ] Backend on VPS/Cloud
- [ ] Frontend on Vercel
- [ ] Database setup
- [ ] Redis setup
- [ ] Domain & SSL

---

## 🎯 Ready to Start!

**First steps:**
1. Create project structure
2. Initialize Git repo
3. Setup contracts (Hardhat)
4. Setup backend (Go)
5. Setup frontend (Next.js)

Let's build! 🚀

