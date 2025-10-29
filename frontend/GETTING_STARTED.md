# Getting Started with EasyFund Frontend

## Quick Start

1. **Navigate to the frontend directory:**
   ```bash
   cd frontend
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

3. **Start the development server:**
   ```bash
   npm run dev
   ```

4. **Access the application:**
   - Open your browser and navigate to `http://localhost:3001`

## Prerequisites

- Node.js 18+ installed
- Backend API running on `http://localhost:8080`

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## Project Overview

This frontend application provides a complete user interface for the EasyFund loan management system. It includes:

### Pages
- **Home** - Landing page with features overview
- **Login** - User authentication
- **Register** - User registration
- **Dashboard** - Overview of loan applications and statistics
- **Applications** - List and manage loan applications

### Components
All UI components are organized in the `components` folder:
- Button - With multiple variants and sizes
- Input - Form inputs with validation
- Card - Container components
- Layout - Header, Footer, and main layout

### Styling
Base styles and theme configuration are in `src/styles/base.css`:
- CSS Variables for colors, fonts, spacing
- Dark theme support
- Responsive design utilities
- Consistent design system

### Authentication
Uses JWT tokens stored in localStorage with:
- AuthContext for state management
- Protected routes
- Automatic token refresh

## Backend Integration

The frontend communicates with the backend API:
- Base URL: `http://localhost:8080/api/v1`
- Authentication: JWT Bearer tokens
- Axios for HTTP requests

## Development Tips

1. **Hot Module Replacement**: Changes to files automatically reload in the browser
2. **TypeScript**: Full type safety across the application
3. **Path Aliases**: Use `@/` to reference the `src` directory
4. **Component Structure**: Each component in its own folder with CSS and index file

## Building for Production

```bash
npm run build
```

The production build will be in the `dist` folder and can be served with any static file server.

