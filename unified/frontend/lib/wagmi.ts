import { getDefaultConfig } from '@rainbow-me/rainbowkit';
import { base, baseSepolia } from 'wagmi/chains';

export const config = getDefaultConfig({
  appName: 'Brainrot Tamagotchi',
  projectId: process.env.NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID || 'YOUR_PROJECT_ID',
  chains: [
    ...(process.env.NEXT_PUBLIC_ENABLE_TESTNETS === 'true' ? [baseSepolia] : []),
    base,
  ],
  ssr: true,
});

