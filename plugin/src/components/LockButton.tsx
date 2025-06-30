import { Button } from 'antd';

export function LockButton({ onLock, disabled }: { onLock: () => void, disabled: boolean }) {
  return (
    <Button onClick={onLock} style={{ marginLeft: 8, marginTop: 16 }} disabled={disabled}>
      Lock Session
    </Button>
  );
} 