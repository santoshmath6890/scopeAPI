import { Component, OnInit } from '@angular/core';
import { ApiService } from '../../../../core/services/api.service';
import { ApiEndpoint, ApiInventory } from '../../../../core/models/api-endpoint.model';

@Component({
  selector: 'app-api-discovery-overview',
  templateUrl: './api-discovery-overview.component.html',
  styleUrls: ['./api-discovery-overview.component.scss']
})
export class ApiDiscoveryOverviewComponent implements OnInit {
  apiInventory: ApiInventory | null = null;
  isLoading = false;

  constructor(private apiService: ApiService) {}

  ngOnInit(): void {
    this.loadApiInventory();
  }

  loadApiInventory(): void {
    this.isLoading = true;
    this.apiService.getApiInventory().subscribe({
      next: (inventory) => {
        this.apiInventory = inventory;
        this.isLoading = false;
      },
      error: (error) => {
        console.error('Error loading API inventory:', error);
        this.isLoading = false;
      }
    });
  }

  startDiscovery(): void {
    this.apiService.startDiscovery().subscribe({
      next: (result) => {
        if (result.success) {
          console.log('API discovery started:', result.message);
          // Refresh inventory after discovery
          setTimeout(() => this.loadApiInventory(), 2000);
        }
      },
      error: (error) => {
        console.error('Error starting API discovery:', error);
      }
    });
  }
}
