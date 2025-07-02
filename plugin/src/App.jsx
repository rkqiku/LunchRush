import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import LunchRush from './components/LunchRush';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchInterval: 2000, // Poll every 2 seconds
      retry: 1,
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