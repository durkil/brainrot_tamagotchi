import { getDefaultConfig } from '@rainbow-me/rainbowkit';
import { base, baseSepolia } from 'wagmi/chains';

const chains = process.env.NEXT_PUBLIC_ENABLE_TESTNETS === 'true' 
  ? [baseSepolia, base] as const
  : [base] as const;

export const config = getDefaultConfig({
  appName: 'Brainrot Tamagotchi',
  projectId: process.env.NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID || 'YOUR_PROJECT_ID',
  chains,
  ssr: true,
});

