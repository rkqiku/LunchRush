import { Form, Input, Button, Modal } from 'antd';

export function OrderForm({ open, onCancel, onOrder }: { open: boolean, onCancel: () => void, onOrder: (values: any) => void }) {
  return (
    <Modal title="Place/Update Order" open={open} onCancel={onCancel} footer={null}>
      <Form onFinish={onOrder} layout="vertical">
        <Form.Item name="restaurant" label="Restaurant" rules={[{ required: true }]}> <Input /> </Form.Item>
        <Form.Item name="dish" label="Dish" rules={[{ required: true }]}> <Input /> </Form.Item>
        <Form.Item name="notes" label="Notes"> <Input /> </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit">Submit</Button>
        </Form.Item>
      </Form>
    </Modal>
  );
} 