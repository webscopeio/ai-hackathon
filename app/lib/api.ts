import { queryOptions, useMutation } from "@tanstack/react-query";
import { fetchGet, fetchPost } from "./api.helpers";

export const api = {
  getStatus: queryOptions({
    queryKey: ["getStatus"],
    queryFn: () => fetchGet<StatusReturn>("status"),
  }),
  getConfig: queryOptions({
    queryKey: ["getConfig"],
    queryFn: () => fetchGet<ConfigReturn>("config"),
  }),
};

export const useUpdateConfig = () => {
  return useMutation({
    mutationFn: (config: Partial<ConfigData>) => 
      fetchPost<Partial<ConfigData>, ConfigUpdateReturn>("config", config),
  });
};

export const useGenerateTests = () => {
  return useMutation({
    mutationFn: (config: ConfigData) => 
      fetchPost<ConfigData, GenerateTestsReturn>("generate-tests", config),
  });
};

export type GenerateTestsReturn = {
  success: boolean;
  message: string;
  jobId: string;
};

export type StatusReturn = {
  status: string;
};

export type ConfigData = {
  anthropicApiKey: string;
  sentryApiKey: string;
  techSpecification: string;
  productSpecification: string;
};

export type ConfigReturn = ConfigData;

export type ConfigUpdateReturn = {
  success: boolean;
  config: ConfigData;
};

export type ErrorReturn = {
  error: string;
};
