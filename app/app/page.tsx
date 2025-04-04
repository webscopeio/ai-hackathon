"use client";

import { useState, useEffect } from "react";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { api } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";
import { ConfigForm } from "@/components/ConfigForm";
import { AppLayout } from "@/components/AppLayout";
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
}

const TestGenerationPanel: React.FC = () => {
  const [config, setConfig] = useState<ConfigState | null>(null);
  const [isGenerating, setIsGenerating] = useState(false);
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

  const handleGenerateTests = async () => {
    setIsGenerating(true);
    setGenerationStatus({ status: 'generating' });
    
    try {
      // In a real implementation, this would call the Go backend API
      // For now, we'll just simulate a successful response
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      setGenerationStatus({ 
        status: 'success', 
        message: 'Test generation initiated successfully', 
        jobId: `job-${Date.now()}` 
      });
      toast.success('Test generation started');
    } catch (error) {
      console.error('Error generating tests:', error);
      setGenerationStatus({ 
        status: 'error', 
        message: error instanceof Error ? error.message : 'Unknown error occurred'
      });
      toast.error('Failed to generate tests');
    } finally {
      setIsGenerating(false);
    }
  };

  return (
    <div className="space-y-6">
      <h2 className="mb-6 scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight transition-colors">
        Generate Tests
      </h2>
      <div className="bg-white rounded-lg shadow-md overflow-hidden border border-gray-200">
        <div className="p-6">
          <h3 className="text-xl font-semibold tracking-tight mb-4">
            Test Generation Ready
          </h3>
          <p className="leading-7 text-gray-600 mb-4">
            Your configuration is complete. You can now generate E2E tests for your application.
          </p>
          <div className="bg-gray-50 p-4 rounded-md mb-4">
            <h4 className="font-medium mb-2">Tech Specification</h4>
            <p className="text-sm text-gray-600 whitespace-pre-line">{config?.techSpecification || ''}</p>
          </div>
          
          {generationStatus.status === 'success' && (
            <div className="bg-green-50 border-l-4 border-green-500 p-4 rounded mb-4">
              <p className="text-green-700">{generationStatus.message}</p>
              <p className="text-sm text-green-600 mt-1">Job ID: {generationStatus.jobId}</p>
            </div>
          )}
          
          {generationStatus.status === 'error' && (
            <div className="bg-red-50 border-l-4 border-red-500 p-4 rounded mb-4">
              <p className="text-red-700">Error: {generationStatus.message}</p>
            </div>
          )}
          
          <Button 
            onClick={handleGenerateTests}
            disabled={isGenerating}
            className="bg-blue-600 hover:bg-blue-700"
          >
            {isGenerating ? (
              <>
                <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Generating...
              </>
            ) : 'Generate Tests'}
          </Button>
        </div>
      </div>
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
