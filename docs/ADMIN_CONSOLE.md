# 🖥️ **Admin Console (Frontend)**

## **📱 Application Overview**
The Admin Console is a modern, responsive web application built with **Angular 16.2+** that provides a comprehensive interface for managing all ScopeAPI services and infrastructure.

## **🏗️ Architecture**
- **Framework**: Angular 16.2+ with TypeScript 5.1+
- **State Management**: RxJS for reactive state management
- **Routing**: Angular Router with lazy-loaded feature modules
- **Styling**: SCSS with modern CSS features and responsive design
- **Testing**: Jasmine/Karma for unit testing

## **📁 Project Structure**
```
adminConsole/
├── 📁 src/                          # Source code
│   ├── 📁 app/                       # Main application
│   │   ├── 📁 core/                  # Core services & guards
│   │   ├── 📁 shared/                # Shared components & pipes
│   │   ├── 📁 features/              # Feature modules
│   │   │   ├── 📁 dashboard/          # Main dashboard
│   │   │   ├── 📁 api-discovery/      # API discovery management
│   │   │   ├── 📁 threat-detection/   # Threat detection interface
│   │   │   ├── 📁 data-protection/    # Data protection controls
│   │   │   ├── 📁 attack-protection/  # Attack blocking interface
│   │   │   ├── 📁 gateway-integration/ # Gateway management
│   │   │   └── 📁 auth/               # Authentication & authorization
│   │   ├── 📄 app.component.*         # Root component
│   │   ├── 📄 app.module.ts           # Root module
│   │   └── 📄 app-routing.module.ts   # Main routing
│   ├── 📁 assets/                     # Static assets
│   ├── 📁 environments/               # Environment configs
│   └── 📄 main.ts                     # Application entry point
├── 📄 package.json                    # Dependencies & scripts
├── 📄 angular.json                    # Angular CLI configuration
├── 📄 tsconfig.json                   # TypeScript configuration
└── 📄 README.md                       # Frontend documentation
```

## **🚀 Development Commands**
```bash
# Install dependencies
cd adminConsole && npm install

# Start development server
npm start
# or
ng serve

# Build for production
npm run build
# or
ng build

# Run tests
npm test
# or
ng test

# Clean and fresh start
npm run fresh-start
```

## **🔧 Key Features**
- **Responsive Design** - Works on desktop, tablet, and mobile
- **Lazy Loading** - Feature modules load on-demand
- **Real-time Updates** - Live data from backend services
- **Role-based Access** - Different views for different user roles
- **Dark/Light Themes** - User preference support
- **Internationalization** - Multi-language support ready

## **📊 Feature Modules**
- **Dashboard** - Overview and system health monitoring
- **API Discovery** - Endpoint catalog and change tracking
- **Threat Detection** - Security monitoring and alerts
- **Data Protection** - PII detection and compliance
- **Attack Protection** - Blocking rules and policies
- **Gateway Integration** - Multi-gateway configuration
- **Authentication** - User management and access control

## **🔌 Integration**
- **Backend APIs** - RESTful communication with Go services
- **WebSocket Support** - Real-time notifications and updates
- **File Upload** - Configuration and policy file management
- **Export/Import** - Data portability and backup

## **🧪 Testing Strategy**
- **Unit Tests** - Component and service testing
- **Integration Tests** - API integration testing
- **E2E Tests** - End-to-end user workflow testing
- **Performance Tests** - Load and stress testing

## **📦 Build & Deployment**
- **Development** - Hot reload with `ng serve`
- **Production** - Optimized builds with `ng build`
- **Docker** - Containerized deployment ready
- **CI/CD** - Automated testing and deployment

## **🔗 Related Documentation**
- [Development Setup](../docs/DEVELOPMENT.md)
- [API Reference](../docs/API.md)
- [Architecture Overview](../docs/ARCHITECTURE.md)
- [Contributing Guide](../docs/CONTRIBUTING.md)

---

**📖 Back to [Main Documentation](../docs/INDEX.md) | [Project README](../README.md)**
