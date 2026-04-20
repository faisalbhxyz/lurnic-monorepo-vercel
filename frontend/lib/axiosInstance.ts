import axios from "axios";

const browserApiBaseUrl = process.env.NEXT_PUBLIC_API_URL;
// Server-side (NextAuth, Server Actions, Route Handlers) should prefer the internal docker-network URL.
const serverApiBaseUrl =
  process.env.API_INTERNAL_URL ?? process.env.NEXT_PUBLIC_API_URL;

const axiosInstance = axios.create({
  baseURL: typeof window === "undefined" ? serverApiBaseUrl : browserApiBaseUrl,
  // withCredentials: true,
});

axiosInstance.interceptors.request.use(async (config) => {
  const isServer = typeof window === "undefined";
  if (isServer && !config.baseURL) {
    // Axios (Node/http adapter) requires absolute URLs. Without a baseURL, requests
    // like "/user/login" will fail with misleading network errors.
    throw new Error(
      "API base URL is not configured. Set API_INTERNAL_URL (server) or NEXT_PUBLIC_API_URL."
    );
  }

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

  // #region agent log
  // Only log from the browser. On the server, `fetch("/api/...")` is not a valid URL.
  if (!isServer) {
    fetch("/api/debug-log", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Debug-Session-Id": "01d620",
      },
      body: JSON.stringify({
        sessionId: "01d620",
        runId: "pre-fix",
        hypothesisId: "H3",
        location: "axiosInstance.ts:request",
        message: "Axios request (no secrets)",
        data: {
          isServer,
          baseURL: String(config.baseURL ?? ""),
          url: String(config.url ?? ""),
          method: String(config.method ?? "get"),
          hasAuthHeader: Boolean(
            config.headers &&
              ("Authorization" in config.headers ||
                "authorization" in (config.headers as Record<string, unknown>))
          ),
          hasAppKey: Boolean(process.env.NEXT_PUBLIC_APP_KEY),
        },
        timestamp: Date.now(),
      }),
    }).catch(() => {});
  }
  // #endregion agent log

  return config;
});

axiosInstance.interceptors.response.use(
  (response) => response,
  async (error) => {
    // #region agent log
    if (typeof window !== "undefined") {
      fetch("/api/debug-log", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Debug-Session-Id": "01d620",
        },
        body: JSON.stringify({
          sessionId: "01d620",
          runId: "pre-fix",
          hypothesisId: "H3",
          location: "axiosInstance.ts:response(error)",
          message: "Axios response error (no secrets)",
          data: {
            status: error?.response?.status ?? null,
            url: String(error?.config?.url ?? ""),
            baseURL: String(error?.config?.baseURL ?? ""),
            method: String(error?.config?.method ?? ""),
          },
          timestamp: Date.now(),
        }),
      }).catch(() => {});
    }
    // #endregion agent log

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
