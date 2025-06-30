import { Select } from 'antd';

export function NominatePanel({ session, onNominate, disabled }: { session: any, onNominate: (id: string) => void, disabled: boolean }) {
  return (
    <Select
      placeholder="Nominate user"
      style={{ width: 200, marginLeft: 8 }}
      onChange={onNominate}
      disabled={disabled}
    >
      {session.participants.map((p: any) => (
        <Select.Option key={p.id} value={p.id}>{p.name}</Select.Option>
      ))}
    </Select>
  );
} 