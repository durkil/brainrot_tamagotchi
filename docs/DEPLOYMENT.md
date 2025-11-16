# 🚀 Deployment Guide - Brainrot Tamagotchi

Цей документ описує процес деплою всього проекту.

---

## 📋 Prerequisites

Перед деплоєм, переконайся що маєш:

- Node.js 18+ встановлено
- Go 1.21+ встановлено
- PostgreSQL database (локально або managed)
- Redis server (локально або managed)
- Base testnet ETH (для деплою контрактів)
- Wallet з private key

---

## 1️⃣ Deploy Smart Contracts (Base Sepolia)

### Setup

```bash
cd contracts
npm install
cp .env.example .env
```

### Edit `.env`:

```env
PRIVATE_KEY=your_private_key_without_0x
BASE_SEPOLIA_RPC=https://sepolia.base.org
BASESCAN_API_KEY=your_basescan_key
```

### Deploy

```bash
npx hardhat compile
npx hardhat run scripts/deploy.js --network base-sepolia
```

Збережи адреси контрактів з output!

### Verify Contracts

```bash
npx hardhat verify --network base-sepolia CONTRACT_ADDRESS
```

---

## 2️⃣ Setup Backend (Golang API)

### Local Development

```bash
cd backend
cp .env.example .env
```

### Edit `.env`:

```env
PORT=8080
DATABASE_URL=postgresql://user:pass@localhost:5432/brainrot
REDIS_URL=redis://localhost:6379
BASE_RPC_URL=https://sepolia.base.org
PRIVATE_KEY=your_private_key
CONTRACT_NFT_ADDRESS=0x...      # From step 1
CONTRACT_CASE_ADDRESS=0x...     # From step 1
CONTRACT_MARKETPLACE_ADDRESS=0x... # From step 1
CONTRACT_BURN_ADDRESS=0x...     # From step 1
```

### Run

```bash
# Install dependencies
go mod download

# Run migrations (automatic on start)
go run cmd/main.go
```

API буде доступний на `http://localhost:8080`

### Production Deploy (Docker)

```bash
cd backend
docker build -t brainrot-backend .
docker run -p 8080:8080 --env-file .env brainrot-backend
```

Або використай docker-compose (див. нижче).

---

## 3️⃣ Setup Frontend (Next.js)

### Local Development

```bash
cd frontend
npm install
cp .env.example .env.local
```

### Edit `.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID=your_wc_project_id
NEXT_PUBLIC_ENABLE_TESTNETS=true
NEXT_PUBLIC_NFT_CONTRACT=0x...
NEXT_PUBLIC_CASE_CONTRACT=0x...
NEXT_PUBLIC_MARKETPLACE_CONTRACT=0x...
NEXT_PUBLIC_BURN_CONTRACT=0x...
```

### Run

```bash
npm run dev
```

Frontend буде на `http://localhost:3000`

### Production Deploy (Vercel)

```bash
# Install Vercel CLI
npm i -g vercel

# Deploy
cd frontend
vercel --prod
```

Або push до GitHub і підключи Vercel автоматично.

---

## 🐳 Docker Compose (All-in-One)

Запустити весь стек локально:

```bash
# В root директорії unified/
docker-compose up -d
```

Це запустить:
- PostgreSQL (port 5432)
- Redis (port 6379)
- Backend API (port 8080)
- Frontend (port 3000)

---

## 📝 Environment Variables Summary

### Contracts

| Variable | Description |
|----------|-------------|
| `PRIVATE_KEY` | Wallet private key (без 0x) |
| `BASE_SEPOLIA_RPC` | Base Sepolia RPC URL |
| `BASESCAN_API_KEY` | Basescan API key для верифікації |

### Backend

| Variable | Description |
|----------|-------------|
| `PORT` | API port (default: 8080) |
| `DATABASE_URL` | PostgreSQL connection string |
| `REDIS_URL` | Redis connection string |
| `BASE_RPC_URL` | Base RPC URL |
| `PRIVATE_KEY` | Wallet private key |
| `CONTRACT_*_ADDRESS` | Smart contract addresses |

### Frontend

| Variable | Description |
|----------|-------------|
| `NEXT_PUBLIC_API_URL` | Backend API URL |
| `NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID` | WalletConnect Project ID |
| `NEXT_PUBLIC_ENABLE_TESTNETS` | Enable testnets (true/false) |
| `NEXT_PUBLIC_*_CONTRACT` | Contract addresses |

---

## 🧪 Testing Deployment

### 1. Test Smart Contracts

```bash
cd contracts
npx hardhat test
```

### 2. Test Backend API

```bash
curl http://localhost:8080/api/v1/health
```

Should return: `{"status":"ok"}`

### 3. Test Frontend

Відкрий `http://localhost:3000` і:
- Connect wallet
- Check that все працює

---

## 🌐 Production Checklist

- [ ] Контракти задеплоєні на Base Mainnet
- [ ] Контракти верифіковані на Basescan
- [ ] Backend задеплоєний (VPS/Cloud/Docker)
- [ ] Database setup з backups
- [ ] Redis setup
- [ ] Frontend задеплоєний на Vercel
- [ ] Domain налаштований
- [ ] SSL certificates
- [ ] Environment variables встановлені
- [ ] Monitoring setup
- [ ] Error tracking (Sentry)
- [ ] Analytics (Google Analytics / Mixpanel)

---

## 🔧 Troubleshooting

### Backend не запускається

- Перевір DATABASE_URL
- Перевір Redis connection
- Перевір що всі env variables встановлені

### Frontend не підключається до Backend

- Перевір NEXT_PUBLIC_API_URL
- Перевір CORS settings в backend
- Відкрий browser console для errors

### Wallet не підключається

- Перевір NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID
- Перевір що використовуєш правильний chain (Base Sepolia)

### Smart contract calls fail

- Перевір що контракти задеплоєні
- Перевір адреси контрактів
- Перевір що wallet має Base ETH

---

## 📊 Monitoring

### Backend Health

```bash
curl http://your-api.com/api/v1/health
```

### Database

```sql
-- Check NFT count
SELECT COUNT(*) FROM nfts;

-- Check active listings
SELECT COUNT(*) FROM market_listings WHERE is_active = true;

-- Check case openings
SELECT COUNT(*) FROM case_openings;
```

---

## 🆘 Support

Якщо є проблеми з деплоєм:
1. Перевір logs
2. Перевір .env файли
3. Перевір network connectivity
4. Перегляни documentation

---

**Ready to deploy! 🚀**

