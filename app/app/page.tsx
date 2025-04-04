"use client";

import { useState, useEffect } from "react";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { api } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";
import { ConfigForm } from "@/components/ConfigForm";
import { AppLayout } from "@/components/AppLayout";
import { TestGenerator } from "@/components/TestGenerator";
import { toast } from "sonner";
import { fetchGet } from "@/lib/api.helpers";

const StatusDemo: React.FC = () => {
  const { status, data, error } = useQuery(api.getStatus);

  return (
    <div className="flex-1">
      <h2 className="mb-6 scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight transition-colors">
        System Status
      </h2>
      {status === "pending" ? (
        <div className="bg-white rounded-lg shadow-md overflow-hidden border border-gray-200">
          <div className="p-6">
            <Skeleton className="h-6 w-3/4 mb-4" />
            <Skeleton className="h-4 w-full mb-2" />
          </div>
        </div>
      ) : status === "error" ? (
        <div className="bg-red-50 border-l-4 border-red-500 p-6 rounded shadow">
          <div className="ml-3">
            <h3 className="text-xl font-medium text-red-800 mb-2">
              Error loading status
            </h3>
            <p className="leading-7 text-red-700">
              {error?.message ?? "Unknown error"}
            </p>
          </div>
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow-md overflow-hidden border border-gray-200 hover:shadow-lg transition-shadow duration-300">
          <div className="p-6">
            <h3 className="text-xl font-semibold tracking-tight mb-2">
              System Status
            </h3>
            <p className="leading-7 text-gray-600">{JSON.stringify(data)}</p>
          </div>
        </div>
      )}
    </div>
  );
};

interface ConfigState {
  anthropicApiKey: string;
  sentryApiKey: string;
  techSpecification: string;
  productSpecification: string;
}

const TestGenerationPanel: React.FC = () => {
  const [config, setConfig] = useState<ConfigState | null>(null);
  const [generationStatus, setGenerationStatus] = useState<{
    status: 'idle' | 'generating' | 'success' | 'error';
    message?: string;
    jobId?: string;
  }>({ status: 'idle' });
  
  // Fetch config on mount
  useEffect(() => {
    const fetchConfig = async () => {
      try {
        const configData = await fetchGet<ConfigState>("config");
        setConfig(configData);
      } catch (error) {
        console.error("Error fetching config:", error);
      }
    };
    
    fetchConfig();
  }, []);

  const handleGenerateTests = async (prompt: string) => {
    setGenerationStatus({ status: 'generating' });
    
    try {
      console.log('Generating tests with prompt:', prompt);
      
      // In a real implementation, this would call the Go backend API with the prompt
      // For now, we'll just simulate a successful response
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      setGenerationStatus({ 
        status: 'success', 
        message: 'Test generation initiated successfully', 
        jobId: `job-${Date.now()}` 
      });
    } catch (error) {
      console.error('Error generating tests:', error);
      setGenerationStatus({ 
        status: 'error', 
        message: error instanceof Error ? error.message : 'Unknown error occurred'
      });
      throw error; // Re-throw to be handled by the TestGenerator component
    }
  };

  return (
    <div className="space-y-8">
      
      {!config ? (
        <div className="space-y-4">
          <Skeleton className="h-8 w-1/3" />
          <Skeleton className="h-32 w-full" />
          <Skeleton className="h-32 w-full" />
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow-md overflow-hidden border border-gray-200">
          <div className="p-6">
            {generationStatus.status === 'success' && (
              <div className="bg-green-50 border-l-4 border-green-500 p-4 rounded mb-6">
                <p className="text-green-700">{generationStatus.message}</p>
                <p className="text-sm text-green-600 mt-1">Job ID: {generationStatus.jobId}</p>
              </div>
            )}
            
            {generationStatus.status === 'error' && (
              <div className="bg-red-50 border-l-4 border-red-500 p-4 rounded mb-6">
                <p className="text-red-700">Error: {generationStatus.message}</p>
              </div>
            )}
            
            <TestGenerator
              techSpecification={config.techSpecification || ''}
              productSpecification={config.productSpecification || ''}
              onGenerateTests={handleGenerateTests}
            />
          </div>
        </div>
      )}
      
      <StatusDemo />
    </div>
  );
};

export default function Home() {
  return (
    <AppLayout 
      configSection={<ConfigForm />}
      mainSection={<TestGenerationPanel />}
    />
  );
}
