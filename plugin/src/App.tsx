import React, { useEffect, useState, useCallback } from 'react';
import { Layout, Typography, Button, message } from 'antd';
import { getTodaySession, getSession, createSession, joinSession, placeOrder, vote, nominate, lockSession } from './api';
import { useLunchRushSocket } from './hooks/useLunchRushSocket';
import { SessionView } from './components/SessionView';
import { JoinForm } from './components/JoinForm';
import { OrderForm } from './components/OrderForm';
import { VotePanel } from './components/VotePanel';
import { NominatePanel } from './components/NominatePanel';
import { LockButton } from './components/LockButton';
import { SummaryView } from './components/SummaryView';
import { Countdown } from './components/Countdown';
import { v4 as uuidv4 } from 'uuid';

const { Header, Content, Footer } = Layout;
const { Title } = Typography;

function useQuery() {
  return new URLSearchParams(window.location.search);
}

function App() {
  const [session, setSession] = useState<any>(null);
  const [user, setUser] = useState<any>(null);
  const [joinModal, setJoinModal] = useState(false);
  const [orderModal, setOrderModal] = useState(false);
  const query = useQuery();
  const sessionIdFromUrl = query.get('session') || undefined;

  // On mount, generate or load userId
  useEffect(() => {
    let userId = localStorage.getItem('lunchrush_userId');
    if (!userId) {
      userId = uuidv4();
      localStorage.setItem('lunchrush_userId', userId);
    }
    setUser((u: any) => u ? u : { id: userId });
  }, []);

  // Fetch session on mount or when sessionIdFromUrl changes
  useEffect(() => {
    if (sessionIdFromUrl) {
      getSession(sessionIdFromUrl).then(setSession);
    } else {
      getTodaySession().then(setSession);
    }
  }, [sessionIdFromUrl]);

  // WebSocket for real-time updates
  useLunchRushSocket(
    useCallback((evt) => {
      if (evt.sessionId === session?.id) {
        getTodaySession().then(setSession);
      }
    }, [session?.id])
  );

  // Join session
  const handleJoin = (values: any) => {
    joinSession(session.id, user.id, values.name).then((s) => {
      setUser({ ...user, name: values.name });
      setSession(s);
      setJoinModal(false);
      message.success('Joined session!');
    });
  };

  // Place order
  const handleOrder = (values: any) => {
    placeOrder(session.id, user.id, values.restaurant, values.dish, values.notes).then((s) => {
      setSession(s);
      setOrderModal(false);
      message.success('Order placed!');
    });
  };

  // Vote
  const handleVote = (restaurant: string) => {
    vote(session.id, restaurant).then((s) => {
      setSession(s);
      message.success('Voted!');
    });
  };

  // Nominate
  const handleNominate = (userId: string) => {
    nominate(session.id, userId).then((s) => {
      setSession(s);
      message.success('User nominated!');
    });
  };

  // Lock session
  const handleLock = () => {
    lockSession(session.id, user.id).then((s) => {
      setSession(s);
      message.success('Session locked!');
    });
  };

  // Helper: is user joined?
  const isJoined = user && session?.participants.some((p: any) => p.id === user.id);

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header>
        <Title style={{ color: 'white', margin: 0 }} level={2}>LunchRush</Title>
      </Header>
      <Content style={{ padding: '2rem' }}>
        <Title level={3}>Today's Lunch Session</Title>
        {session && (
          <div style={{ marginBottom: 16 }}>
            <b>Share this session:</b> <input style={{ width: 350 }} value={`${window.location.origin}?session=${session.id}`} readOnly onFocus={e => e.target.select()} />
          </div>
        )}
        {!session ? (
          <div>
            <Button type="primary" onClick={() => {
              const today = new Date().toISOString().slice(0, 10);
              createSession(today).then(setSession);
            }}>Create Session</Button>
          </div>
        ) : (
          <>
            {session.lockAt && !session.locked && <Countdown lockAt={session.lockAt} />}
            <SessionView session={session} />
            {!isJoined ? (
              <Button type="primary" onClick={() => setJoinModal(true)} style={{ marginTop: 16 }}>Join Session</Button>
            ) : (
              <>
                {!session.locked && (
                  <>
                    <Button onClick={() => setOrderModal(true)} style={{ marginTop: 16 }}>Place/Update Order</Button>
                    <LockButton onLock={handleLock} disabled={session.locked} />
                    <div style={{ marginTop: 16 }}>
                      <VotePanel session={session} onVote={handleVote} disabled={session.locked} />
                      <NominatePanel session={session} onNominate={handleNominate} disabled={session.locked} />
                    </div>
                  </>
                )}
                {session.locked && <SummaryView session={session} />}
              </>
            )}
            <JoinForm open={joinModal} onCancel={() => setJoinModal(false)} onJoin={handleJoin} requireUserId={false} />
            <OrderForm open={orderModal} onCancel={() => setOrderModal(false)} onOrder={handleOrder} />
          </>
        )}
      </Content>
      <Footer style={{ textAlign: 'center' }}>LunchRush ©2025</Footer>
    </Layout>
  );
}

export default App;
