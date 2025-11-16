'use client';

import { ConnectButton } from '@rainbow-me/rainbowkit';
import Link from 'next/link';
import { motion } from 'framer-motion';

export default function Home() {
  return (
    <main className="min-h-screen bg-gradient-to-b from-gray-900 via-purple-900 to-gray-900">
      {/* Header */}
      <nav className="container mx-auto px-4 py-6">
        <div className="flex justify-between items-center">
          <h1 className="text-3xl font-bold bg-gradient-to-r from-purple-400 to-pink-400 bg-clip-text text-transparent">
            🧠 Brainrot Tamagotchi
          </h1>
          <ConnectButton />
        </div>
      </nav>

      {/* Hero Section */}
      <div className="container mx-auto px-4 py-20">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          className="text-center max-w-4xl mx-auto"
        >
          <h2 className="text-6xl font-extrabold mb-6 bg-gradient-to-r from-cyan-400 via-purple-400 to-pink-400 bg-clip-text text-transparent animate-gradient">
            The Most Brainrot Tamagotchi Ever 🎮
          </h2>
          
          <p className="text-xl text-gray-300 mb-8">
            Collect meme NFTs, take care of your brainrot pet, and trade on the marketplace.
            <br />
            Built on Base blockchain 🔵
          </p>

          <div className="flex justify-center gap-4 mb-12">
            <Link href="/pet">
              <button className="brainrot-button">
                🐸 My Pet
              </button>
            </Link>
            <Link href="/cases">
              <button className="brainrot-button">
                🎁 Open Cases
              </button>
            </Link>
            <Link href="/marketplace">
              <button className="brainrot-button">
                🛒 Marketplace
              </button>
            </Link>
          </div>

          {/* Feature Cards */}
          <div className="grid md:grid-cols-3 gap-6 mt-16">
            <motion.div
              whileHover={{ scale: 1.05 }}
              className="brainrot-card text-left"
            >
              <div className="text-4xl mb-4">🎨</div>
              <h3 className="text-2xl font-bold mb-2 text-purple-300">Meme NFTs</h3>
              <p className="text-gray-400">
                Collect rare meme characters from Pepe to Gigachad. Each NFT is unique!
              </p>
            </motion.div>

            <motion.div
              whileHover={{ scale: 1.05 }}
              className="brainrot-card text-left"
            >
              <div className="text-4xl mb-4">🎮</div>
              <h3 className="text-2xl font-bold mb-2 text-pink-300">Tamagotchi Care</h3>
              <p className="text-gray-400">
                Feed, play, and level up your pet. Don't let them die!
              </p>
            </motion.div>

            <motion.div
              whileHover={{ scale: 1.05 }}
              className="brainrot-card text-left"
            >
              <div className="text-4xl mb-4">💎</div>
              <h3 className="text-2xl font-bold mb-2 text-cyan-300">Trade & Earn</h3>
              <p className="text-gray-400">
                Buy, sell, and trade NFTs on the marketplace. Build your collection!
              </p>
            </motion.div>
          </div>

          {/* Stats */}
          <div className="mt-16 grid grid-cols-3 gap-8 max-w-2xl mx-auto">
            <div>
              <div className="text-4xl font-bold text-purple-400">8</div>
              <div className="text-gray-400">Meme Types</div>
            </div>
            <div>
              <div className="text-4xl font-bold text-pink-400">35+</div>
              <div className="text-gray-400">Unique NFTs</div>
            </div>
            <div>
              <div className="text-4xl font-bold text-cyan-400">Base</div>
              <div className="text-gray-400">Blockchain</div>
            </div>
          </div>
        </motion.div>
      </div>

      {/* Footer */}
      <footer className="container mx-auto px-4 py-8 text-center text-gray-500">
        <p>Built with 💙 for Base ecosystem</p>
        <p className="mt-2">🧠 Stay brainrot, frens 🎮</p>
      </footer>
    </main>
  );
}

