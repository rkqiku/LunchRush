import React, { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { ShoppingCartIcon, LockClosedIcon, UserIcon } from '@heroicons/react/24/outline';
import api from '../api/client';

const OrderSummary = ({ sessionId, session, currentUser, winningRestaurant }) => {
  const [meal, setMeal] = useState('');
  const queryClient = useQueryClient();

  const orderMutation = useMutation({
    mutationFn: (mealText) => api.orderMeal(sessionId, currentUser, mealText),
    onSuccess: () => {
      queryClient.invalidateQueries(['session', sessionId]);
      setMeal('');
    },
  });

  const selectOrderPlacerMutation = useMutation({
    mutationFn: () => api.selectOrderPlacer(sessionId, currentUser),
    onSuccess: () => {
      queryClient.invalidateQueries(['session', sessionId]);
    },
  });

  const lockSessionMutation = useMutation({
    mutationFn: () => api.lockSession(sessionId),
    onSuccess: () => {
      queryClient.invalidateQueries(['session', sessionId]);
    },
  });

  const currentParticipant = session?.participants?.find(
    p => p.username === currentUser
  );
  const hasOrdered = currentParticipant?.meal && currentParticipant.meal.trim() !== '';
  const orderPlacer = session?.participants?.find(p => p.isOrderPlacer);
  const allHaveOrdered = session?.participants?.every(p => p.meal && p.meal.trim() !== '');

  const handleOrder = (e) => {
    e.preventDefault();
    if (meal.trim() && !session?.locked) {
      orderMutation.mutate(meal.trim());
    }
  };

  if (!winningRestaurant) {
    return (
      <div className="bg-white rounded-lg shadow-sm p-4">
        <div className="text-center py-6">
          <ShoppingCartIcon className="h-10 w-10 text-gray-400 mx-auto mb-3" />
          <p className="text-gray-500 text-sm">
            Waiting for restaurant selection...
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-sm p-4">
      <h2 className="text-xl font-semibold text-gray-800 mb-3">
        Order from {winningRestaurant.name}
      </h2>

      {!session?.locked && !hasOrdered && (
        <form onSubmit={handleOrder} className="mb-6">
          <div className="space-y-3">
            <textarea
              value={meal}
              onChange={(e) => setMeal(e.target.value)}
              placeholder="Enter your order..."
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 resize-none text-sm"
            />
            <button
              type="submit"
              disabled={!meal.trim() || orderMutation.isLoading}
              className="w-full bg-green-500 text-white px-4 py-2 rounded-lg hover:bg-green-600 disabled:bg-gray-300 transition-colors flex items-center justify-center gap-2 text-sm"
            >
              <ShoppingCartIcon className="h-4 w-4" />
              Confirm Order
            </button>
          </div>
        </form>
      )}

      {hasOrdered && (
        <div className="mb-6 p-4 bg-green-50 border border-green-200 rounded-lg">
          <h3 className="font-semibold text-green-800 mb-2 text-sm">Your order:</h3>
          <p className="text-gray-700">{currentParticipant.meal}</p>
        </div>
      )}

      <div className="space-y-4">
        {/* Order Placer Section */}
        <div className="p-3 bg-gray-50 rounded-lg border border-gray-200">
          <h3 className="font-semibold text-gray-800 mb-2 flex items-center gap-2 text-sm">
            <UserIcon className="h-4 w-4" />
            Order Placer
          </h3>
          
          {orderPlacer ? (
            <div className="flex items-center justify-between">
              <p className="text-gray-700">
                <span className="font-medium">{orderPlacer.username}</span>
                {orderPlacer.username === currentUser && (
                  <span className="text-sm text-gray-500 ml-2">(You)</span>
                )}
              </p>
              {orderPlacer.username === currentUser && allHaveOrdered && !session?.locked && (
                <button
                  onClick={() => lockSessionMutation.mutate()}
                  disabled={lockSessionMutation.isLoading}
                  className="bg-red-500 text-white px-4 py-2 rounded-lg hover:bg-red-600 disabled:bg-gray-300 transition-colors flex items-center gap-2 text-sm"
                >
                  <LockClosedIcon className="h-4 w-4" />
                  Lock Order
                </button>
              )}
            </div>
          ) : (
            <div>
              <p className="text-gray-600 mb-3 text-sm">
                No one has volunteered to place the order yet.
              </p>
              {!session?.locked && (
                <button
                  onClick={() => selectOrderPlacerMutation.mutate()}
                  disabled={selectOrderPlacerMutation.isLoading}
                  className="bg-orange-500 text-white px-4 py-2 rounded-lg hover:bg-orange-600 disabled:bg-gray-300 transition-colors text-sm"
                >
                  Volunteer to Order
                </button>
              )}
            </div>
          )}
        </div>

        {/* Session Status */}
        {session?.locked && (
          <div className="p-3 bg-red-50 border border-red-200 rounded-lg">
            <div className="flex items-center gap-2 text-red-800">
              <LockClosedIcon className="h-4 w-4" />
              <p className="font-semibold text-sm">Session Locked</p>
            </div>
            <p className="text-xs text-red-700 mt-1">
              The order has been finalized and sent to the restaurant.
            </p>
          </div>
        )}

        {/* Order Summary */}
        <div className="pt-4 border-t border-gray-200">
          <h3 className="font-semibold text-gray-800 mb-2 text-sm">Order Summary</h3>
          <div className="space-y-2">
            {session?.participants?.filter(p => p.meal && p.meal.trim() !== '').map(p => (
              <div key={p.username} className="flex justify-between text-sm">
                <span className="font-medium text-gray-700">{p.username}:</span>
                <span className="text-gray-600">{p.meal}</span>
              </div>
            ))}
          </div>
          
          {!allHaveOrdered && !session?.locked && (
            <p className="text-xs text-yellow-600 mt-3">
              Waiting for all participants to order...
            </p>
          )}
        </div>
      </div>
    </div>
  );
};

export default OrderSummary;