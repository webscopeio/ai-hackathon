"use client";

import { createContext, useContext, useState, ReactNode } from "react";

export interface ConfigState {
  anthropicApiKey: string;
  sentryApiKey: string;
  umamiAPIKey: string;
  umamiWebsiteId: string;
  techSpecification: string;
}

interface ConfigContextType {
  config: ConfigState;
  updateConfig: (newConfig: Partial<ConfigState>) => void;
  isConfigComplete: boolean;
}

const defaultConfig: ConfigState = {
  anthropicApiKey: "",
  sentryApiKey: "",
  umamiAPIKey: "",
  umamiWebsiteId: "",
  techSpecification: "",
};

const ConfigContext = createContext<ConfigContextType | undefined>(undefined);

export function ConfigProvider({ children }: { children: ReactNode }) {
  const [config, setConfig] = useState<ConfigState>(defaultConfig);

  const updateConfig = (newConfig: Partial<ConfigState>) => {
    setConfig((prev) => ({ ...prev, ...newConfig }));
  };

  // Check if all required fields are filled
  const isConfigComplete = Boolean(
    config.anthropicApiKey && config.sentryApiKey && 
    config.umamiAPIKey && config.umamiWebsiteId && 
    config.techSpecification
  );

  return (
    <ConfigContext.Provider value={{ config, updateConfig, isConfigComplete }}>
      {children}
    </ConfigContext.Provider>
  );
}

export function useConfig() {
  const context = useContext(ConfigContext);
  if (context === undefined) {
    throw new Error("useConfig must be used within a ConfigProvider");
  }
  return context;
}
