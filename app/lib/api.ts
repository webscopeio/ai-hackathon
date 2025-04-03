import { queryOptions } from "@tanstack/react-query";
import { fetchGet, fetchPost } from "./api.helpers";

export const api = {
  getStatus: queryOptions({
    queryKey: ["getStatus"],
    queryFn: () => fetchGet<StatusReturn>("status"),
  }),
};

export type StatusReturn = {
  status: string;
};

export type ErrorReturn = {
  error: string;
};
