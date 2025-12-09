export interface ApiEndpoint {
  id: string;
  path: string;
  method: HttpMethod;
  parameters: EndpointParameter[];
  responseStructure: any;
  lastSeen: Date;
  riskScore: number;
  isActive: boolean;
  discoveredAt: Date;
  updatedAt: Date;
}

export enum HttpMethod {
  GET = 'GET',
  POST = 'POST',
  PUT = 'PUT',
  DELETE = 'DELETE',
  PATCH = 'PATCH',
  HEAD = 'HEAD',
  OPTIONS = 'OPTIONS'
}

export interface EndpointParameter {
  name: string;
  type: ParameterType;
  location: ParameterLocation;
  required: boolean;
  description?: string;
}

export enum ParameterType {
  STRING = 'string',
  NUMBER = 'number',
  BOOLEAN = 'boolean',
  OBJECT = 'object',
  ARRAY = 'array'
}

export enum ParameterLocation {
  QUERY = 'query',
  HEADER = 'header',
  BODY = 'body',
  PATH = 'path'
}

export interface ApiInventory {
  totalEndpoints: number;
  activeEndpoints: number;
  highRiskEndpoints: number;
  endpoints: ApiEndpoint[];
} 