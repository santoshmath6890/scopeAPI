export interface Threat {
  id: string;
  type: ThreatType;
  severity: ThreatSeverity;
  status: ThreatStatus;
  title: string;
  description: string;
  sourceIp: string;
  endpointId?: string;
  endpointPath?: string;
  detectionMethod: DetectionMethod;
  confidence: number;
  riskScore: number;
  indicators: ThreatIndicator[];
  requestData?: any;
  responseData?: any;
  firstSeen: Date;
  lastSeen: Date;
  count: number;
  createdAt: Date;
  updatedAt: Date;
}

export enum ThreatType {
  SQL_INJECTION = 'sql_injection',
  XSS = 'xss',
  CSRF = 'csrf',
  BOLA = 'bola',
  BROKEN_AUTH = 'broken_auth',
  DDoS = 'ddos',
  BRUTE_FORCE = 'brute_force',
  DATA_EXFILTRATION = 'data_exfiltration',
  ANOMALY = 'anomaly',
  SIGNATURE = 'signature'
}

export enum ThreatSeverity {
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high',
  CRITICAL = 'critical'
}

export enum ThreatStatus {
  NEW = 'new',
  INVESTIGATING = 'investigating',
  CONFIRMED = 'confirmed',
  FALSE_POSITIVE = 'false_positive',
  RESOLVED = 'resolved',
  BLOCKED = 'blocked'
}

export enum DetectionMethod {
  SIGNATURE = 'signature',
  BEHAVIORAL = 'behavioral',
  ANOMALY = 'anomaly',
  ML = 'machine_learning'
}

export interface ThreatIndicator {
  type: string;
  value: string;
  description: string;
  severity: ThreatSeverity;
  confidence: number;
}

export interface ThreatStatistics {
  totalThreats: number;
  threatsByType: { [key: string]: number };
  threatsBySeverity: { [key: string]: number };
  threatsByStatus: { [key: string]: number };
  recentThreats: Threat[];
} 