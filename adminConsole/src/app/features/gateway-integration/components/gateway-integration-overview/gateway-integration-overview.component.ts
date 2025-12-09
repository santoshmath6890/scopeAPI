import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { GatewayIntegrationService } from '../../services/gateway-integration.service';
import { Integration, GatewayType, IntegrationStatus } from '../../../../core/models/gateway-integration.model';

@Component({
  selector: 'app-gateway-integration-overview',
  templateUrl: './gateway-integration-overview.component.html',
  styleUrls: ['./gateway-integration-overview.component.scss']
})
export class GatewayIntegrationOverviewComponent implements OnInit {
  integrations: Integration[] = [];
  loading = false;
  error: string | null = null;

  // Statistics
  totalIntegrations = 0;
  activeIntegrations = 0;
  errorIntegrations = 0;
  pendingIntegrations = 0;

  // Gateway type distribution
  gatewayTypeStats: { [key: string]: number } = {};

  // Make service public for template access and add Object/Math references
  constructor(
    public gatewayIntegrationService: GatewayIntegrationService,
    private router: Router
  ) { }

  // Add Object and Math references for template access
  Object = Object;
  Math = Math;

  ngOnInit(): void {
    this.loadIntegrations();
  }

  loadIntegrations(): void {
    this.loading = true;
    this.error = null;

    // For development, use mock data
    this.integrations = this.gatewayIntegrationService.getMockIntegrations();
    this.calculateStatistics();
    this.loading = false;

    // Uncomment for production
    /*
    this.gatewayIntegrationService.getIntegrations().subscribe({
      next: (integrations) => {
        this.integrations = integrations;
        this.calculateStatistics();
        this.loading = false;
      },
      error: (error) => {
        this.error = 'Failed to load integrations';
        this.loading = false;
        console.error('Error loading integrations:', error);
      }
    });
    */
  }

  calculateStatistics(): void {
    this.totalIntegrations = this.integrations.length;
    this.activeIntegrations = this.integrations.filter(i => i.status === IntegrationStatus.ACTIVE).length;
    this.errorIntegrations = this.integrations.filter(i => i.status === IntegrationStatus.ERROR).length;
    this.pendingIntegrations = this.integrations.filter(i => i.status === IntegrationStatus.PENDING).length;

    // Calculate gateway type distribution
    this.gatewayTypeStats = {};
    this.integrations.forEach(integration => {
      const type = integration.type;
      this.gatewayTypeStats[type] = (this.gatewayTypeStats[type] || 0) + 1;
    });
  }

  onAddIntegration(): void {
    this.router.navigate(['/gateway-integration/integrations/new']);
  }

  onViewIntegration(integration: Integration): void {
    this.router.navigate(['/gateway-integration/integrations', integration.id]);
  }

  onEditIntegration(integration: Integration): void {
    this.router.navigate(['/gateway-integration/integrations', integration.id, 'edit']);
  }

  onDeleteIntegration(integration: Integration): void {
    if (confirm(`Are you sure you want to delete the integration "${integration.name}"?`)) {
      this.gatewayIntegrationService.deleteIntegration(integration.id).subscribe({
        next: () => {
          this.loadIntegrations();
        },
        error: (error) => {
          this.error = 'Failed to delete integration';
          console.error('Error deleting integration:', error);
        }
      });
    }
  }

  onTestIntegration(integration: Integration): void {
    this.gatewayIntegrationService.testIntegration(integration.id).subscribe({
      next: (health) => {
        // Health status is not stored in the integration model
        console.log('Integration health:', health);
      },
      error: (error) => {
        this.error = 'Failed to test integration';
        console.error('Error testing integration:', error);
      }
    });
  }

  onSyncIntegration(integration: Integration): void {
    this.gatewayIntegrationService.syncIntegration(integration.id).subscribe({
      next: (result) => {
        if (result.success) {
          // Update last sync time
          const index = this.integrations.findIndex(i => i.id === integration.id);
          if (index !== -1) {
            this.integrations[index].lastSyncAt = new Date(result.timestamp);
          }
        } else {
          this.error = `Sync failed: ${result.message}`;
        }
      },
      error: (error) => {
        this.error = 'Failed to sync integration';
        console.error('Error syncing integration:', error);
      }
    });
  }

  getGatewayTypeDisplayName(type: GatewayType): string {
    return this.gatewayIntegrationService.getGatewayTypeDisplayName(type);
  }

  getGatewayTypeIcon(type: GatewayType): string {
    return this.gatewayIntegrationService.getGatewayTypeIcon(type);
  }

  getStatusDisplayName(status: IntegrationStatus): string {
    return this.gatewayIntegrationService.getStatusDisplayName(status);
  }

  getStatusColor(status: IntegrationStatus): string {
    return this.gatewayIntegrationService.getStatusColor(status);
  }

  getHealthStatusColor(status?: string): string {
    if (!status) return 'secondary';
    return this.gatewayIntegrationService.getHealthStatusColor(status);
  }

  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleString();
  }

  formatLatency(latency: number): string {
    return `${latency}ms`;
  }

  getGatewayTypeRoute(type: GatewayType): string {
    return `/gateway-integration/${type}`;
  }

  onGatewayTypeClick(type: GatewayType): void {
    const route = this.getGatewayTypeRoute(type);
    this.router.navigate([route]);
  }

  clearError(): void {
    this.error = null;
  }
} 