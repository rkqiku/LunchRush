import { Result } from 'antd';

export function SummaryView({ session }: { session: any }) {
  return (
    <div style={{ marginTop: 24 }}>
      <Result
        status="success"
        title="Lunch Order Locked!"
        subTitle={`Nominated: ${
          session.participants.find((p: any) => p.id === session.nominatedUserId)?.name || 'N/A'
        }`}
      />
    </div>
  );
} 