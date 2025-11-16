import axios from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add wallet address to requests
export const setWalletAddress = (address: string) => {
  api.defaults.headers.common['X-Wallet-Address'] = address;
};

// API functions
export const petAPI = {
  getPet: (tokenId: number) => api.get(`/pets/${tokenId}`),
  feedPet: (tokenId: number, isPaid: boolean = false) => 
    api.post(`/pets/${tokenId}/feed`, { is_paid: isPaid }),
  playWithPet: (tokenId: number) => api.post(`/pets/${tokenId}/play`),
};

export const casesAPI = {
  getPrices: () => api.get('/cases/prices'),
  buyCase: (caseType: string) => api.post('/cases/buy', { case_type: caseType }),
  openCase: (caseId: string) => api.post(`/cases/${caseId}/open`),
};

export const marketplaceAPI = {
  getListings: (params?: { rarity?: string; min_level?: number; limit?: number; offset?: number }) => 
    api.get('/marketplace', { params }),
  listNFT: (tokenId: number, price: number) => 
    api.post('/marketplace/list', { token_id: tokenId, price }),
  buyNFT: (tokenId: number) => api.post(`/marketplace/${tokenId}/buy`),
  cancelListing: (tokenId: number) => api.delete(`/marketplace/${tokenId}`),
};

export const userAPI = {
  getUser: (address: string) => api.get(`/users/${address}`),
  getInventory: (address: string) => api.get(`/users/${address}/inventory`),
};

