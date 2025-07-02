import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Pizza, Clock, Users, ChefHat } from 'lucide-react';
import api from '../api/client';
import useUser from '../hooks/useUser';
import { useHeartbeat } from '../hooks/useHeartbeat';
import SessionHeader from './SessionHeader';
import JoinSession from './JoinSession';
import RestaurantSection from './RestaurantSection';
import ParticipantsList from './ParticipantsList';
import OrderSummary from './OrderSummary';
import Toast from './Toast';

const LunchRush = () => {
  const { user } = useUser();
  const queryClient = useQueryClient();
  const [toast, setToast] = useState(null);

  // Query for today's session
  const { data: session, isLoading, error } = useQuery({
    queryKey: ['session', 'today'],
    queryFn: async () => {
      return await api.getTodaySession();
    },
  });

  // Create session mutation
  const createSessionMutation = useMutation({
    mutationFn: api.createSession,
    onSuccess: () => {
      queryClient.invalidateQueries(['session']);
      showToast('Session created successfully!', 'success');
    },
    onError: (error) => {
      const message = error.response?.data?.error || 'Failed to create session';
      showToast(message, 'error');
    },
  });

  const showToast = (message, type = 'info') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 5000);
  };

  // Calculate isParticipant early for hook
  const isParticipant = user && session && 
    session.participants?.some(p => p.username === user?.username);
  
  // Use heartbeat hook - must be called before any returns
  useHeartbeat(session?.id, user?.username, isParticipant);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <Pizza className="w-12 h-12 animate-spin text-orange-500 mx-auto mb-4" />
          <p className="text-gray-600">Loading lunch session...</p>
        </div>
      </div>
    );
  }
  
  const winningRestaurant = session?.restaurants
    ?.sort((a, b) => (b.voters?.length || 0) - (a.voters?.length || 0))[0];

  return (
    <div className="max-w-6xl mx-auto p-3 space-y-4">
      {/* Header */}
      <header className="text-center py-4">
        <div className="flex items-center justify-center gap-2 mb-2">
          <Pizza className="w-8 h-8 text-orange-500" />
          <h1 className="text-3xl font-bold text-gray-900">LunchRush</h1>
        </div>
        <p className="text-sm text-gray-600">
          Coordinate your team's lunch orders efficiently
        </p>
      </header>

      {/* Session Status */}
      {session && <SessionHeader session={session} />}

      {/* Main Content */}
      {error || !session ? (
        <div className="bg-white rounded-lg shadow-sm p-6 text-center">
          <ChefHat className="w-12 h-12 text-gray-400 mx-auto mb-3" />
          <h2 className="text-xl font-semibold mb-2">No lunch session today</h2>
          <p className="text-sm text-gray-600 mb-4">
            Start coordinating today's lunch orders
          </p>
          <button
            onClick={() => createSessionMutation.mutate()}
            disabled={createSessionMutation.isLoading}
            className="bg-orange-500 text-white px-5 py-2.5 rounded-lg font-medium text-sm
                     hover:bg-orange-600 transition-colors disabled:opacity-50"
          >
            {createSessionMutation.isLoading ? 'Creating...' : 'Create Today\'s Session'}
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
          {/* Left Column */}
          <div className="lg:col-span-2 space-y-4">
            {/* Join Section */}
            {!isParticipant && session.status === 'open' && (
              <JoinSession session={session} onToast={showToast} />
            )}

            {/* Restaurant Section */}
            {(isParticipant || session.status === 'locked') && (
              <RestaurantSection 
                sessionId={session.id}
                currentUser={user?.username}
                isLocked={session.locked}
              />
            )}

            {/* Participants List */}
            {(isParticipant || session.status === 'locked') && (
              <ParticipantsList 
                sessionId={session.id}
                currentUser={user?.username}
                winningRestaurant={winningRestaurant}
              />
            )}
          </div>

          {/* Right Column - Order Summary */}
          <div className="lg:col-span-1">
            {(isParticipant || session.status === 'locked') && (
              <OrderSummary 
                sessionId={session.id}
                session={session}
                currentUser={user?.username}
                winningRestaurant={winningRestaurant}
              />
            )}
          </div>
        </div>
      )}

      {/* Toast Notifications */}
      {toast && (
        <Toast 
          message={toast.message} 
          type={toast.type} 
          onClose={() => setToast(null)} 
        />
      )}
    </div>
  );
};

export default LunchRush;