# Cal Chatbot React UI

This is a React-based UI for the Cal Chatbot application, built with Next.js and shadcn/ui.

## Features

- Modern React component architecture
- Responsive design with Tailwind CSS
- Accessible UI components from shadcn/ui
- Real-time chat interface
- Integration with Cal.com scheduling

## Getting Started

### Prerequisites

- Node.js 18+ and npm/yarn
- Cal Chatbot backend running on http://localhost:8080

### Installation

1. Clone the repository
2. Install dependencies:

\`\`\`bash
npm install
# or
yarn install
\`\`\`

3. Run the development server:

\`\`\`bash
npm run dev
# or
yarn dev
\`\`\`

4. Open [http://localhost:3000](http://localhost:3000) in your browser

## Project Structure

- `app/` - Next.js app router pages and API routes
- `components/` - React components
- `hooks/` - Custom React hooks
- `types/` - TypeScript type definitions
- `lib/` - Utility functions

## Backend Integration

The React UI communicates with the Go backend via API calls to `/api/chat`. Make sure the backend is running on http://localhost:8080 before using the chat interface.
