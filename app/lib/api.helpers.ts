import { ErrorReturn } from "./api";

export const API_BASE = "/_api";

export async function throwServerError(response: Response): Promise<void> {
  const errorData: ErrorReturn = await response
    .json()
    .catch(() => ({ error: "Unknown error" }));
  throw new Error(
    errorData.error || `Request failed with status ${response.status}`,
  );
}

export async function fetchGet<TReturn, TArgs = undefined>(
  path: string,
  args?: TArgs,
  options?: RequestInit,
): Promise<TReturn> {
  let url = `${API_BASE}/${path}`;

  if (args) {
    const queryParams = new URLSearchParams();
    Object.entries(args).forEach(([key, value]) => {
      if (value !== undefined) {
        queryParams.append(key, String(value));
      }
    });

    const queryString = queryParams.toString();
    if (queryString) {
      url += `?${queryString}`;
    }
  }

  const response = await fetch(url, options);
  if (!response.ok) await throwServerError(response);
  return await response.json();
}

export async function fetchPost<TArgs, TReturn>(
  path: string,
  args: TArgs,
  options?: RequestInit,
): Promise<TReturn> {
  const response = await fetch(`${API_BASE}/${path}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
    ...options,
    body: JSON.stringify(args),
  });
  if (!response.ok) await throwServerError(response);
  return await response.json();
}
