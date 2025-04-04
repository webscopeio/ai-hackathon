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
  analyze: (args: AnalyzeArgs) =>
    fetchPost<AnalyzeArgs, AnalyzeReturn>("analyze", args),
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

export type AnalyzeArgs = {
  url: string;
  prompt: string;
};

export type AnalyzeReturn = {
  techSpec: string;
  siteMap: Record<string, string>;
  criteria: string;
};

export type ErrorReturn = {
  error: string;
};
