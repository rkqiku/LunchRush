import { useState, useEffect } from 'react';

const useUser = () => {
  const [user, setUser] = useState(null);

  useEffect(() => {
    const savedUser = localStorage.getItem('lunchRushUser');
    if (savedUser) {
      setUser(JSON.parse(savedUser));
    }
  }, []);

  const createUser = (username) => {
    const newUser = {
      userId: `user_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      username,
    };
    
    setUser(newUser);
    localStorage.setItem('lunchRushUser', JSON.stringify(newUser));
    
    return newUser;
  };

  const clearUser = () => {
    setUser(null);
    localStorage.removeItem('lunchRushUser');
  };

  return { user, createUser, clearUser };
};

export default useUser;