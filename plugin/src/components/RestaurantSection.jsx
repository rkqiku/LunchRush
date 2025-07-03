import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { PlusIcon, TrashIcon, HandThumbUpIcon } from '@heroicons/react/24/outline';
import { HandThumbUpIcon as HandThumbUpSolidIcon } from '@heroicons/react/24/solid';
import api from '../api/client';

const RestaurantSection = ({ sessionId, currentUser, isLocked }) => {
  const [newRestaurant, setNewRestaurant] = useState('');
  const queryClient = useQueryClient();

  const { data: restaurants = [] } = useQuery({
    queryKey: ['restaurants', sessionId],
    queryFn: () => api.getRestaurants(sessionId),
  });

  const proposeMutation = useMutation({
    mutationFn: (name) => api.proposeRestaurant(sessionId, name, currentUser),
    onSuccess: () => {
      queryClient.invalidateQueries(['restaurants', sessionId]);
      setNewRestaurant('');
    },
  });

  const voteMutation = useMutation({
    mutationFn: (restaurantId) => api.voteRestaurant(sessionId, restaurantId, currentUser),
    onSuccess: () => {
      queryClient.invalidateQueries(['restaurants', sessionId]);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (restaurantId) => api.deleteRestaurant(sessionId, restaurantId),
    onSuccess: () => {
      queryClient.invalidateQueries(['restaurants', sessionId]);
    },
  });

  const handlePropose = (e) => {
    e.preventDefault();
    if (newRestaurant.trim() && !isLocked) {
      proposeMutation.mutate(newRestaurant.trim());
    }
  };

  const hasVoted = (restaurant) => {
    return restaurant.voters?.includes(currentUser);
  };

  const sortedRestaurants = Array.isArray(restaurants) 
    ? [...restaurants].sort((a, b) => (b.voters?.length || 0) - (a.voters?.length || 0))
    : [];

  const winningRestaurant = sortedRestaurants[0];
  const hasWinner = winningRestaurant && (winningRestaurant.voters?.length || 0) > 0;

  return (
    <div className="bg-white rounded-lg shadow-sm p-4">
      <h2 className="text-xl font-semibold text-gray-800 mb-3">Restaurants</h2>
      
      {!isLocked && (
        <form onSubmit={handlePropose} className="mb-6">
          <div className="flex gap-2">
            <input
              type="text"
              value={newRestaurant}
              onChange={(e) => setNewRestaurant(e.target.value)}
              placeholder="Suggest a restaurant..."
              className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 text-sm"
            />
            <button
              type="submit"
              disabled={!newRestaurant.trim() || proposeMutation.isLoading}
              className="bg-orange-500 text-white px-4 py-2 rounded-lg hover:bg-orange-600 disabled:bg-gray-300 transition-colors flex items-center gap-2 text-sm"
            >
              <PlusIcon className="h-4 w-4" />
              Add
            </button>
          </div>
        </form>
      )}

      <div className="space-y-3">
        {sortedRestaurants.length === 0 ? (
          <p className="text-gray-500 text-center py-4 text-sm">
            No restaurants suggested yet. Be the first!
          </p>
        ) : (
          sortedRestaurants.map((restaurant, index) => (
            <div
              key={restaurant.id}
              className={`p-3 rounded-lg border transition-all ${
                hasWinner && index === 0
                  ? 'border-green-500 bg-green-50 shadow-sm'
                  : 'border-gray-200 hover:border-gray-300'
              }`}
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <h3 className="font-semibold text-base text-gray-800">
                    {restaurant.name}
                  </h3>
                  {hasWinner && index === 0 && (
                    <span className="bg-green-500 text-white text-xs px-2 py-1 rounded-full">
                      Winner
                    </span>
                  )}
                </div>
                
                <div className="flex items-center gap-2">
                  <span className="text-gray-600 font-medium text-sm">
                    {restaurant.voters?.length || 0} votes
                  </span>
                  
                  {!isLocked && (
                    <>
                      <button
                        onClick={() => voteMutation.mutate(restaurant.id)}
                        disabled={voteMutation.isLoading}
                        className={`flex items-center gap-1 px-3 py-1.5 rounded-lg transition-all text-sm font-medium ${
                          hasVoted(restaurant)
                            ? 'bg-orange-500 text-white hover:bg-orange-600'
                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                        }`}
                      >
                        {hasVoted(restaurant) ? (
                          <HandThumbUpSolidIcon className="h-4 w-4" />
                        ) : (
                          <HandThumbUpIcon className="h-4 w-4" />
                        )}
                        Vote
                      </button>
                      
                      {restaurant.proposedBy === currentUser && (
                        <button
                          onClick={() => deleteMutation.mutate(restaurant.id)}
                          disabled={deleteMutation.isLoading}
                          className="p-1.5 rounded-lg hover:bg-red-100 transition-colors"
                          title="Delete restaurant"
                        >
                          <TrashIcon className="h-4 w-4 text-red-600" />
                        </button>
                      )}
                    </>
                  )}
                </div>
              </div>
              
              {restaurant.voters && restaurant.voters.length > 0 && (
                <div className="mt-2 flex flex-wrap gap-1">
                  {restaurant.voters.map((voter) => (
                    <span
                      key={voter}
                      className="text-xs bg-gray-200 px-2 py-1 rounded-full text-gray-600"
                    >
                      {voter}
                    </span>
                  ))}
                </div>
              )}
              
              <p className="text-xs text-gray-500 mt-1">
                Suggested by {restaurant.proposedBy}
              </p>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default RestaurantSection;