import { Clock, Lock, Unlock } from 'lucide-react';
import { useEffect, useState } from 'react';

const SessionHeader = ({ session }) => {
  const [timeLeft, setTimeLeft] = useState('');

  useEffect(() => {
    const updateTimer = () => {
      if (session.status === 'locked') {
        setTimeLeft('Order Locked');
        return;
      }

      const lockTime = new Date(session.lockTime);
      const now = new Date();
      const diff = lockTime - now;

      if (diff <= 0) {
        setTimeLeft('Lock time passed');
      } else {
        const hours = Math.floor(diff / (1000 * 60 * 60));
        const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
        setTimeLeft(`${hours}h ${minutes}m`);
      }
    };

    updateTimer();
    const interval = setInterval(updateTimer, 60000); // Update every minute

    return () => clearInterval(interval);
  }, [session]);

  return (
    <div className="bg-white rounded-lg shadow-sm p-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className={`flex items-center gap-2 px-3 py-1.5 rounded-full font-medium text-sm
                          ${session.status === 'open' 
                            ? 'bg-green-100 text-green-800' 
                            : 'bg-red-100 text-red-800'}`}>
            {session.status === 'open' ? (
              <>
                <Unlock className="w-4 h-4" />
                <span>Open</span>
              </>
            ) : (
              <>
                <Lock className="w-4 h-4" />
                <span>Locked</span>
              </>
            )}
          </div>
          
          <div className="text-gray-600 text-sm">
            <span className="font-medium">{session.participants.length}</span> participants
          </div>
        </div>

        <div className="flex items-center gap-2 text-gray-700">
          <Clock className="w-4 h-4" />
          <span className="font-medium text-sm">{timeLeft}</span>
        </div>
      </div>
    </div>
  );
};

export default SessionHeader;