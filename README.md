# AI Hackathon

Full-stack demo application with Next frontend and Go backend.

## Tech Stack

- **Frontend**: Next with React Query, Zod, and Shadcn UI
- **Backend**: Go with Chi router

## Getting Started

```bash
# Run only backend with hot reload
make dev/server

# Run only frontend
make dev/app

# Run both frontend and backend
make dev
```

## API Integration

Frontend uses typed API client (`lib/api.ts`) with React Query integration for data fetching:

```typescript
// Query example
const { data } = useQuery(api.getPosts);

// Mutation example
const mutation = useMutation({
  mutationFn: api.greet,
  onSuccess: (data) => toast.success(data.message),
});
```

> [!WARNING]  
> Endpoints are manually typed between Go backend and TypeScript frontend.
