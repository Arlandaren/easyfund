# EasyFund Frontend

A modern React frontend application for the EasyFund loan management system.

## 🚀 Quick Start

### Test Account (No Backend Required)

**Email:** `admin@easyfund.com`  
**Password:** `admin123`

Use these credentials to test all features without a backend!

### Installation

```bash
cd frontend
npm install
npm run dev
```

Visit `http://localhost:3001` and log in with the test account above.

## Features

- **User Authentication**: Login and registration system
- **Mock Authentication**: Test without backend using admin credentials
- **Loan Applications**: Create and manage loan applications
- **Dashboard**: Overview of applications and statistics
- **Role-based Access**: Different views for borrowers, bank managers, and analysts
- **Responsive Design**: Modern UI with mobile support
- **Component Library**: Reusable UI components
- **Design Matching**: Colors (#189CF4, #082131), 20px border radius

## Tech Stack

- React 18.2
- TypeScript
- Vite
- React Router
- Axios
- Custom CSS with CSS Variables

## Getting Started

### Installation

```bash
# Install dependencies
npm install
```

### Development

```bash
# Start development server
npm run dev
```

The application will be available at `http://localhost:3001`

### Build

```bash
# Build for production
npm run build
```

### Preview Production Build

```bash
npm run preview
```

## Project Structure

```
frontend/
├── src/
│   ├── components/       # Reusable UI components
│   │   ├── Button/
│   │   ├── Input/
│   │   ├── Card/
│   │   └── Layout/
│   ├── pages/            # Page components
│   │   ├── Home/
│   │   ├── Login/
│   │   ├── Register/
│   │   ├── Dashboard/
│   │   └── Applications/
│   ├── context/          # React Context
│   │   └── AuthContext.tsx
│   ├── styles/           # Global styles
│   │   └── base.css
│   ├── utils/            # Utility functions
│   │   └── api.ts
│   ├── App.tsx           # Main App component
│   └── main.tsx          # Entry point
├── public/               # Static assets
├── package.json
├── tsconfig.json
├── vite.config.ts
└── README.md
```

## UI Components

### Button
- Variants: primary, secondary, outline, ghost, danger
- Sizes: sm, md, lg
- Supports loading state

### Input
- Label and helper text support
- Error states
- Left and right icon support
- Full width option

### Card
- Variants: default, outlined, elevated
- Clickable option

### Layout
- Header with navigation
- Footer
- Main content area

## API Integration

The frontend connects to the backend API running on `http://localhost:8080`.

API endpoints:
- `/api/v1/auth/*` - Authentication
- `/api/v1/users/*` - User management
- `/api/v1/loan-applications/*` - Loan applications
- `/api/v1/banks/*` - Bank information

## Styling

Styles are managed through CSS variables defined in `src/styles/base.css`:
- Colors (primary, secondary, semantic)
- Typography (fonts, sizes, weights)
- Spacing system
- Shadows and borders
- Transitions

## Authentication

Authentication is handled through:
- JWT tokens stored in localStorage
- AuthContext for state management
- Protected routes for authenticated pages
- Automatic redirect to login on 401 errors

## Development Notes

- The app uses Vite proxy configuration to forward `/api` requests to the backend
- All components are in separate folders with their CSS files
- TypeScript is used for type safety
- React Router handles client-side routing

