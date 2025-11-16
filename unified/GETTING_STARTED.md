# 🚀 Getting Started - Brainrot Tamagotchi

Ласкаво просимо до **Brainrot Tamagotchi**! Цей гайд допоможе тобі швидко запустити проект локально.

---

## 📦 Що вже зроблено

✅ **Smart Contracts** (Solidity)
- BrainrotNFT.sol - ERC-721 NFT з levels та рідкістю
- CaseOpening.sol - Система кейсів
- Marketplace.sol - P2P торгівля
- BurnUpgrade.sol - Обмін NFT

✅ **Backend API** (Golang)
- Gin web server
- PostgreSQL + GORM
- Redis caching
- go-ethereum blockchain integration
- Tamagotchi service (hunger/mood/energy)
- Marketplace service
- Background jobs

✅ **Frontend** (Next.js)
- React 18 + TypeScript
- RainbowKit wallet integration
- Wagmi for Base chain
- TailwindCSS styling
- Home, Pet, Cases, Marketplace pages

---

## 🛠️ Quick Start (5 хвилин)

### 1. Clone & Install

```bash
cd /Users/durkil/brainrot\ tamagotchi/unified

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

### 2. Setup Environment

**Contracts:**
```bash
cd contracts
cp .env.example .env
# Edit .env - додай PRIVATE_KEY і RPC URLs
```

**Backend:**
```bash
cd backend
cp .env.example .env
# Edit .env - налаштуй DATABASE_URL, REDIS_URL
```

**Frontend:**
```bash
cd frontend
cp .env.example .env.local
# Edit .env.local - встанови API_URL та WalletConnect ID
```

### 3. Start Services

**Option A: Docker (найпростіше)**
```bash
# В unified/ root
docker-compose up -d
```

**Option B: Manual**

Terminal 1 - PostgreSQL:
```bash
# Встанови PostgreSQL локально або використай Docker:
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres
```

Terminal 2 - Redis:
```bash
docker run -d -p 6379:6379 redis
```

Terminal 3 - Backend:
```bash
cd backend
go run cmd/main.go
```

Terminal 4 - Frontend:
```bash
cd frontend
npm run dev
```

### 4. Deploy Contracts (Base Sepolia)

```bash
cd contracts
npx hardhat run scripts/deploy.js --network base-sepolia
```

Збережи адреси контрактів і додай їх в .env файли!

---

## 🎯 Next Steps

### Що потрібно зробити далі:

#### 1. **Завершити blockchain інтеграцію**

Backend має заглушки для blockchain calls. Потрібно:

- Згенерувати Go bindings з контрактів:
```bash
cd contracts
npx hardhat compile
# Згенерувати ABI
cd ../backend
abigen --abi=../contracts/artifacts/contracts/BrainrotNFT.sol/BrainrotNFT.json --pkg=contracts --out=contracts/BrainrotNFT.go
```

- Імплементувати функції в `blockchain/client.go`:
  - MintNFT()
  - BurnNFT()
  - ListOnMarketplace()
  - BuyFromMarketplace()

- Додати event listeners для синхронізації з blockchain

#### 2. **IPFS для NFT metadata**

- Setup Pinata або власний IPFS node
- Створити JSON metadata для кожного мема
- Upload зображення мемів
- Згенерувати metadata URIs

#### 3. **Frontend покращення**

- Додати case opening анімацію
- Додати loading states
- Покращити error handling
- Додати toast notifications
- Inventory page (зараз placeholder)

#### 4. **Testing**

- Unit tests для смарт-контрактів (hardhat)
- Integration tests для backend (Go tests)
- E2E tests для frontend (Playwright)

#### 5. **Production готовність**

- Error tracking (Sentry)
- Analytics (Mixpanel)
- Monitoring (Grafana)
- Rate limiting
- Security audit

---

## 📁 Project Structure

```
unified/
├── contracts/          # Solidity smart contracts
│   ├── src/
│   │   ├── BrainrotNFT.sol
│   │   ├── CaseOpening.sol
│   │   ├── Marketplace.sol
│   │   └── BurnUpgrade.sol
│   ├── scripts/deploy.js
│   └── test/
│
├── backend/           # Golang API
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── api/          # HTTP handlers
│   │   ├── blockchain/   # Web3 integration
│   │   ├── models/       # Database models
│   │   ├── services/     # Business logic
│   │   └── repository/   # Data access
│   └── pkg/
│
├── frontend/          # Next.js app
│   ├── app/
│   │   ├── page.tsx      # Home
│   │   ├── pet/          # Pet page
│   │   ├── cases/        # Cases page
│   │   └── marketplace/  # Marketplace
│   ├── components/
│   └── lib/
│
├── docs/              # Documentation
├── docker-compose.yml
└── README.md
```

---

## 🧪 Development Workflow

### Розробка контрактів:
```bash
cd contracts
npx hardhat compile
npx hardhat test
npx hardhat node  # Local blockchain
```

### Розробка backend:
```bash
cd backend
go run cmd/main.go
# API на http://localhost:8080
```

### Розробка frontend:
```bash
cd frontend
npm run dev
# App на http://localhost:3000
```

---

## 🔥 Hot Reload

- Frontend: Next.js автоматично hot reload
- Backend: Використай `air` для Go hot reload:
```bash
go install github.com/cosmtrek/air@latest
cd backend
air
```

---

## 📚 Useful Commands

### Contracts:
```bash
npx hardhat compile                    # Compile contracts
npx hardhat test                       # Run tests
npx hardhat run scripts/deploy.js      # Deploy locally
npx hardhat verify --network base-sepolia ADDRESS  # Verify
```

### Backend:
```bash
go run cmd/main.go                     # Run
go test ./...                          # Test
go build -o brainrot cmd/main.go       # Build
```

### Frontend:
```bash
npm run dev                            # Development
npm run build                          # Build for production
npm run start                          # Start production build
npm run lint                           # Lint
```

---

## 🐛 Debugging

### Backend logs:
Backend виводить детальні logs. Шукай помилки в console.

### Frontend logs:
Відкрий browser DevTools (F12) → Console

### Database:
```bash
docker exec -it brainrot-postgres psql -U postgres -d brainrot
# SQL queries here
```

### Redis:
```bash
docker exec -it brainrot-redis redis-cli
# Redis commands here
```

---

## 💡 Tips

1. **Використовуй Base Sepolia** для testing (безкоштовні ETH з faucet)
2. **Встанови MetaMask** для локальної розробки
3. **Використовуй Hardhat console** для testing контрактів
4. **Перевіряй API через Postman** перед frontend інтеграцією
5. **Commit often** - робота велика, не втрачай прогрес!

---

## 🆘 Need Help?

- Smart Contracts: `/Users/durkil/brainrot tamagotchi/unified/contracts/`
- Backend: `/Users/durkil/brainrot tamagotchi/unified/backend/`
- Frontend: `/Users/durkil/brainrot tamagotchi/unified/frontend/`
- Documentation: `/Users/durkil/brainrot tamagotchi/unified/docs/`

---

## ✨ MVP Ready!

Базова структура готова! Тепер можеш:
1. Deploy контракти
2. Start backend
3. Start frontend
4. Почати розробку фіч

**Let's build the most brainrot tamagotchi ever! 🧠🎮**

