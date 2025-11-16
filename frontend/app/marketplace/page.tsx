'use client';

import { ConnectButton } from '@rainbow-me/rainbowkit';
import Link from 'next/link';
import { useState, useEffect } from 'react';
import { useAccount } from 'wagmi';
import { marketplaceAPI, setWalletAddress } from '@/lib/api';
import { motion } from 'framer-motion';

interface NFTListing {
  token_id: number;
  seller_address: string;
  price: number;
  nft?: {
    meme_type: string;
    rarity: string;
    level: number;
  };
}

export default function MarketplacePage() {
  const { address, isConnected } = useAccount();
  const [listings, setListings] = useState<NFTListing[]>([]);
  const [loading, setLoading] = useState(false);
  const [filter, setFilter] = useState<string>('all');

  useEffect(() => {
    if (address) {
      setWalletAddress(address);
      loadListings();
    }
  }, [address, filter]);

  const loadListings = async () => {
    setLoading(true);
    try {
      const params = filter !== 'all' ? { rarity: filter } : {};
      const response = await marketplaceAPI.getListings(params);
      setListings(response.data.listings || []);
    } catch (error) {
      console.error('Failed to load listings:', error);
    }
    setLoading(false);
  };

  const handleBuy = async (tokenId: number) => {
    try {
      await marketplaceAPI.buyNFT(tokenId);
      loadListings();
    } catch (error) {
      console.error('Failed to buy NFT:', error);
    }
  };

  const getMemeEmoji = (memeType: string) => {
    const emojis: Record<string, string> = {
      pepe: '🐸',
      doge: '🐕',
      gigachad: '💪',
      wojak: '😢',
      cheems: '🐕',
      drake: '🎤',
      vibing_cat: '😺',
      pikachu: '⚡',
    };
    return emojis[memeType] || '🎭';
  };

  const getRarityColor = (rarity: string) => {
    const colors: Record<string, string> = {
      common: 'text-gray-400',
      rare: 'text-blue-400',
      epic: 'text-purple-400',
      legendary: 'text-yellow-400',
    };
    return colors[rarity] || 'text-gray-400';
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
        <div className="text-center mb-8">
          <h2 className="text-5xl font-bold mb-4 bg-gradient-to-r from-cyan-400 to-pink-400 bg-clip-text text-transparent">
            🛒 Marketplace
          </h2>
          <p className="text-xl text-gray-300">
            Buy and sell meme NFTs
          </p>
        </div>

        {!isConnected ? (
          <div className="text-center py-20">
            <h3 className="text-2xl font-bold text-gray-300 mb-4">
              Connect your wallet to browse marketplace
            </h3>
            <ConnectButton />
          </div>
        ) : (
          <>
            {/* Filters */}
            <div className="flex justify-center gap-4 mb-8">
              {['all', 'common', 'rare', 'epic', 'legendary'].map((rarity) => (
                <button
                  key={rarity}
                  onClick={() => setFilter(rarity)}
                  className={`px-6 py-2 rounded-xl capitalize transition ${
                    filter === rarity
                      ? 'bg-purple-600'
                      : 'bg-gray-700 hover:bg-gray-600'
                  }`}
                >
                  {rarity}
                </button>
              ))}
            </div>

            {/* Listings Grid */}
            {loading ? (
              <div className="text-center py-20">
                <div className="text-4xl mb-4 animate-spin">🔄</div>
                <p className="text-gray-400">Loading listings...</p>
              </div>
            ) : listings.length === 0 ? (
              <div className="text-center py-20">
                <p className="text-2xl text-gray-400">No listings found</p>
              </div>
            ) : (
              <div className="grid md:grid-cols-3 lg:grid-cols-4 gap-6">
                {listings.map((listing, index) => (
                  <motion.div
                    key={listing.token_id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.05 }}
                    className="brainrot-card"
                  >
                    <div className="text-6xl mb-4 text-center">
                      {listing.nft && getMemeEmoji(listing.nft.meme_type)}
                    </div>
                    
                    <div className="mb-4">
                      <h3 className="text-xl font-bold capitalize mb-2">
                        {listing.nft?.meme_type} #{listing.token_id}
                      </h3>
                      <div className="flex gap-2 text-sm">
                        <span className={`capitalize ${getRarityColor(listing.nft?.rarity || 'common')}`}>
                          {listing.nft?.rarity}
                        </span>
                        <span className="text-gray-400">
                          • Lv {listing.nft?.level}
                        </span>
                      </div>
                    </div>

                    <div className="border-t border-gray-700 pt-4">
                      <div className="text-2xl font-bold text-purple-400 mb-4">
                        {listing.price} BASE
                      </div>
                      <button
                        onClick={() => handleBuy(listing.token_id)}
                        className="brainrot-button w-full text-sm"
                      >
                        Buy Now
                      </button>
                    </div>
                  </motion.div>
                ))}
              </div>
            )}
          </>
        )}
      </div>
    </main>
  );
}

