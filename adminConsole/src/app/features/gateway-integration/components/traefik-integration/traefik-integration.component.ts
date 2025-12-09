import { Component, OnInit } from '@angular/core';
import { GatewayIntegrationService } from '../../services/gateway-integration.service';
import { Integration, GatewayType, IntegrationStatus, CredentialType, Endpoint } from '../../../../core/models/gateway-integration.model';

@Component({
  selector: 'app-traefik-integration',
  templateUrl: './traefik-integration.component.html',
  styleUrls: ['./traefik-integration.component.scss']
})
export class TraefikIntegrationComponent implements OnInit {
  traefikIntegrations: Integration[] = [];
  loading = false;
  error: string | null = null;
  selectedIntegration: Integration | null = null;

  constructor(private gatewayIntegrationService: GatewayIntegrationService) { }

  ngOnInit(): void {
    this.loadTraefikIntegrations();
  }

  loadTraefikIntegrations(): void {
    this.loading = true;
    this.error = null;

    this.gatewayIntegrationService.getIntegrations().subscribe({
      next: (integrations) => {
        this.traefikIntegrations = integrations.filter(integration => integration.type === GatewayType.TRAEFIK);
        this.loading = false;
      },
      error: (error) => {
        this.error = 'Failed to load Traefik integrations: ' + error.message;
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
        alert(`Traefik connection test completed: ${health.status}`);
        this.selectedIntegration!.health = health;
      },
      error: (error) => {
        this.error = 'Failed to test Traefik connection: ' + error.message;
      }
    });
  }

  onSyncConfiguration(): void {
    if (!this.selectedIntegration) return;

    this.gatewayIntegrationService.syncIntegration(this.selectedIntegration.id).subscribe({
      next: (result) => {
        alert(`Traefik configuration sync completed: ${result.message}`);
      },
      error: (error) => {
        this.error = 'Failed to sync Traefik configuration: ' + error.message;
      }
    });
  }

  clearError(): void {
    this.error = null;
  }
} 