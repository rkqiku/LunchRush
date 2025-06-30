import React, { useEffect } from 'react';
import { Form, Input, Button, Modal } from 'antd';

export function JoinForm({ open, onCancel, onJoin, requireUserId = true }: { open: boolean, onCancel: () => void, onJoin: (values: any) => void, requireUserId?: boolean }) {
  const [form] = Form.useForm();

  useEffect(() => {
    if (open) {
      form.resetFields();
    }
  }, [open, form]);

  return (
    <Modal title="Join Session" open={open} onCancel={onCancel} footer={null} forceRender>
      <Form form={form} onFinish={onJoin} layout="vertical" autoComplete="off">
        {requireUserId && (
          <Form.Item name="userId" label="User ID" rules={[{ required: true }]} key="userId"> <Input autoComplete="off" /> </Form.Item>
        )}
        <Form.Item name="name" label="Name" rules={[{ required: true, message: 'Please enter your name' }]} key="name">
          <Input autoComplete="off" />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit">Join</Button>
        </Form.Item>
      </Form>
    </Modal>
  );
} 