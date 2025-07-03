import { useEffect } from 'react';
import api from '../api/client';

const HEARTBEAT_INTERVAL = 30000; // 30 seconds

export const useHeartbeat = (sessionId, username, enabled = true) => {
  useEffect(() => {
    if (!enabled || !sessionId || !username) return;

    const sendHeartbeat = async () => {
      try {
        await api.sendHeartbeat(sessionId, username);
      } catch (error) {
        console.error('Heartbeat failed:', error);
      }
    };

    // Send initial heartbeat
    sendHeartbeat();

    // Set up interval for periodic heartbeats
    const intervalId = setInterval(sendHeartbeat, HEARTBEAT_INTERVAL);

    // Cleanup on unmount or dependency change
    return () => {
      clearInterval(intervalId);
    };
  }, [sessionId, username, enabled]);
};