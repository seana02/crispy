const API_BASE_URL = '/api';

interface FetchError extends Error {
  status: number;
  message: string;
}

// Normalise backend snake_case pagination meta to camelCase totalPages.
function normalisePagination(data: unknown): unknown {
  if (data && typeof data === 'object' && !Array.isArray(data)) {
    const obj = data as Record<string, unknown>;
    if (obj.meta && typeof obj.meta === 'object') {
      const meta = obj.meta as Record<string, unknown>;
      if ('total_pages' in meta && !('totalPages' in meta)) {
        meta.totalPages = meta.total_pages;
      }
    }
  }
  return data;
}

export async function fetchApi<T>(url: string, options?: RequestInit): Promise<T> {
  const fullUrl = `${API_BASE_URL}${url.startsWith('/') ? url : `/${url}`}`;
  const method = options?.method ?? 'GET';

  const defaultHeaders: HeadersInit = {
    'Content-Type': 'application/json',
  };

  const mergedOptions: RequestInit = {
    ...options,
    headers: {
      ...defaultHeaders,
      ...options?.headers,
    },
  };

  console.debug(`[api] ${method} ${fullUrl}`);

  let response: Response;
  try {
    response = await fetch(fullUrl, mergedOptions);
  } catch (networkErr) {
    console.error(`[api] Network error on ${method} ${fullUrl}:`, networkErr);
    throw networkErr;
  }

  if (!response.ok) {
    let bodyText = '';
    let errorMessage = `${method} ${fullUrl} → ${response.status} ${response.statusText}`;
    try {
      bodyText = await response.text();
      const body = JSON.parse(bodyText);
      if (body.message) errorMessage = body.message;
    } catch {
      if (bodyText) errorMessage += ` — ${bodyText}`;
    }
    console.error(`[api] ${errorMessage}`, { method, url: fullUrl, status: response.status, body: bodyText });
    const error = new Error(errorMessage) as FetchError;
    error.status = response.status;
    throw error;
  }

  if (response.status === 204) {
    console.debug(`[api] ${method} ${fullUrl} → 204 No Content`);
    return undefined as T;
  }

  const json = await response.json();
  console.debug(`[api] ${method} ${fullUrl} → ${response.status}`, json);
  return normalisePagination(json) as T;
}
