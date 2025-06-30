const API_URL = 'http://localhost:8080';

export async function getTodaySession() {
  const res = await fetch(`${API_URL}/session/today`);
  if (!res.ok) return null;
  return res.json();
}

export async function getSession(id: string) {
  const res = await fetch(`${API_URL}/session/${id}`);
  if (!res.ok) return null;
  return res.json();
}

export async function createSession(date: string, lockAt?: string) {
  return fetch(`${API_URL}/session`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ date, lockAt }),
  }).then(res => res.json());
}

export async function joinSession(sessionId: string, userId: string, name: string) {
  return fetch(`${API_URL}/session/${sessionId}/join`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ userId, name }),
  }).then(res => res.json());
}

export async function placeOrder(sessionId: string, userId: string, restaurant: string, dish: string, notes?: string) {
  return fetch(`${API_URL}/session/${sessionId}/order`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ userId, restaurant, dish, notes }),
  }).then(res => res.json());
}

export async function vote(sessionId: string, restaurant: string) {
  return fetch(`${API_URL}/session/${sessionId}/vote`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ restaurant }),
  }).then(res => res.json());
}

export async function nominate(sessionId: string, userId: string) {
  return fetch(`${API_URL}/session/${sessionId}/nominate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ userId }),
  }).then(res => res.json());
}

export async function lockSession(sessionId: string, userId: string) {
  return fetch(`${API_URL}/session/${sessionId}/lock`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ userId }),
  }).then(res => res.json());
} 