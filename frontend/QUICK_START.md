# Quick Start Guide - EasyFund Frontend

## Test Account

You can now test the application without a backend using these credentials:

**Email:** `admin@easyfund.com`  
**Password:** `admin123`

This will grant you full admin access to all pages.

## Getting Started

1. **Install dependencies:**
   ```bash
   cd frontend
   npm install
   ```

2. **Start the development server:**
   ```bash
   npm run dev
   ```

3. **Access the application:**
   - Open your browser to `http://localhost:3001`
   - Go to the Login page
   - Use the test account credentials above
   - You'll be automatically logged in and redirected to the Dashboard

## Design Features

The application now matches the design specifications:

- **Colors:**
  - Primary Blue: `#189CF4`
  - Dark Text: `#082131`
  - Light Background: `#F5F7FA`

- **Border Radius:** 20px on all interactive elements

- **Pages:**
  - Login page with role selection (Bank Client / Bank Employee)
  - Dashboard with statistics and application overview
  - Applications list with filtering

## Mock Authentication

The login now uses mock authentication when using the test account. This means:
- No JWT token is required
- No backend connection needed
- Full access to all protected pages
- User data is stored in localStorage

## Testing All Pages

Once logged in with the test account, you can access:
- `/` - Home page
- `/login` - Login page (you're already logged in)
- `/dashboard` - User dashboard
- `/applications` - Applications list

## Features

- Modern, clean UI matching the design
- Responsive layout for mobile and desktop
- Role-based UI elements
- Card-based layout with shadows
- Smooth transitions and hover effects

## Troubleshooting

If you encounter issues:
1. Clear your browser cache
2. Delete localStorage: `localStorage.clear()` in browser console
3. Restart the dev server

