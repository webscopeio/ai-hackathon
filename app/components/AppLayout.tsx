"use client";

import { ReactNode, useState, useEffect } from "react";
import { fetchGet } from "@/lib/api.helpers";

interface AppLayoutProps {
  configSection: ReactNode;
  mainSection: ReactNode;
}

interface ConfigState {
  anthropicApiKey: string;
  sentryApiKey: string;
  techSpecification: string;
}

export function AppLayout({ configSection, mainSection }: AppLayoutProps) {
  const [isConfigComplete, setIsConfigComplete] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  // Check if config is complete on mount
  useEffect(() => {
    const checkConfig = async () => {
      try {
        const config = await fetchGet<ConfigState>("config");
        const complete = Boolean(
          config?.anthropicApiKey && 
          config?.sentryApiKey && 
          config?.techSpecification
        );
        setIsConfigComplete(complete);
      } catch (error) {
        console.error("Error checking config:", error);
        setIsConfigComplete(false);
      } finally {
        setIsLoading(false);
      }
    };

    checkConfig();
  }, []);

  return (
    <div className="container mx-auto py-8 px-4">
      <h1 className="text-4xl font-bold mb-8">E2E Test Generator</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-12 gap-8">
        <div className="md:col-span-4 space-y-6">
          <div className="bg-white rounded-lg shadow-md p-6 border border-gray-200">
            <h2 className="text-2xl font-semibold mb-4">Configuration</h2>
            {configSection}
          </div>
        </div>

        <div className="md:col-span-8">
          {isLoading ? (
            <div className="bg-gray-50 border border-gray-200 p-6 rounded shadow animate-pulse">
              <div className="h-4 bg-gray-200 rounded w-3/4 mb-4"></div>
              <div className="h-4 bg-gray-200 rounded w-1/2"></div>
            </div>
          ) : !isConfigComplete ? (
            <div className="bg-amber-50 border-l-4 border-amber-500 p-6 rounded shadow">
              <div className="ml-3">
                <h3 className="text-xl font-medium text-amber-800 mb-2">
                  Configuration Required
                </h3>
                <p className="leading-7 text-amber-700">
                  Please complete the configuration form to enable test generation.
                </p>
              </div>
            </div>
          ) : (
            mainSection
          )}
        </div>
      </div>
    </div>
  );
}
