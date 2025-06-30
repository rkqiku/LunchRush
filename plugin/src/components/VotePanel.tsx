import { Select } from 'antd';

export function VotePanel({ session, onVote, disabled }: { session: any, onVote: (r: string) => void, disabled: boolean }) {
  return (
    <Select
      placeholder="Vote for a restaurant"
      style={{ width: 200 }}
      onChange={onVote}
      disabled={disabled}
    >
      {session.participants.map((p: any) => p.order?.restaurant).filter(Boolean).map((r: string, i: number) => (
        <Select.Option key={i} value={r}>{r}</Select.Option>
      ))}
    </Select>
  );
} 