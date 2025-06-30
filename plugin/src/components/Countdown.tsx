import { Statistic, message } from 'antd';

export function Countdown({ lockAt }: { lockAt: string }) {
  const deadline = new Date(lockAt).getTime();
  return (
    <Statistic.Countdown
      title="Auto-lock in"
      value={deadline}
      onFinish={() => message.info('Session auto-locked!')}
      style={{ marginBottom: 24 }}
    />
  );
} 