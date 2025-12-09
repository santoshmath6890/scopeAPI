import { Component, OnInit } from '@angular/core';
import { GatewayIntegrationService } from '../../services/gateway-integration.service';
import { Integration, GatewayType, IntegrationStatus, CredentialType, Endpoint } from '../../../../core/models/gateway-integration.model';

@Component({
  selector: 'app-nginx-integration',
  templateUrl: './nginx-integration.component.html',
  styleUrls: ['./nginx-integration.component.scss']
})
export class NginxIntegrationComponent implements OnInit {
  nginxIntegrations: Integration[] = [];
  loading = false;
  error: string | null = null;
  selectedIntegration: Integration | null = null;

  constructor(private gatewayIntegrationService: GatewayIntegrationService) { }

  ngOnInit(): void {
    this.loadNginxIntegrations();
  }

  loadNginxIntegrations(): void {
    this.loading = true;
    this.error = null;

    this.gatewayIntegrationService.getIntegrations().subscribe({
      next: (integrations) => {
        this.nginxIntegrations = integrations.filter(integration => integration.type === GatewayType.NGINX);
        this.loading = false;
      },
      error: (error) => {
        this.error = 'Failed to load NGINX integrations: ' + error.message;
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
        alert(`NGINX connection test completed: ${health.status}`);
        this.selectedIntegration!.health = health;
      },
      error: (error) => {
        this.error = 'Failed to test NGINX connection: ' + error.message;
      }
    });
  }

  onSyncConfiguration(): void {
    if (!this.selectedIntegration) return;

    this.gatewayIntegrationService.syncIntegration(this.selectedIntegration.id).subscribe({
      next: (result) => {
        alert(`NGINX configuration sync completed: ${result.message}`);
      },
      error: (error) => {
        this.error = 'Failed to sync NGINX configuration: ' + error.message;
      }
    });
  }

  clearError(): void {
    this.error = null;
  }
} 