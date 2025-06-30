import { Card, List, Avatar } from 'antd';

export function SessionView({ session }: { session: any }) {
  return (
    <Card title="Participants" style={{ marginBottom: 24 }}>
      <List
        itemLayout="horizontal"
        dataSource={session.participants}
        renderItem={(p: any) => (
          <List.Item>
            <List.Item.Meta
              avatar={<Avatar>{p.name[0]}</Avatar>}
              title={p.name}
              description={p.order ? `${p.order.dish} @ ${p.order.restaurant}` : 'No order yet'}
            />
          </List.Item>
        )}
      />
    </Card>
  );
} 