import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import LunchRush from './components/LunchRush';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchInterval: 2000, // Poll every 2 seconds
      retry: (failureCount, error) => {
        // Don't retry on 404 (no session found)
        if (error?.response?.status === 404) return false;
        return failureCount < 1;
      },
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="min-h-screen bg-gray-50">
        <LunchRush />
      </div>
    </QueryClientProvider>
  );
}

export default App;