import { Component, OnInit } from '@angular/core';
import { GatewayIntegrationService } from '../../services/gateway-integration.service';
import { Integration, GatewayType, IntegrationStatus, CredentialType, Endpoint } from '../../../../core/models/gateway-integration.model';

@Component({
  selector: 'app-envoy-integration',
  templateUrl: './envoy-integration.component.html',
  styleUrls: ['./envoy-integration.component.scss']
})
export class EnvoyIntegrationComponent implements OnInit {
  envoyIntegrations: Integration[] = [];
  loading = false;
  error: string | null = null;
  selectedIntegration: Integration | null = null;

  constructor(private gatewayIntegrationService: GatewayIntegrationService) { }

  ngOnInit(): void {
    this.loadEnvoyIntegrations();
  }

  loadEnvoyIntegrations(): void {
    this.loading = true;
    this.error = null;

    this.gatewayIntegrationService.getIntegrations().subscribe({
      next: (integrations) => {
        this.envoyIntegrations = integrations.filter(integration => integration.type === GatewayType.ENVOY);
        this.loading = false;
      },
      error: (error) => {
        this.error = 'Failed to load Envoy integrations: ' + error.message;
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
        alert(`Envoy connection test completed: ${health.status}`);
        this.selectedIntegration!.health = health;
      },
      error: (error) => {
        this.error = 'Failed to test Envoy connection: ' + error.message;
      }
    });
  }

  onSyncConfiguration(): void {
    if (!this.selectedIntegration) return;

    this.gatewayIntegrationService.syncIntegration(this.selectedIntegration.id).subscribe({
      next: (result) => {
        alert(`Envoy configuration sync completed: ${result.message}`);
      },
      error: (error) => {
        this.error = 'Failed to sync Envoy configuration: ' + error.message;
      }
    });
  }

  clearError(): void {
    this.error = null;
  }
} 