"use client";
import type * as React from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "@/components/ui/sonner";
import { ConfigProvider } from "@/lib/config-context";

let browserQueryClient: QueryClient | undefined = undefined;

function getQueryClient() {
  if (!browserQueryClient) {
    browserQueryClient = new QueryClient({
      defaultOptions: {
        queries: {
          staleTime: 60 * 1000,
        },
      },
    });
  }
  return browserQueryClient;
}

export default function Providers({ children }: { children: React.ReactNode }) {
  const queryClient = getQueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      <ConfigProvider>
        {children}
        <Toaster />
      </ConfigProvider>
    </QueryClientProvider>
  );
}
