# 🎮 Brainrot Tamagotchi

**Base mini app** - Meмний тамагочі з NFT скінами, кейсами та маркетплейсом

---

## 🌟 Що це?

**Brainrot Tamagotchi** - це Web3 гра на Base blockchain, де ти:
- 🐸 Доглядаєш за мемним тамагочі (NFT персонаж)
- 🎁 Відкриваєш кейси для нових скінів
- ⬆️ Прокачуєш персонажа до Level 30
- 🛒 Торгуєш на маркетплейсі
- 🔥 Обмінюєш 3 NFT на шанс отримати кращий

---

## 🏗️ Архітектура

```
/brainrot tamagotchi/
├── /contracts/        # Solidity смарт-контракти (Base chain)
├── /backend/          # Golang API (Gin + GORM + go-ethereum)
├── /frontend/         # Next.js mini app (React + Base SDK)
└── /docs/            # Документація
```

---

## 🛠️ Tech Stack

### **Blockchain:**
- Base (Ethereum L2)
- Solidity ^0.8.20
- Hardhat
- OpenZeppelin

### **Backend:**
- Go 1.21+
- Gin (web framework)
- PostgreSQL + GORM
- Redis (caching)
- go-ethereum (Web3)

### **Frontend:**
- Next.js 14
- React 18 + TypeScript
- Wagmi + Viem
- RainbowKit
- TailwindCSS

---

## 🚀 Quick Start

### **1. Clone & Install**

```bash
cd "/Users/durkil/brainrot tamagotchi"

# Contracts
cd contracts
npm install

# Backend
cd ../backend
go mod download

# Frontend
cd ../frontend
npm install
```

### **2. Setup Environment**

```bash
# Backend (.env)
BASE_RPC_URL=https://sepolia.base.org
DATABASE_URL=postgresql://user:pass@localhost:5432/brainrot
REDIS_URL=redis://localhost:6379
CONTRACT_NFT_ADDRESS=0x...
CONTRACT_CASE_ADDRESS=0x...
CONTRACT_MARKETPLACE_ADDRESS=0x...

# Frontend (.env.local)
NEXT_PUBLIC_BASE_RPC=https://sepolia.base.org
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_NFT_CONTRACT=0x...
```

### **3. Run Development**

```bash
# Terminal 1: Backend
cd backend
go run cmd/main.go

# Terminal 2: Frontend
cd frontend
npm run dev

# Terminal 3: Contracts (local node)
cd contracts
npx hardhat node
```

---

## 📦 What's Included

### **Smart Contracts:**
- ✅ **BrainrotNFT.sol** - ERC-721 NFT з рівнями та рідкістю
- ✅ **CaseOpening.sol** - Система кейсів (Bronze, Silver, Gold)
- ✅ **Marketplace.sol** - P2P торгівля NFT
- ✅ **BurnUpgrade.sol** - Обмін 3 NFT на кращий

### **Backend API:**
- ✅ **Tamagotchi Service** - Hunger/Mood/Energy mechanics
- ✅ **Case Service** - Купівля і відкриття кейсів
- ✅ **Marketplace Service** - List/Buy NFT
- ✅ **Blockchain Service** - Взаємодія з Base chain
- ✅ **Background Jobs** - Автоматичний decay метрик

### **Frontend:**
- ✅ **Home** - Твій мемний тамагочі
- ✅ **Cases** - Купуй і відкривай кейси
- ✅ **Marketplace** - Browse і торгуй NFT
- ✅ **Inventory** - Твоя колекція мемів

---

## 🎨 NFT Меми

### **8 базових мемів:**
1. 🐸 **Pepe** (Common → Legendary)
2. 🐕 **Doge** (Rare → Legendary)
3. 💪 **Gigachad** (Rare → Legendary)
4. 😢 **Wojak** (Common → Epic)
5. 🐕 **Cheems** (Common → Rare)
6. 🎤 **Drake** (Rare → Epic)
7. 😺 **Vibing Cat** (Common → Legendary)
8. ⚡ **Pikachu** (Rare → Epic)

Кожен мем має **4-5 варіантів** (різні фони/кольори/ефекти)

**Всього: ~35 унікальних NFT**

---

## 💰 Економіка

### **Кейси:**
- 🥉 Bronze: $0.50
- 🥈 Silver: $2.00
- 🥇 Gold: $10.00

### **Догляд:**
- Feed: Free (1x/day) або $0.10
- Energy: $0.10

### **Апгрейд:**
- Level up: $1-15 (залежно від рівня)

### **Marketplace:**
- 5% platform fee

---

## 📖 Documentation

- [MVP Plan](./MVP_PLAN.md) - Повний план розробки
- [Smart Contracts](./docs/CONTRACTS.md) - Опис контрактів
- [API Documentation](./docs/API.md) - Backend endpoints
- [Deployment Guide](./docs/DEPLOYMENT.md) - Як задеплоїти

---

## 🧪 Testing

```bash
# Smart contracts
cd contracts
npx hardhat test

# Backend
cd backend
go test ./...

# Frontend
cd frontend
npm run test
```

---

## 🚀 Deployment

### **Testnet (Base Sepolia):**
```bash
cd contracts
npx hardhat run scripts/deploy.js --network base-sepolia
```

### **Mainnet (Base):**
```bash
# Deploy contracts
npx hardhat run scripts/deploy.js --network base-mainnet

# Deploy backend (Docker)
cd backend
docker build -t brainrot-backend .
docker run -p 8080:8080 brainrot-backend

# Deploy frontend (Vercel)
cd frontend
vercel --prod
```

---

## 📊 Status

**Current Phase:** 🟡 Development

- [x] MVP Plan
- [ ] Smart Contracts
- [ ] Backend API
- [ ] Frontend UI
- [ ] Integration
- [ ] Testing
- [ ] Deployment

---

## 🤝 Contributing

1. Fork the repo
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

---

## 📝 License

MIT License - see [LICENSE](./LICENSE)

---

## 🔗 Links

- **Base Chain:** https://www.base.org/
- **Documentation:** [docs/](./docs/)
- **Issues:** [GitHub Issues](#)

---

## 🎯 Roadmap

### **MVP (Week 1-4):**
- ✅ Core mechanics
- ✅ 8 memes
- ✅ Cases & Marketplace
- ✅ Base integration

### **v2 (Future):**
- [ ] Evolution stones
- [ ] Leaderboards
- [ ] Social features
- [ ] Mini-games
- [ ] Achievements
- [ ] Staking

---

Made with 💙 for Base ecosystem

**Let's build the most brainrot tamagotchi ever! 🧠🎮**

