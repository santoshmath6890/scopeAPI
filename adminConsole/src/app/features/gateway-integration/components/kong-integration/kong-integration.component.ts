import { Component, OnInit } from '@angular/core';
import { GatewayIntegrationService } from '../../services/gateway-integration.service';
import { Integration, GatewayType, IntegrationStatus, CredentialType, Endpoint } from '../../../../core/models/gateway-integration.model';

@Component({
  selector: 'app-kong-integration',
  templateUrl: './kong-integration.component.html',
  styleUrls: ['./kong-integration.component.scss']
})
export class KongIntegrationComponent implements OnInit {
  kongIntegrations: Integration[] = [];
  loading = false;
  error: string | null = null;
  selectedIntegration: Integration | null = null;

  // Kong-specific data
  services: any[] = [];
  routes: any[] = [];
  plugins: any[] = [];
  consumers: any[] = [];
  upstreams: any[] = [];

  constructor(private gatewayIntegrationService: GatewayIntegrationService) { }

  // Add JSON reference for template access
  JSON = JSON;

  ngOnInit(): void {
    this.loadKongIntegrations();
  }

  loadKongIntegrations(): void {
    this.loading = true;
    this.error = null;

    this.gatewayIntegrationService.getIntegrations().subscribe({
      next: (integrations) => {
        this.kongIntegrations = integrations.filter(integration => integration.type === GatewayType.KONG);
        this.loading = false;
      },
      error: (error) => {
        this.error = 'Failed to load Kong integrations: ' + error.message;
        this.loading = false;
      }
    });
  }

  selectIntegration(integration: Integration): void {
    this.selectedIntegration = integration;
    this.loadKongData(integration);
  }

  loadKongData(integration: Integration): void {
    // In a real implementation, this would call Kong-specific API endpoints
    // For now, we'll simulate the data
    this.services = [
      { id: '1', name: 'user-service', url: 'http://user-service:8080', protocol: 'http', host: 'user-service', port: 8080, path: '/', connect_timeout: 60000, write_timeout: 60000, read_timeout: 60000 },
      { id: '2', name: 'order-service', url: 'http://order-service:8081', protocol: 'http', host: 'order-service', port: 8081, path: '/', connect_timeout: 60000, write_timeout: 60000, read_timeout: 60000 }
    ];

    this.routes = [
      { id: '1', name: 'user-routes', protocols: ['http', 'https'], methods: ['GET', 'POST'], hosts: ['api.example.com'], paths: ['/users'], strip_path: true, preserve_host: false, regex_priority: 0, https_redirect_status_code: 426, path_handling: 'v0' },
      { id: '2', name: 'order-routes', protocols: ['http', 'https'], methods: ['GET', 'POST', 'PUT'], hosts: ['api.example.com'], paths: ['/orders'], strip_path: true, preserve_host: false, regex_priority: 0, https_redirect_status_code: 426, path_handling: 'v0' }
    ];

    this.plugins = [
      { id: '1', name: 'rate-limiting', service_id: '1', config: { minute: 100, hour: 1000 }, enabled: true },
      { id: '2', name: 'cors', service_id: '1', config: { origins: ['*'], methods: ['GET', 'POST'], headers: ['Content-Type'] }, enabled: true },
      { id: '3', name: 'jwt', service_id: '2', config: { secret: 'your-secret', key_claim_name: 'iss' }, enabled: true }
    ];

    this.consumers = [
      { id: '1', username: 'api-consumer', custom_id: 'consumer-1', tags: ['api'] },
      { id: '2', username: 'mobile-app', custom_id: 'consumer-2', tags: ['mobile'] }
    ];

    this.upstreams = [
      { id: '1', name: 'user-upstream', algorithm: 'round-robin', hash_on: 'none', hash_fallback: 'none', hash_on_header: null, hash_fallback_header: null, hash_on_cookie: null, hash_on_cookie_path: '/', slots: 10000, health_checks: { active: { type: 'http', timeout: 1, concurrency: 10, http_path: '/', https_verify_certificate: true, https_sni: null, healthy: { interval: 0, successes: 0, http_statuses: [200, 302], failures: 0 }, unhealthy: { interval: 0, http_failures: 0, tcp_failures: 0, timeouts: 0, http_statuses: [429, 404, 500, 501, 502, 503, 504, 505] } }, passive: { type: 'http', healthy: { http_statuses: [200, 201, 202, 203, 204, 205, 206, 207, 208, 226, 300, 301, 302, 303, 304, 305, 306, 307, 308], successes: 0 }, unhealthy: { http_failures: 0, tcp_failures: 0, timeouts: 0, http_statuses: [429, 500, 503] } } } }
    ];
  }

  onTestConnection(): void {
    if (!this.selectedIntegration) return;

    this.gatewayIntegrationService.testIntegration(this.selectedIntegration.id).subscribe({
      next: (health) => {
        alert(`Kong connection test completed: ${health.status}`);
        // Update the integration's health status
        this.selectedIntegration!.health = health;
      },
      error: (error) => {
        this.error = 'Failed to test Kong connection: ' + error.message;
      }
    });
  }

  onSyncConfiguration(): void {
    if (!this.selectedIntegration) return;

    this.gatewayIntegrationService.syncIntegration(this.selectedIntegration.id).subscribe({
      next: (result) => {
        alert(`Kong configuration sync completed: ${result.message}`);
        // Refresh Kong data
        this.loadKongData(this.selectedIntegration!);
      },
      error: (error) => {
        this.error = 'Failed to sync Kong configuration: ' + error.message;
      }
    });
  }

  onAddService(): void {
    // In a real implementation, this would open a modal to add a new service
    alert('Add Service functionality would be implemented here');
  }

  onAddRoute(): void {
    // In a real implementation, this would open a modal to add a new route
    alert('Add Route functionality would be implemented here');
  }

  onAddPlugin(): void {
    // In a real implementation, this would open a modal to add a new plugin
    alert('Add Plugin functionality would be implemented here');
  }

  onAddConsumer(): void {
    // In a real implementation, this would open a modal to add a new consumer
    alert('Add Consumer functionality would be implemented here');
  }

  onAddUpstream(): void {
    // In a real implementation, this would open a modal to add a new upstream
    alert('Add Upstream functionality would be implemented here');
  }

  onEditService(service: any): void {
    // In a real implementation, this would open a modal to edit the service
    alert(`Edit Service: ${service.name}`);
  }

  onDeleteService(service: any): void {
    if (confirm(`Are you sure you want to delete the service "${service.name}"?`)) {
      // In a real implementation, this would call the Kong API to delete the service
      alert(`Service "${service.name}" would be deleted`);
    }
  }

  onEditRoute(route: any): void {
    // In a real implementation, this would open a modal to edit the route
    alert(`Edit Route: ${route.name}`);
  }

  onDeleteRoute(route: any): void {
    if (confirm(`Are you sure you want to delete the route "${route.name}"?`)) {
      // In a real implementation, this would call the Kong API to delete the route
      alert(`Route "${route.name}" would be deleted`);
    }
  }

  onEditPlugin(plugin: any): void {
    // In a real implementation, this would open a modal to edit the plugin
    alert(`Edit Plugin: ${plugin.name}`);
  }

  onDeletePlugin(plugin: any): void {
    if (confirm(`Are you sure you want to delete the plugin "${plugin.name}"?`)) {
      // In a real implementation, this would call the Kong API to delete the plugin
      alert(`Plugin "${plugin.name}" would be deleted`);
    }
  }

  getPluginIcon(pluginName: string): string {
    switch (pluginName) {
      case 'rate-limiting': return '⏱️';
      case 'cors': return '🌐';
      case 'jwt': return '🔐';
      case 'oauth2': return '🔑';
      case 'key-auth': return '🗝️';
      case 'basic-auth': return '👤';
      case 'ip-restriction': return '🚫';
      case 'acl': return '📋';
      default: return '🔌';
    }
  }

  getServiceStatus(service: any): string {
    // In a real implementation, this would check the actual service health
    return 'healthy';
  }

  getRouteStatus(route: any): string {
    // In a real implementation, this would check the actual route status
    return 'active';
  }

  getPluginStatus(plugin: any): string {
    return plugin.enabled ? 'enabled' : 'disabled';
  }

  clearError(): void {
    this.error = null;
  }
} 