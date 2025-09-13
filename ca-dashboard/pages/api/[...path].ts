// pages/api/[...path].ts
import type { NextApiRequest, NextApiResponse } from "next";
import http from "http";
import { API_URL } from "@/config/constants";

export default function handler(req: NextApiRequest, res: NextApiResponse) {
  // Map Next.js /api/* requests to Go server's /api/v1/* paths
  const path = API_URL + (req.query.path as string[]).join("/");

  // Prepare headers, excluding problematic ones
  const { ...cleanHeaders } = req.headers;
  
  // Check if we should use Unix socket or HTTP
  const useUnixSocket = process.env.USE_UNIX_SOCKET === "true";
  
  const options: http.RequestOptions = useUnixSocket
    ? {
        socketPath: "/tmp/go.sock",
        path,
        method: req.method,
        headers: cleanHeaders,
      }
    : {
        hostname: "localhost",
        port: process.env.GO_PORT || 8080,
        path,
        method: req.method,
        headers: cleanHeaders,
      };

  const proxy = http.request(options, (backendRes) => {
    // Copy status code
    res.statusCode = backendRes.statusCode || 500;
    
    // Copy headers from backend response
    Object.entries(backendRes.headers).forEach(([key, value]) => {
      if (value) {
        res.setHeader(key, value);
      }
    });

    // Stream the response
    backendRes.pipe(res);
  });

  proxy.on("error", (err) => {
    console.error("Proxy error:", err);
    if (!res.headersSent) {
      res.status(500).json({ error: "Backend unavailable" });
    }
  });

  // Handle request timeout
  proxy.setTimeout(30000, () => {
    proxy.destroy();
    if (!res.headersSent) {
      res.status(504).json({ error: "Backend timeout" });
    }
  });

  // Handle request body for POST/PUT/PATCH
  if (req.body && ["GET", "POST", "UPDATE", "DELETE"].includes(req.method || "")) {
    const bodyData = typeof req.body === 'string' ? req.body : JSON.stringify(req.body);
    proxy.write(bodyData);
  }

  proxy.end();
}

// Important: Disable Next.js body parsing to handle raw requests
export const config = {
  api: {
    bodyParser: {
      sizeLimit: '10mb',
    },
  },
}
