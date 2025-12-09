import { Component, OnInit } from '@angular/core';
import { GatewayIntegrationService } from '../../services/gateway-integration.service';
import { Integration, GatewayType, IntegrationStatus, CredentialType, Endpoint } from '../../../../core/models/gateway-integration.model';

@Component({
  selector: 'app-haproxy-integration',
  templateUrl: './haproxy-integration.component.html',
  styleUrls: ['./haproxy-integration.component.scss']
})
export class HAProxyIntegrationComponent implements OnInit {
  haproxyIntegrations: Integration[] = [];
  loading = false;
  error: string | null = null;
  selectedIntegration: Integration | null = null;

  constructor(private gatewayIntegrationService: GatewayIntegrationService) { }

  ngOnInit(): void {
    this.loadHAProxyIntegrations();
  }

  loadHAProxyIntegrations(): void {
    this.loading = true;
    this.error = null;

    this.gatewayIntegrationService.getIntegrations().subscribe({
      next: (integrations) => {
        this.haproxyIntegrations = integrations.filter(integration => integration.type === GatewayType.HAPROXY);
        this.loading = false;
      },
      error: (error) => {
        this.error = 'Failed to load HAProxy integrations: ' + error.message;
        this.loading = false;
      }
    });
  }

  selectIntegration(integration: Integration): void {
    this.selectedIntegration = integration;
  }

  onTestConnection(): void {
    if (!this.selectedIntegration) return;

    this.gatewayIntegrationService.testIntegration(this.selectedIntegration.id).subscribe({
      next: (health) => {
        alert(`HAProxy connection test completed: ${health.status}`);
        this.selectedIntegration!.health = health;
      },
      error: (error) => {
        this.error = 'Failed to test HAProxy connection: ' + error.message;
      }
    });
  }

  onSyncConfiguration(): void {
    if (!this.selectedIntegration) return;

    this.gatewayIntegrationService.syncIntegration(this.selectedIntegration.id).subscribe({
      next: (result) => {
        alert(`HAProxy configuration sync completed: ${result.message}`);
      },
      error: (error) => {
        this.error = 'Failed to sync HAProxy configuration: ' + error.message;
      }
    });
  }

  clearError(): void {
    this.error = null;
  }
} 