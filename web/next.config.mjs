/** @type {import('next').NextConfig} */
const nextConfig = {
  eslint: {
    ignoreDuringBuilds: true,
  },
  typescript: {
    ignoreBuildErrors: true,
  },
  images: {
    unoptimized: true,
  },
  async rewrites() {
    return [
      {
        source: '/api/cal/:path*',
        destination: 'http://localhost:8080/api/cal/:path*',
      },
    ]
  },
}

export default nextConfig
