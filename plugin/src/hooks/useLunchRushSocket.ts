import { useEffect } from 'react';

export function useLunchRushSocket(onEvent: (evt: any) => void) {
  useEffect(() => {
    const ws = new WebSocket('ws://localhost:8080/ws');
    ws.onmessage = (event) => {
      try {
        const evt = JSON.parse(event.data);
        onEvent(evt);
      } catch {}
    };
    return () => ws.close();
  }, [onEvent]);
} 