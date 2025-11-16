'use client';

import { ConnectButton } from '@rainbow-me/rainbowkit';
import Link from 'next/link';
import { useState, useEffect } from 'react';
import { useAccount } from 'wagmi';
import { petAPI, setWalletAddress } from '@/lib/api';
import { motion } from 'framer-motion';

interface PetState {
  token_id: number;
  meme_type: string;
  rarity: string;
  level: number;
  hunger: number;
  mood: number;
  energy: number;
}

export default function PetPage() {
  const { address, isConnected } = useAccount();
  const [pet, setPet] = useState<PetState | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (address) {
      setWalletAddress(address);
      // TODO: Load user's primary pet
    }
  }, [address]);

  const handleFeed = async () => {
    if (!pet) return;
    setLoading(true);
    try {
      await petAPI.feedPet(pet.token_id);
      // Refresh pet state
      const response = await petAPI.getPet(pet.token_id);
      setPet(response.data);
    } catch (error) {
      console.error('Failed to feed pet:', error);
    }
    setLoading(false);
  };

  const handlePlay = async () => {
    if (!pet) return;
    setLoading(true);
    try {
      await petAPI.playWithPet(pet.token_id);
      // Refresh pet state
      const response = await petAPI.getPet(pet.token_id);
      setPet(response.data);
    } catch (error) {
      console.error('Failed to play with pet:', error);
    }
    setLoading(false);
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
        {!isConnected ? (
          <div className="text-center py-20">
            <h2 className="text-3xl font-bold text-gray-300 mb-4">
              Connect your wallet to see your pet!
            </h2>
            <ConnectButton />
          </div>
        ) : !pet ? (
          <div className="text-center py-20">
            <h2 className="text-3xl font-bold text-gray-300 mb-4">
              You don't have a pet yet!
            </h2>
            <Link href="/cases">
              <button className="brainrot-button">
                🎁 Open a Case to Get One
              </button>
            </Link>
          </div>
        ) : (
          <div className="max-w-2xl mx-auto">
            {/* Pet Display */}
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              className="brainrot-card text-center mb-8"
            >
              <div className="text-8xl mb-4 animate-bounce-slow">
                {pet.meme_type === 'pepe' && '🐸'}
                {pet.meme_type === 'doge' && '🐕'}
                {pet.meme_type === 'gigachad' && '💪'}
                {pet.meme_type === 'wojak' && '😢'}
                {pet.meme_type === 'cheems' && '🐕'}
                {pet.meme_type === 'drake' && '🎤'}
                {pet.meme_type === 'vibing_cat' && '😺'}
                {pet.meme_type === 'pikachu' && '⚡'}
              </div>
              <h2 className="text-3xl font-bold capitalize mb-2">
                {pet.meme_type} #{pet.token_id}
              </h2>
              <div className="flex justify-center gap-4 text-sm">
                <span className="px-3 py-1 bg-purple-500/30 rounded-full capitalize">
                  {pet.rarity}
                </span>
                <span className="px-3 py-1 bg-pink-500/30 rounded-full">
                  Level {pet.level}
                </span>
              </div>
            </motion.div>

            {/* Stats */}
            <div className="brainrot-card mb-8">
              <h3 className="text-2xl font-bold mb-6">Pet Stats</h3>
              
              {/* Hunger */}
              <div className="mb-4">
                <div className="flex justify-between mb-2">
                  <span>🍔 Hunger</span>
                  <span>{pet.hunger}%</span>
                </div>
                <div className="stat-bar">
                  <div
                    className="stat-bar-fill from-green-500 to-green-600"
                    style={{ width: `${pet.hunger}%` }}
                  />
                </div>
              </div>

              {/* Mood */}
              <div className="mb-4">
                <div className="flex justify-between mb-2">
                  <span>😊 Mood</span>
                  <span>{pet.mood}%</span>
                </div>
                <div className="stat-bar">
                  <div
                    className="stat-bar-fill from-yellow-500 to-yellow-600"
                    style={{ width: `${pet.mood}%` }}
                  />
                </div>
              </div>

              {/* Energy */}
              <div className="mb-4">
                <div className="flex justify-between mb-2">
                  <span>⚡ Energy</span>
                  <span>{pet.energy}%</span>
                </div>
                <div className="stat-bar">
                  <div
                    className="stat-bar-fill from-blue-500 to-blue-600"
                    style={{ width: `${pet.energy}%` }}
                  />
                </div>
              </div>
            </div>

            {/* Actions */}
            <div className="grid grid-cols-2 gap-4">
              <button
                onClick={handleFeed}
                disabled={loading || pet.hunger >= 100}
                className="brainrot-button disabled:opacity-50 disabled:cursor-not-allowed"
              >
                🍔 Feed
              </button>
              <button
                onClick={handlePlay}
                disabled={loading || pet.energy < 10}
                className="brainrot-button disabled:opacity-50 disabled:cursor-not-allowed"
              >
                🎮 Play
              </button>
            </div>

            {/* Navigation */}
            <div className="mt-8 flex justify-center gap-4">
              <Link href="/inventory">
                <button className="px-6 py-3 bg-gray-700 hover:bg-gray-600 rounded-xl transition">
                  📦 Inventory
                </button>
              </Link>
              <Link href="/marketplace">
                <button className="px-6 py-3 bg-gray-700 hover:bg-gray-600 rounded-xl transition">
                  🛒 Marketplace
                </button>
              </Link>
            </div>
          </div>
        )}
      </div>
    </main>
  );
}

