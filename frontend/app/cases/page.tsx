'use client';

import { ConnectButton } from '@rainbow-me/rainbowkit';
import Link from 'next/link';
import { motion } from 'framer-motion';
import { useAccount } from 'wagmi';

const cases = [
  {
    id: 'bronze',
    name: 'Bronze Case',
    price: '$0.50',
    emoji: '🥉',
    gradient: 'from-amber-700 to-amber-900',
    description: '80% Common, 20% Rare',
  },
  {
    id: 'silver',
    name: 'Silver Case',
    price: '$2.00',
    emoji: '🥈',
    gradient: 'from-gray-400 to-gray-600',
    description: '70% Rare, 25% Epic, 5% Legendary',
  },
  {
    id: 'gold',
    name: 'Gold Case',
    price: '$10.00',
    emoji: '🥇',
    gradient: 'from-yellow-400 to-yellow-600',
    description: '60% Epic, 40% Legendary',
  },
];

export default function CasesPage() {
  const { isConnected } = useAccount();

  const handleOpenCase = async (caseType: string) => {
    // TODO: Implement case opening
    console.log(`Opening ${caseType} case...`);
  };

  return (
    <main className="min-h-screen bg-gradient-to-b from-gray-900 via-purple-900 to-gray-900">
      {/* Header */}
      <nav className="container mx-auto px-4 py-6">
        <div className="flex justify-between items-center">
          <Link href="/">
            <h1 className="text-3xl font-bold bg-gradient-to-r from-purple-400 to-pink-400 bg-clip-text text-transparent cursor-pointer">
              🧠 Brainrot Tamagotchi
            </h1>
          </Link>
          <ConnectButton />
        </div>
      </nav>

      <div className="container mx-auto px-4 py-8">
        <div className="text-center mb-12">
          <h2 className="text-5xl font-bold mb-4 bg-gradient-to-r from-cyan-400 to-pink-400 bg-clip-text text-transparent">
            🎁 Open Cases
          </h2>
          <p className="text-xl text-gray-300">
            Get rare meme NFTs by opening cases!
          </p>
        </div>

        {!isConnected ? (
          <div className="text-center py-20">
            <h3 className="text-2xl font-bold text-gray-300 mb-4">
              Connect your wallet to open cases
            </h3>
            <ConnectButton />
          </div>
        ) : (
          <div className="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
            {cases.map((caseItem, index) => (
              <motion.div
                key={caseItem.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
                whileHover={{ scale: 1.05 }}
                className="brainrot-card text-center"
              >
                <div
                  className={`text-8xl mb-4 bg-gradient-to-br ${caseItem.gradient} rounded-full w-32 h-32 flex items-center justify-center mx-auto`}
                >
                  {caseItem.emoji}
                </div>
                
                <h3 className="text-2xl font-bold mb-2">{caseItem.name}</h3>
                <p className="text-3xl font-bold text-purple-400 mb-4">
                  {caseItem.price}
                </p>
                <p className="text-gray-400 mb-6">{caseItem.description}</p>
                
                <button
                  onClick={() => handleOpenCase(caseItem.id)}
                  className="brainrot-button w-full"
                >
                  Open Case
                </button>
              </motion.div>
            ))}
          </div>
        )}

        {/* Info Section */}
        <div className="max-w-3xl mx-auto mt-16 brainrot-card">
          <h3 className="text-2xl font-bold mb-4">How it works</h3>
          <ol className="list-decimal list-inside space-y-2 text-gray-300">
            <li>Choose a case type based on your budget</li>
            <li>Buy the case with Base (ETH)</li>
            <li>Watch the opening animation</li>
            <li>Receive your random meme NFT!</li>
            <li>Take care of your new pet or trade it on marketplace</li>
          </ol>
        </div>
      </div>
    </main>
  );
}

