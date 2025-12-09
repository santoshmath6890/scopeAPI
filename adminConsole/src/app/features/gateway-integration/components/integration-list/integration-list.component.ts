import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { GatewayIntegrationService } from '../../services/gateway-integration.service';
import { Integration, GatewayType, IntegrationStatus } from '../../../../core/models/gateway-integration.model';

@Component({
  selector: 'app-integration-list',
  templateUrl: './integration-list.component.html',
  styleUrls: ['./integration-list.component.scss']
})
export class IntegrationListComponent implements OnInit {
  integrations: Integration[] = [];
  filteredIntegrations: Integration[] = [];
  loading = false;
  error: string | null = null;
  
  // Filter properties
  selectedGatewayType: GatewayType | 'all' = 'all';
  selectedStatus: IntegrationStatus | 'all' = 'all';
  searchTerm = '';
  
  // Pagination
  currentPage = 1;
  itemsPerPage = 10;
  totalItems = 0;
  
  // Available filter options
  gatewayTypes = [
    { value: 'all', label: 'All Gateways' },
    { value: 'kong', label: 'Kong' },
    { value: 'nginx', label: 'NGINX' },
    { value: 'traefik', label: 'Traefik' },
    { value: 'envoy', label: 'Envoy' },
    { value: 'haproxy', label: 'HAProxy' }
  ];
  
  statuses = [
    { value: 'all', label: 'All Statuses' },
    { value: 'active', label: 'Active' },
    { value: 'inactive', label: 'Inactive' },
    { value: 'error', label: 'Error' },
    { value: 'pending', label: 'Pending' }
  ];

  constructor(
    private gatewayIntegrationService: GatewayIntegrationService,
    public router: Router
  ) { }

  // Add Math reference for template access
  Math = Math;

  ngOnInit(): void {
    this.loadIntegrations();
  }

  loadIntegrations(): void {
    this.loading = true;
    this.error = null;

    this.gatewayIntegrationService.getIntegrations().subscribe({
      next: (integrations) => {
        this.integrations = integrations;
        this.applyFilters();
        this.loading = false;
      },
      error: (error) => {
        this.error = 'Failed to load integrations: ' + error.message;
        this.loading = false;
      }
    });
  }

  applyFilters(): void {
    let filtered = [...this.integrations];

    // Filter by gateway type
    if (this.selectedGatewayType !== 'all') {
      filtered = filtered.filter(integration => integration.type === this.selectedGatewayType);
    }

    // Filter by status
    if (this.selectedStatus !== 'all') {
      filtered = filtered.filter(integration => integration.status === this.selectedStatus);
    }

    // Filter by search term
    if (this.searchTerm.trim()) {
      const term = this.searchTerm.toLowerCase();
      filtered = filtered.filter(integration => 
        integration.name.toLowerCase().includes(term) ||
        integration.type.toLowerCase().includes(term)
      );
    }

    this.filteredIntegrations = filtered;
    this.totalItems = filtered.length;
    this.currentPage = 1;
  }

  onGatewayTypeChange(): void {
    this.applyFilters();
  }

  onStatusChange(): void {
    this.applyFilters();
  }

  onSearchChange(): void {
    this.applyFilters();
  }

  clearFilters(): void {
    this.selectedGatewayType = 'all';
    this.selectedStatus = 'all';
    this.searchTerm = '';
    this.applyFilters();
  }

  get paginatedIntegrations(): Integration[] {
    const startIndex = (this.currentPage - 1) * this.itemsPerPage;
    const endIndex = startIndex + this.itemsPerPage;
    return this.filteredIntegrations.slice(startIndex, endIndex);
  }

  get totalPages(): number {
    return Math.ceil(this.totalItems / this.itemsPerPage);
  }

  get pages(): number[] {
    const pages: number[] = [];
    for (let i = 1; i <= this.totalPages; i++) {
      pages.push(i);
    }
    return pages;
  }

  onPageChange(page: number): void {
    this.currentPage = page;
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
          this.error = 'Failed to delete integration: ' + error.message;
        }
      });
    }
  }

  onTestIntegration(integration: Integration): void {
    this.gatewayIntegrationService.testIntegration(integration.id).subscribe({
      next: (result) => {
        alert(`Integration test completed: ${result.status}`);
        this.loadIntegrations(); // Refresh to get updated health status
      },
      error: (error) => {
        this.error = 'Failed to test integration: ' + error.message;
      }
    });
  }

  onSyncIntegration(integration: Integration): void {
    this.gatewayIntegrationService.syncIntegration(integration.id).subscribe({
      next: (result) => {
        alert(`Integration sync completed: ${result.message}`);
        this.loadIntegrations(); // Refresh to get updated sync status
      },
      error: (error) => {
        this.error = 'Failed to sync integration: ' + error.message;
      }
    });
  }

  getStatusColor(status: IntegrationStatus): string {
    switch (status) {
      case 'active': return 'success';
      case 'inactive': return 'secondary';
      case 'error': return 'danger';
      case 'pending': return 'warning';
      default: return 'secondary';
    }
  }

  getGatewayTypeIcon(type: GatewayType): string {
    switch (type) {
      case 'kong': return '🔗';
      case 'nginx': return '⚡';
      case 'traefik': return '🚦';
      case 'envoy': return '🛡️';
      case 'haproxy': return '⚖️';
      default: return '🔧';
    }
  }

  formatDate(date: string | Date): string {
    return new Date(date).toLocaleDateString();
  }

  clearError(): void {
    this.error = null;
  }
} 