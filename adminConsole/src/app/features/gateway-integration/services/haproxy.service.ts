import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../../environments/environment';
import { Integration } from '../../../core/models/gateway-integration.model';

@Injectable({
  providedIn: 'root'
})
export class HAProxyService {
  private apiUrl = `${environment.gatewayApiUrl}/haproxy`;

  constructor(private http: HttpClient) {}

  getIntegrations(): Observable<Integration[]> {
    return this.http.get<Integration[]>(`${this.apiUrl}/integrations`);
  }

  getIntegration(id: string): Observable<Integration> {
    return this.http.get<Integration>(`${this.apiUrl}/integrations/${id}`);
  }

  createIntegration(integration: Partial<Integration>): Observable<Integration> {
    return this.http.post<Integration>(`${this.apiUrl}/integrations`, integration);
  }

  updateIntegration(id: string, integration: Partial<Integration>): Observable<Integration> {
    return this.http.put<Integration>(`${this.apiUrl}/integrations/${id}`, integration);
  }

  deleteIntegration(id: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/integrations/${id}`);
  }

  testConnection(integration: Partial<Integration>): Observable<{ success: boolean; message: string }> {
    return this.http.post<{ success: boolean; message: string }>(`${this.apiUrl}/test-connection`, integration);
  }

  getBackends(integrationId: string): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/integrations/${integrationId}/backends`);
  }

  getFrontends(integrationId: string): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/integrations/${integrationId}/frontends`);
  }

  getServers(integrationId: string): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/integrations/${integrationId}/servers`);
  }
} 