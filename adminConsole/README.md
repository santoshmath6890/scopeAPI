# 🖥️ **ScopeAPI Admin Console**

[![Angular Version](https://img.shields.io/badge/Angular-16.2+-red.svg)](https://angular.io)
[![TypeScript Version](https://img.shields.io/badge/TypeScript-5.1+-blue.svg)](https://www.typescriptlang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](../LICENSE)

**Modern, responsive web application** built with **Angular 16.2+** that provides a comprehensive interface for managing all ScopeAPI services and infrastructure.

## 🏗️ **Architecture Overview**

```
┌──────────────────────────────────────────────────────────────────┐
│                    ScopeAPI Admin Console                        │
├──────────────────────────────────────────────────────────────────┤
│                    Presentation Layer                            │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │
│  │   Header    │ │   Sidebar   │ │   Content   │ │   Footer    │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │
├──────────────────────────────────────────────────────────────────┤
│                    Feature Modules                               │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │
│  │  Dashboard  │ │API Discovery│ │Threat Detect│ │Data Protect │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐                 │
│  │Attack Block │ │Gateway Integ│ │   Auth      │                 │
│  └─────────────┘ └─────────────┘ └─────────────┘                 │
├──────────────────────────────────────────────────────────────────┤
│                    Core Services                                 │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │
│  │   Routing   │ │   Guards    │ │ Interceptor │ │   Models    │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │
└──────────────────────────────────────────────────────────────────┘
```

## 📁 **Project Structure**

```
adminConsole/
├── 📁 src/                          # Source code
│   ├── 📁 app/                       # Main application
│   │   ├── 📁 core/                  # Core services & guards
│   │   │   ├── 📁 guards/            # Route protection
│   │   │   ├── 📁 interceptors/      # HTTP request/response handling
│   │   │   ├── 📁 models/            # Data models & interfaces
│   │   │   └── 📁 services/          # Core business logic
│   │   │
│   │   ├── 📁 shared/                # Shared components & utilities
│   │   │   ├── 📁 components/        # Reusable UI components
│   │   │   ├── 📁 pipes/             # Data transformation pipes
│   │   │   ├── 📁 directives/        # Custom directives
│   │   │   └── 📁 utils/             # Utility functions
│   │   │
│   │   ├── 📁 features/              # Feature modules (lazy-loaded)
│   │   │   ├── 📁 dashboard/          # Main dashboard & overview
│   │   │   │   ├── 📁 components/     # Dashboard components
│   │   │   │   └── 📁 dashboard.module.ts
│   │   │   │
│   │   │   ├── 📁 api-discovery/      # API discovery management
│   │   │   │   ├── 📁 components/     # API catalog components
│   │   │   │   └── 📁 api-discovery.module.ts
│   │   │   │
│   │   │   ├── 📁 threat-detection/   # Security monitoring
│   │   │   │   ├── 📁 components/     # Threat visualization
│   │   │   │   └── 📁 threat-detection.module.ts
│   │   │   │
│   │   │   ├── 📁 data-protection/    # PII & compliance
│   │   │   │   ├── 📁 components/     # Data protection UI
│   │   │   │   └── 📁 data-protection.module.ts
│   │   │   │
│   │   │   ├── 📁 attack-protection/  # Attack blocking interface
│   │   │   │   ├── 📁 components/     # Blocking rules UI
│   │   │   │   └── 📁 attack-protection.module.ts
│   │   │   │
│   │   │   ├── 📁 gateway-integration/ # Gateway management
│   │   │   │   ├── 📁 components/     # Gateway config UI
│   │   │   │   └── 📁 gateway-integration.module.ts
│   │   │   │
│   │   │   └── 📁 auth/               # Authentication & authorization
│   │   │       ├── 📁 components/     # Login, user management
│   │   │       └── 📁 auth.module.ts
│   │   │
│   │   ├── 📄 app.component.*         # Root component
│   │   ├── 📄 app.module.ts           # Root module
│   │   ├── 📄 app-routing.module.ts   # Main routing configuration
│   │   └── 📄 app.component.html      # Root template
│   │
│   ├── 📁 assets/                     # Static assets
│   │   ├── 📁 images/                 # Application images
│   │   ├── 📁 icons/                  # UI icons
│   │   └── 📁 styles/                 # Global styles
│   │
│   ├── 📁 environments/               # Environment configurations
│   │   ├── 📄 environment.ts           # Development environment
│   │   └── 📄 environment.prod.ts     # Production environment
│   │
│   ├── 📄 main.ts                     # Application entry point
│   ├── 📄 styles.scss                 # Global styles
│   └── 📄 index.html                  # Main HTML template
│
├── 📄 package.json                    # Dependencies & scripts
├── 📄 angular.json                    # Angular CLI configuration
├── 📄 tsconfig.json                   # TypeScript configuration
├── 📄 tsconfig.spec.json              # Test TypeScript config
├── 📄 tsconfig.app.json               # App TypeScript config
├── 📄 .editorconfig                   # Editor configuration
└── 📄 README.md                       # This file
```

## 🔧 **Technology Stack**

### **Frontend Framework**
- **Angular**: 16.2+ - Full-featured web application framework
- **TypeScript**: 5.1+ - Type-safe JavaScript development
- **RxJS**: 7.8+ - Reactive programming for state management
- **Angular CLI**: 16.2+ - Command-line interface and build tools

### **UI & Styling**
- **SCSS**: Advanced CSS with variables, mixins, and nesting
- **Angular Material**: Material Design components (ready to integrate)
- **Responsive Design**: Mobile-first, tablet, and desktop support
- **CSS Grid & Flexbox**: Modern layout techniques

### **Development Tools**
- **Webpack**: Module bundling and optimization
- **Karma**: Unit testing framework
- **Jasmine**: Testing framework
- **ESLint**: Code quality and style enforcement

### **Build & Deployment**
- **Angular CLI**: Build, serve, and deployment commands
- **Docker**: Containerized deployment ready
- **CI/CD**: GitHub Actions integration ready

## 🚀 **Quick Start**

### **Prerequisites**
```bash
# Install Node.js 18+ and npm
node --version
npm --version

# Install Angular CLI globally
npm install -g @angular/cli@16.2
```

### **Development Setup**
```bash
# Clone and navigate
git clone https://github.com/your-org/scopeapi.git
cd scopeapi/adminConsole

# Install dependencies
npm install

# Start development server
npm start
# or
ng serve

# Open browser to http://localhost:4200
```

### **Available Scripts**
```bash
# Development
npm start              # Start dev server with live reload
npm run build          # Build for production
npm run watch          # Build with file watching
npm run test           # Run unit tests
npm run e2e            # Run end-to-end tests

# Maintenance
npm run clean          # Clean all generated files
npm run clean:all      # Clean + remove node_modules
npm run clean:install  # Clean + fresh install
npm run fresh-start    # Complete fresh start
```

## 📊 **Feature Modules**

### **🏠 Dashboard Module**
- **Purpose**: System overview and health monitoring
- **Features**: 
  - Real-time service status
  - Performance metrics
  - System alerts and notifications
  - Quick access to all features

### **🔍 API Discovery Module**
- **Purpose**: Manage and monitor API endpoints
- **Features**:
  - Endpoint catalog and search
  - Change detection and versioning
  - API documentation generation
  - Metadata management

### **🛡️ Threat Detection Module**
- **Purpose**: Security monitoring and threat analysis
- **Features**:
  - Real-time threat alerts
  - Security dashboard
  - Threat intelligence feeds
  - Incident response tools

### **🔒 Data Protection Module**
- **Purpose**: Sensitive data management and compliance
- **Features**:
  - PII detection and classification
  - Compliance monitoring
  - Data flow visualization
  - Audit logging

### **⚡ Attack Protection Module**
- **Purpose**: Configure and monitor security rules
- **Features**:
  - Blocking rules management
  - Rate limiting configuration
  - IP whitelist/blacklist
  - Attack pattern analysis

### **🌐 Gateway Integration Module**
- **Purpose**: Multi-gateway configuration management
- **Features**:
  - Kong, Envoy, HAProxy support
  - Policy configuration
  - Health monitoring
  - Configuration sync

### **🔐 Authentication Module**
- **Purpose**: User management and access control
- **Features**:
  - User authentication
  - Role-based access control
  - Permission management
  - Session management

## 🎨 **UI Components**

### **Shared Components**
- **Header**: Navigation and user controls
- **Sidebar**: Feature navigation and quick actions
- **Loading Spinner**: Visual feedback during operations
- **Data Tables**: Sortable, filterable data display
- **Forms**: Validation and error handling
- **Modals**: Confirmation dialogs and forms
- **Charts**: Data visualization components

### **Design System**
- **Color Palette**: Consistent color scheme
- **Typography**: Readable font hierarchy
- **Spacing**: Consistent margins and padding
- **Icons**: Unified icon library
- **Animations**: Smooth transitions and feedback

## 🔌 **Backend Integration**

### **API Communication**
- **RESTful APIs**: HTTP-based service communication
- **WebSocket Support**: Real-time updates and notifications
- **Error Handling**: Graceful error management
- **Loading States**: Visual feedback during API calls

### **Authentication & Security**
- **JWT Tokens**: Secure authentication
- **HTTP Interceptors**: Automatic token management
- **Route Guards**: Protected route access
- **CSRF Protection**: Cross-site request forgery prevention

### **Data Management**
- **State Management**: RxJS-based reactive state
- **Caching**: Local storage and memory caching
- **Offline Support**: Basic offline functionality
- **Data Validation**: Client-side validation

## 🧪 **Testing Strategy**

### **Unit Testing**
- **Framework**: Jasmine + Karma
- **Coverage**: Target 80%+ code coverage
- **Components**: Individual component testing
- **Services**: Business logic testing
- **Pipes**: Data transformation testing

### **Integration Testing**
- **API Integration**: Backend service testing
- **Component Interaction**: Multi-component testing
- **Routing**: Navigation and route testing

### **End-to-End Testing**
- **User Workflows**: Complete user journey testing
- **Cross-browser**: Multiple browser compatibility
- **Performance**: Load time and responsiveness

## 📦 **Build & Deployment**

### **Development Build**
```bash
# Development server with hot reload
ng serve

# Development build
ng build --configuration development
```

### **Production Build**
```bash
# Production build with optimization
ng build --configuration production

# Build with specific environment
ng build --configuration production --environment prod
```

### **Docker Deployment**
```dockerfile
# Multi-stage build for production
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
```

## 🔧 **Configuration**

### **Environment Variables**
```typescript
// environment.ts
export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080',
  wsUrl: 'ws://localhost:8080',
  version: '1.0.0'
};
```

### **Angular Configuration**
```json
// angular.json
{
  "projects": {
    "admin-console": {
      "architect": {
        "build": {
          "options": {
            "outputPath": "dist/admin-console",
            "index": "src/index.html",
            "main": "src/main.ts"
          }
        }
      }
    }
  }
}
```

## 📚 **Documentation & Resources**

### **Project Documentation**
- **[Main Project README](../README.md)** - Project overview
- **[Backend Services](../backend/README.md)** - Go microservices
- **[Architecture Guide](../docs/ARCHITECTURE.md)** - System design
- **[API Reference](../docs/API.md)** - Service endpoints

### **Angular Resources**
- **[Angular Documentation](https://angular.io/docs)** - Official guides
- **[Angular CLI](https://cli.angular.io/)** - Command reference
- **[Angular Material](https://material.angular.io/)** - UI components
- **[RxJS Documentation](https://rxjs.dev/)** - Reactive programming

## 🤝 **Contributing**

### **Development Workflow**
1. **Fork** the repository
2. **Create** a feature branch
3. **Follow** Angular style guide
4. **Add** tests for new features
5. **Update** documentation
6. **Submit** pull request

### **Code Standards**
- Follow [Angular Style Guide](https://angular.io/guide/styleguide)
- Use TypeScript strict mode
- Write meaningful commit messages
- Maintain test coverage
- Document complex logic

### **Component Guidelines**
- Use OnPush change detection strategy
- Implement OnDestroy for cleanup
- Use async pipe with observables
- Follow single responsibility principle

## 📄 **License**

This project is licensed under the **MIT License** - see the [LICENSE](../LICENSE) file for details.

## 🔗 **Related Repositories**

- **[Backend Services](../backend/)** - Go microservices
- **[Infrastructure Scripts](../scripts/)** - Development automation
- **[Project Documentation](../docs/)** - Comprehensive guides

---

**🎯 Ready to build amazing user interfaces?**
- **Star** this repository if you find it useful
- **Fork** to contribute or customize
- **Share** with your team and community

**Happy coding! 🚀✨**
