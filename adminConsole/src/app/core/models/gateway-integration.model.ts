export enum GatewayType {
  KONG = 'kong',
  NGINX = 'nginx',
  TRAEFIK = 'traefik',
  ENVOY = 'envoy',
  HAPROXY = 'haproxy'
}

export enum IntegrationStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  ERROR = 'error',
  PENDING = 'pending'
}

export enum CredentialType {
  BASIC = 'basic',
  TOKEN = 'token',
  API_KEY = 'api_key',
  TLS = 'tls'
}

export interface Credentials {
  type: CredentialType;
  username?: string;
  password?: string;
  token?: string;
  apiKey?: string;
  certificate?: string;
  privateKey?: string;
}

export interface Endpoint {
  id: string;
  name: string;
  url: string;
  method: string;
  protocol: string;
  port: number;
  timeout: number;
  description?: string;
  enabled: boolean;
}

export interface Integration {
  id: string;
  name: string;
  type: GatewayType;
  status: IntegrationStatus;
  config: Record<string, any>;
  credentials: Credentials;
  endpoints: Endpoint[];
  health?: HealthStatus;
  createdAt: Date;
  updatedAt: Date;
  lastSyncAt?: Date;
  lastSync?: string;
  errorMessage?: string;
}

export interface HealthStatus {
  status: string;
  message: string;
  lastCheck: string;
  latency: number;
}

export interface GatewayStats {
  totalIntegrations: number;
  activeIntegrations: number;
  errorIntegrations: number;
  pendingIntegrations: number;
  integrationsByType: Record<GatewayType, number>;
}

export interface Plugin {
  id: string;
  name: string;
  enabled: boolean;
  config: Record<string, any>;
  createdAt: Date;
}

export interface Route {
  id: string;
  name: string;
  protocols: string[];
  methods: string[];
  hosts: string[];
  paths: string[];
  strip_path: boolean;
  preserve_host: boolean;
  regex_priority: number;
  https_redirect_status_code: number;
  path_handling: string;
  created_at: number;
  updated_at: number;
  service: {
    id: string;
  };
}

export interface Service {
  id: string;
  name: string;
  protocol: string;
  host: string;
  port: number;
  path: string;
  retries: number;
  connect_timeout: number;
  write_timeout: number;
  read_timeout: number;
  created_at: number;
  updated_at: number;
} 