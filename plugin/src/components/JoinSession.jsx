import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { UserPlus } from 'lucide-react';
import api from '../api/client';
import useUser from '../hooks/useUser';

const JoinSession = ({ session, onToast }) => {
  const [username, setUsername] = useState('');
  const { createUser } = useUser();
  const queryClient = useQueryClient();

  const joinMutation = useMutation({
    mutationFn: (data) => api.joinSession(session.id, data.username),
    onSuccess: async (data) => {
      console.log('Join successful, response:', data);
      onToast('Successfully joined the session!', 'success');
      // Force page refresh to sync all state from localStorage
      setTimeout(() => {
        window.location.reload();
      }, 500);
    },
    onError: (error) => {
      console.error('Join failed:', error);
      const message = error.response?.data?.error || 'Failed to join session';
      onToast(message, 'error');
    },
  });

  const handleJoin = (e) => {
    e.preventDefault();
    if (!username.trim()) {
      onToast('Please enter your name', 'error');
      return;
    }

    createUser(username);
    joinMutation.mutate({ username });
  };

  return (
    <div className="bg-white rounded-lg shadow-lg p-6">
      <div className="flex items-center gap-3 mb-4">
        <UserPlus className="w-6 h-6 text-orange-500" />
        <h2 className="text-xl font-semibold">Join Today's Lunch</h2>
      </div>
      
      <form onSubmit={handleJoin} className="space-y-4">
        <div>
          <label htmlFor="username" className="block text-sm font-medium text-gray-700 mb-1">
            Your Name
          </label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="Enter your name"
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 
                     focus:ring-orange-500 focus:border-transparent"
            disabled={joinMutation.isLoading}
          />
        </div>
        
        <button
          type="submit"
          disabled={joinMutation.isLoading}
          className="w-full bg-orange-500 text-white py-3 rounded-lg font-medium
                   hover:bg-orange-600 transition-colors disabled:opacity-50"
        >
          {joinMutation.isLoading ? 'Joining...' : 'Join Session'}
        </button>
      </form>
    </div>
  );
};

export default JoinSession;