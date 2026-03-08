import axios, { type AxiosInstance } from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

// Create axios instance for agent API calls
export function createAgentClient(apiKey: string): AxiosInstance {
  const client = axios.create({
    baseURL: API_BASE_URL,
    headers: {
      'Content-Type': 'application/json',
      'X-API-KEY': apiKey,
    },
  });

  return client;
}

// Or get a singleton instance with stored API key
export function getAgentClient(): AxiosInstance | null {
  const apiKey = localStorage.getItem('agent_api_key');
  if (!apiKey) return null;

  return axios.create({
    baseURL: API_BASE_URL,
    headers: {
      'Content-Type': 'application/json',
      'X-API-KEY': apiKey,
    },
  });
}

// Store API key for agent mode
export function setAgentApiKey(apiKey: string) {
  localStorage.setItem('agent_api_key', apiKey);
}

// Clear API key
export function clearAgentApiKey() {
  localStorage.removeItem('agent_api_key');
}

// Check if in agent mode
export function isAgentMode(): boolean {
  return !!localStorage.getItem('agent_api_key');
}
