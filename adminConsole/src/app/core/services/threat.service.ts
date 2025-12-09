import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { Threat, ThreatType, ThreatSeverity, ThreatStatus, ThreatStatistics, DetectionMethod } from '../models/threat.model';

@Injectable({
  providedIn: 'root'
})
export class ThreatService {

  constructor() { }

  // Threat Detection methods
  getThreats(): Observable<Threat[]> {
    // Simulate API call
    return of([]);
  }

  getThreatDetails(threatId: string): Observable<Threat> {
    // Simulate API call
    return of({
      id: threatId,
      type: ThreatType.ANOMALY,
      severity: ThreatSeverity.MEDIUM,
      status: ThreatStatus.NEW,
      title: 'Suspicious API Request',
      description: 'Unusual pattern detected in API traffic',
      sourceIp: '192.168.1.100',
      endpointId: 'endpoint-1',
      endpointPath: '/api/users',
      detectionMethod: DetectionMethod.ML,
      confidence: 0.85,
      riskScore: 75,
      indicators: [],
      firstSeen: new Date(),
      lastSeen: new Date(),
      count: 1,
      createdAt: new Date(),
      updatedAt: new Date()
    });
  }

  getThreatStatistics(): Observable<ThreatStatistics> {
    // Simulate API call
    return of({
      totalThreats: 0,
      threatsByType: {},
      threatsBySeverity: {},
      threatsByStatus: {},
      recentThreats: []
    });
  }

  updateThreatStatus(threatId: string, status: ThreatStatus): Observable<{ success: boolean; message: string }> {
    // Simulate API call
    return of({ success: true, message: 'Threat status updated successfully' });
  }

  blockThreat(threatId: string): Observable<{ success: boolean; message: string }> {
    // Simulate API call
    return of({ success: true, message: 'Threat blocked successfully' });
  }
} 