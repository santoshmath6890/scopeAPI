import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { GatewayIntegrationService } from '../../services/gateway-integration.service';
import { Integration } from '../../../../core/models/gateway-integration.model';
import { HealthStatus, SyncResult } from '../../services/gateway-integration.service';

@Component({
  selector: 'app-integration-details',
  templateUrl: './integration-details.component.html',
  styleUrls: ['./integration-details.component.scss']
})
export class IntegrationDetailsComponent implements OnInit {
  integration: Integration | null = null;
  loading = false;
  error: string | null = null;
  testing = false;
  syncing = false;
  lastSyncResult: SyncResult | null = null;

  constructor(
    private gatewayIntegrationService: GatewayIntegrationService,
    private route: ActivatedRoute,
    private router: Router
  ) { }

  // Add Object reference for template access
  Object = Object;

  ngOnInit(): void {
    this.loadIntegration();
  }

  loadIntegration(): void {
    this.loading = true;
    this.error = null;

    const integrationId = this.route.snapshot.paramMap.get('id');
    if (!integrationId) {
      this.error = 'Integration ID is required';
      this.loading = false;
      return;
    }

    this.gatewayIntegrationService.getIntegration(integrationId).subscribe({
      next: (integration) => {
        this.integration = integration;
        this.loading = false;
      },
      error: (error) => {
        this.error = 'Failed to load integration: ' + error.message;
        this.loading = false;
      }
    });
  }

  onTestConnection(): void {
    if (!this.integration) return;

    this.testing = true;
    this.error = null;

    this.gatewayIntegrationService.testIntegration(this.integration.id).subscribe({
      next: (health) => {
        this.integration!.health = health;
        this.testing = false;
        // Show success message
        alert(`Connection test completed: ${health.status}`);
      },
      error: (error) => {
        this.error = 'Failed to test connection: ' + error.message;
        this.testing = false;
      }
    });
  }

  onSyncConfiguration(): void {
    if (!this.integration) return;

    this.syncing = true;
    this.error = null;

    this.gatewayIntegrationService.syncIntegration(this.integration.id).subscribe({
      next: (result) => {
        this.lastSyncResult = result;
        this.syncing = false;
        // Refresh integration data
        this.loadIntegration();
        // Show success message
        alert(`Configuration sync completed: ${result.message}`);
      },
      error: (error) => {
        this.error = 'Failed to sync configuration: ' + error.message;
        this.syncing = false;
      }
    });
  }

  onEditIntegration(): void {
    if (!this.integration) return;
    this.router.navigate(['/gateway-integration/integrations', this.integration.id, 'edit']);
  }

  onDeleteIntegration(): void {
    if (!this.integration) return;

    if (confirm(`Are you sure you want to delete the integration "${this.integration.name}"? This action cannot be undone.`)) {
      this.gatewayIntegrationService.deleteIntegration(this.integration.id).subscribe({
        next: () => {
          this.router.navigate(['/gateway-integration/integrations']);
        },
        error: (error) => {
          this.error = 'Failed to delete integration: ' + error.message;
        }
      });
    }
  }

  onBackToList(): void {
    this.router.navigate(['/gateway-integration/integrations']);
  }

  getStatusColor(status: string): string {
    switch (status) {
      case 'active': return 'success';
      case 'inactive': return 'secondary';
      case 'error': return 'danger';
      case 'pending': return 'warning';
      default: return 'secondary';
    }
  }

  getHealthStatusColor(health: HealthStatus): string {
    switch (health.status) {
      case 'healthy': return 'success';
      case 'unhealthy': return 'danger';
      case 'degraded': return 'warning';
      default: return 'secondary';
    }
  }

  getGatewayTypeIcon(type: string): string {
    switch (type) {
      case 'kong': return '🔗';
      case 'nginx': return '⚡';
      case 'traefik': return '🚦';
      case 'envoy': return '🛡️';
      case 'haproxy': return '⚖️';
      default: return '🔧';
    }
  }

  getGatewayTypeDisplayName(type: string): string {
    switch (type) {
      case 'kong': return 'Kong';
      case 'nginx': return 'NGINX';
      case 'traefik': return 'Traefik';
      case 'envoy': return 'Envoy';
      case 'haproxy': return 'HAProxy';
      default: return type;
    }
  }

  formatDate(date: string | Date): string {
    return new Date(date).toLocaleString();
  }

  formatLatency(latency: number): string {
    if (latency < 1000) {
      return `${latency}ms`;
    } else {
      return `${(latency / 1000).toFixed(2)}s`;
    }
  }

  getConfigDisplayValue(value: any): string {
    if (typeof value === 'boolean') {
      return value ? 'Yes' : 'No';
    }
    if (typeof value === 'object') {
      return JSON.stringify(value, null, 2);
    }
    return String(value);
  }

  isConfigValueComplex(value: any): boolean {
    return typeof value === 'object' && value !== null;
  }

  getCredentialDisplayValue(value: string | undefined): string {
    if (!value) return 'Not configured';
    return value.length > 20 ? value.substring(0, 20) + '...' : value;
  }

  clearError(): void {
    this.error = null;
  }
} 