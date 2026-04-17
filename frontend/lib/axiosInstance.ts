import axios from "axios";

const axiosInstance = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  // withCredentials: true,
});

axiosInstance.interceptors.request.use(async (config) => {
  const appKey = process.env.NEXT_PUBLIC_APP_KEY;
  if (appKey) {
    config.headers = config.headers ?? {};
    config.headers["app-key"] = appKey;
  }

  // Client-side only: attach NextAuth access token if not already provided.
  // Server Components / Route Handlers should pass Authorization explicitly.
  if (typeof window !== "undefined") {
    config.headers = config.headers ?? {};
    const hasAuthHeader =
      "Authorization" in config.headers ||
      "authorization" in (config.headers as Record<string, unknown>);

    if (!hasAuthHeader) {
      try {
        const { getSession } = await import("next-auth/react");
        const session = await getSession();
        if (session?.accessToken) {
          config.headers["Authorization"] = `Bearer ${session.accessToken}`;
        }
      } catch {
        // ignore: best-effort header attach
      }
    }
  }

  return config;
});

axiosInstance.interceptors.response.use(
  (response) => response,
  async (error) => {
    // Centralized handling for auth failures.
    if (typeof window !== "undefined" && error?.response?.status === 401) {
      try {
        const { signOut } = await import("next-auth/react");
        await signOut({ callbackUrl: "/" });
      } catch {
        // ignore
      }
    }
    return Promise.reject(error);
  }
);

export default axiosInstance;
