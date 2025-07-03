import React from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { UserGroupIcon, CheckCircleIcon, ClockIcon, XCircleIcon } from '@heroicons/react/24/outline';
import { WifiIcon, XMarkIcon } from '@heroicons/react/24/solid';
import api from '../api/client';

const ParticipantsList = ({ sessionId, currentUser, winningRestaurant }) => {
  const queryClient = useQueryClient();
  const { data: session } = useQuery({
    queryKey: ['session', 'today'],
    queryFn: () => api.getSession(sessionId),
  });
  
  const removeParticipantMutation = useMutation({
    mutationFn: (username) => api.removeParticipant(sessionId, username),
    onSuccess: () => {
      queryClient.invalidateQueries(['session', 'today']);
    },
  });

  const participants = session?.participants || [];
  const hasOrdered = (participant) => {
    return participant.meal && participant.meal.trim() !== '';
  };

  const getParticipantStatus = (participant) => {
    if (participant.isOrderPlacer) {
      return { icon: UserGroupIcon, color: 'text-blue-600', bg: 'bg-blue-100' };
    }
    if (hasOrdered(participant)) {
      return { icon: CheckCircleIcon, color: 'text-green-600', bg: 'bg-green-100' };
    }
    return { icon: ClockIcon, color: 'text-yellow-600', bg: 'bg-yellow-100' };
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-4">
      <h2 className="text-xl font-semibold text-gray-800 mb-3">
        Participants ({participants.length})
      </h2>

      <div className="space-y-3">
        {participants.length === 0 ? (
          <p className="text-gray-500 text-center py-4 text-sm">
            No participants yet.
          </p>
        ) : (
          participants.map((participant) => {
            const status = getParticipantStatus(participant);
            const StatusIcon = status.icon;
            
            return (
              <div
                key={participant.username}
                className={`p-3 rounded-lg border transition-all ${
                  participant.username === currentUser
                    ? 'border-blue-300 bg-blue-50'
                    : 'border-gray-200 bg-gray-50'
                } ${
                  participant.isActive === false ? 'opacity-60' : ''
                }`}
              >
                <div className="flex items-start justify-between">
                  <div className="flex items-center gap-3">
                    <div className={`p-2 rounded-full ${status.bg}`}>
                      <StatusIcon className={`h-5 w-5 ${status.color}`} />
                    </div>
                    <div>
                      <h3 className="font-semibold text-gray-800 flex items-center gap-2">
                        {participant.username}
                        {participant.username === currentUser && (
                          <span className="text-sm text-gray-500">(You)</span>
                        )}
                        {participant.isActive === true && (
                          <WifiIcon className="h-4 w-4 text-green-500" title="Online" />
                        )}
                        {participant.isActive === false && (
                          <span className="text-xs text-gray-500 italic">Inactive</span>
                        )}
                      </h3>
                      {participant.isOrderPlacer && (
                        <span className="text-xs text-blue-600 font-medium">
                          Order Placer
                        </span>
                      )}
                    </div>
                  </div>
                  
                  <div className="flex items-center gap-2">
                    {hasOrdered(participant) && (
                      <CheckCircleIcon className="h-5 w-5 text-green-500" />
                    )}
                    {participant.isActive === false && !session?.locked && (
                      <button
                        onClick={() => removeParticipantMutation.mutate(participant.username)}
                        className="p-1 hover:bg-red-100 rounded transition-colors"
                        title="Remove inactive participant"
                      >
                        <XMarkIcon className="h-4 w-4 text-red-600" />
                      </button>
                    )}
                  </div>
                </div>

                {participant.meal && (
                  <div className="mt-2 p-2 bg-gray-50 rounded border border-gray-200">
                    <p className="text-xs font-medium text-gray-700">Order:</p>
                    <p className="text-sm text-gray-600">{participant.meal}</p>
                    {winningRestaurant && (
                      <p className="text-xs text-gray-500 mt-1">
                        from {winningRestaurant.name}
                      </p>
                    )}
                  </div>
                )}

                {!hasOrdered(participant) && participant.username === currentUser && (
                  <p className="text-xs text-yellow-600 mt-2">
                    You haven't ordered yet!
                  </p>
                )}
              </div>
            );
          })
        )}
      </div>

      <div className="mt-4 pt-4 border-t border-gray-200">
        <div className="flex justify-between text-sm">
          <span className="text-gray-600">Total participants:</span>
          <span className="font-semibold">{participants.length}</span>
        </div>
        <div className="flex justify-between text-sm mt-1">
          <span className="text-gray-600">Ordered:</span>
          <span className="font-semibold text-green-600">
            {participants.filter(hasOrdered).length}
          </span>
        </div>
        <div className="flex justify-between text-sm mt-1">
          <span className="text-gray-600">Pending:</span>
          <span className="font-semibold text-yellow-600">
            {participants.filter(p => !hasOrdered(p)).length}
          </span>
        </div>
        <div className="flex justify-between text-sm mt-1">
          <span className="text-gray-600">Active:</span>
          <span className="font-semibold text-green-600">
            {participants.filter(p => p.isActive !== false).length}
          </span>
        </div>
      </div>
    </div>
  );
};

export default ParticipantsList;