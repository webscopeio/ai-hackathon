import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  async rewrites() {
    return [
      {
        source: "/_api/:path*",
        destination:
          process.env.NODE_ENV == "development"
            ? "http://localhost:8080/:path*"
            : "",
      },
    ];
  },
};

export default nextConfig;
