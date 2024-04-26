const API_URL = "/api/v1";

async function apiRequest(
  path: string,
  method: "GET" | "POST" | "PUT" | "DELETE",
  body?: unknown,
) {
  const res = await fetch(API_URL + path, {
    method,
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    body: JSON.stringify(body),
  });

  if (!res.ok) {
    throw new Error(
      `Failed to fetch '${API_URL + path}': ${res.status} ${res.statusText}`,
    );
  }

  return res.json();
}

export default apiRequest;
